package repository

import (
	"context"
	"database/sql"
	"strings"

	"lis-khanza-mapper/internal/model"
)

type SimrsRepo struct {
	db *sql.DB
}

func NewSimrsRepo(db *sql.DB) *SimrsRepo {
	return &SimrsRepo{db: db}
}

func (r *SimrsRepo) ListPanels(ctx context.Context) ([]model.Panel, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT kd_jenis_prw, nm_perawatan FROM jns_perawatan_lab ORDER BY nm_perawatan`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.Panel
	for rows.Next() {
		var p model.Panel
		if err := rows.Scan(&p.KdJenisPrw, &p.NmPerawatan); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *SimrsRepo) ListTemplates(ctx context.Context, kdJenisPrw, q string, lisTestsPK uint64, limit int) ([]model.Template, error) {
	if limit <= 0 || limit > 500 {
		limit = 500
	}
	sqlQ := `
SELECT t.id_template, t.Pemeriksaan, t.satuan, IFNULL(t.urut,0), t.kd_jenis_prw, j.nm_perawatan,
       EXISTS(
         SELECT 1 FROM lis_mapping_tests m
         WHERE m.id_template = t.id_template AND m.status = 'aktif'
           AND m.lis_tests_pk = ?
       ) AS mapped
FROM template_laboratorium t
INNER JOIN jns_perawatan_lab j ON t.kd_jenis_prw = j.kd_jenis_prw
WHERE 1=1`
	args := []any{lisTestsPK}
	if kdJenisPrw != "" {
		sqlQ += " AND t.kd_jenis_prw = ?"
		args = append(args, kdJenisPrw)
	}
	if q != "" {
		sqlQ += " AND t.Pemeriksaan LIKE ?"
		args = append(args, "%"+q+"%")
	}
	sqlQ += " ORDER BY j.nm_perawatan, t.urut LIMIT ?"
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, sqlQ, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.Template
	for rows.Next() {
		var t model.Template
		var mapped int
		if err := rows.Scan(&t.IDTemplate, &t.Pemeriksaan, &t.Satuan, &t.Urut, &t.KdJenisPrw, &t.NmPerawatan, &mapped); err != nil {
			return nil, err
		}
		t.Mapped = mapped == 1
		list = append(list, t)
	}
	return list, rows.Err()
}

func (r *SimrsRepo) TemplateExists(ctx context.Context, idTemplate int) (bool, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM template_laboratorium WHERE id_template=?`, idTemplate).Scan(&n)
	return n > 0, err
}

func (r *SimrsRepo) TemplateKdJenisPrw(ctx context.Context, idTemplate int) (string, error) {
	var kd string
	err := r.db.QueryRowContext(ctx, `
SELECT kd_jenis_prw FROM template_laboratorium WHERE id_template=?`, idTemplate).Scan(&kd)
	return strings.TrimSpace(kd), err
}

func (r *SimrsRepo) CountTemplates(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM template_laboratorium`).Scan(&n)
	return n, err
}
