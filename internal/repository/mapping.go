package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"lis-khanza-mapper/internal/model"
)

type MappingRepo struct {
	db *sql.DB
}

func NewMappingRepo(db *sql.DB) *MappingRepo {
	return &MappingRepo{db: db}
}

func (r *MappingRepo) List(ctx context.Context, lisTestsPK uint64, idTemplate int, kdJenisPrw, status string) ([]model.Mapping, error) {
	sqlQ := `
SELECT m.id, m.lis_tests_pk, t.lis_test_id, t.test_name, m.id_template,
       IFNULL(tl.Pemeriksaan,''), m.kd_jenis_prw, IFNULL(j.nm_perawatan,''),
       m.status, m.created_by, m.created_at, m.updated_at
FROM lis_mapping_tests m
INNER JOIN lis_tests t ON t.id = m.lis_tests_pk
LEFT JOIN template_laboratorium tl ON tl.id_template = m.id_template
LEFT JOIN jns_perawatan_lab j ON j.kd_jenis_prw = m.kd_jenis_prw
WHERE 1=1`
	var args []any
	if lisTestsPK > 0 {
		sqlQ += " AND m.lis_tests_pk = ?"
		args = append(args, lisTestsPK)
	}
	if idTemplate > 0 {
		sqlQ += " AND m.id_template = ?"
		args = append(args, idTemplate)
	}
	if kdJenisPrw != "" {
		sqlQ += " AND m.kd_jenis_prw = ?"
		args = append(args, kdJenisPrw)
	}
	if status != "" {
		sqlQ += " AND m.status = ?"
		args = append(args, status)
	}
	sqlQ += " ORDER BY t.test_name, m.id_template"

	rows, err := r.db.QueryContext(ctx, sqlQ, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.Mapping
	for rows.Next() {
		var m model.Mapping
		if err := rows.Scan(&m.ID, &m.LisTestsPK, &m.LisTestID, &m.TestName, &m.IDTemplate,
			&m.Pemeriksaan, &m.KdJenisPrw, &m.NmPerawatan, &m.Status, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

func (r *MappingRepo) BulkUpsert(ctx context.Context, lisTestsPK uint64, idTemplates []int, createdBy string) (model.BulkMapResult, error) {
	var result model.BulkMapResult
	if lisTestsPK == 0 {
		return result, fmt.Errorf("lis_tests_pk is required")
	}
	if len(idTemplates) == 0 {
		return result, fmt.Errorf("no templates selected")
	}
	if len(idTemplates) > 200 {
		return result, fmt.Errorf("maximum 200 templates per request")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback() }()

	simrs := NewSimrsRepo(r.db)
	for _, idTpl := range idTemplates {
		kdJenisPrw, err := simrs.TemplateKdJenisPrw(ctx, idTpl)
		if errors.Is(err, sql.ErrNoRows) {
			result.Skipped++
			continue
		}
		if err != nil {
			return result, err
		}
		if kdJenisPrw == "" {
			result.Skipped++
			continue
		}

		var existingID uint64
		err = tx.QueryRowContext(ctx, `
SELECT id FROM lis_mapping_tests
WHERE lis_tests_pk=? AND id_template=? AND kd_jenis_prw=?`,
			lisTestsPK, idTpl, kdJenisPrw).Scan(&existingID)

		if errors.Is(err, sql.ErrNoRows) {
			_, err = tx.ExecContext(ctx, `
INSERT INTO lis_mapping_tests (lis_tests_pk, id_template, kd_jenis_prw, status, created_by)
VALUES (?,?,?,'aktif',?)`, lisTestsPK, idTpl, kdJenisPrw, createdBy)
			if err != nil {
				return result, err
			}
			result.Created++
		} else if err != nil {
			return result, err
		} else {
			_, err = tx.ExecContext(ctx, `
UPDATE lis_mapping_tests SET status='aktif', created_by=?, updated_at=NOW() WHERE id=?`,
				createdBy, existingID)
			if err != nil {
				return result, err
			}
			result.Updated++
		}
	}

	if err := tx.Commit(); err != nil {
		return result, err
	}
	return result, nil
}

func (r *MappingRepo) Deactivate(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE lis_mapping_tests SET status='nonaktif', updated_at=NOW() WHERE id=?`, id)
	return err
}

func (r *MappingRepo) Delete(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM lis_mapping_tests WHERE id=?`, id)
	return err
}

func (r *MappingRepo) CountActive(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM lis_mapping_tests WHERE status='aktif'`).Scan(&n)
	return n, err
}

func (r *MappingRepo) ListByLisTestPK(ctx context.Context, lisTestsPK uint64) ([]model.Mapping, error) {
	return r.List(ctx, lisTestsPK, 0, "", "aktif")
}

// FilterLisTestsPKByLisTestID supports API query by LIS testId string.
func FilterLisTestsPKByLisTestID(ctx context.Context, db *sql.DB, lisTestID string) (uint64, error) {
	lisTestID = strings.TrimSpace(lisTestID)
	var pk uint64
	err := db.QueryRowContext(ctx, `SELECT id FROM lis_tests WHERE lis_test_id=? AND status='aktif'`, lisTestID).Scan(&pk)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return pk, err
}
