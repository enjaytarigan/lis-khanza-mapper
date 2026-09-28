package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"lis-khanza-mapper/internal/model"
)

type BridgingLogRepo struct {
	db *sql.DB
}

func NewBridgingLogRepo(db *sql.DB) *BridgingLogRepo {
	return &BridgingLogRepo{db: db}
}

func (r *BridgingLogRepo) List(ctx context.Context, q, status string, page, limit int) (model.BridgingLogListResult, error) {
	out := model.BridgingLogListResult{Page: page, Limit: limit}
	if page < 1 {
		page = 1
		out.Page = 1
	}
	if limit <= 0 || limit > 200 {
		limit = 50
		out.Limit = 50
	}
	offset := (page - 1) * limit

	where := "WHERE 1=1"
	var args []any
	q = strings.TrimSpace(q)
	status = strings.TrimSpace(status)
	if status != "" {
		where += " AND status = ?"
		args = append(args, status)
	}
	if q != "" {
		where += " AND (no_order LIKE ? OR no_laboratorium LIKE ? OR IFNULL(error_message,'') LIKE ?)"
		like := "%" + q + "%"
		args = append(args, like, like, like)
	}

	countQ := "SELECT COUNT(*) FROM lis_hasil_inbox " + where
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&out.Total); err != nil {
		return out, err
	}

	listQ := `
SELECT id, no_order, no_laboratorium, status, mapped_count, unmapped_count, detail_written,
       IFNULL(error_message,''), created_at, updated_at
FROM lis_hasil_inbox ` + where + `
ORDER BY id DESC
LIMIT ? OFFSET ?`
	listArgs := append(append([]any{}, args...), limit, offset)
	rows, err := r.db.QueryContext(ctx, listQ, listArgs...)
	if err != nil {
		return out, err
	}
	defer rows.Close()

	for rows.Next() {
		var m model.BridgingLog
		if err := rows.Scan(
			&m.ID, &m.NoOrder, &m.NoLaboratorium, &m.Status,
			&m.MappedCount, &m.UnmappedCount, &m.DetailWritten,
			&m.ErrorMessage, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return out, err
		}
		out.Items = append(out.Items, m)
	}
	return out, rows.Err()
}

func (r *BridgingLogRepo) GetByID(ctx context.Context, id uint64) (*model.BridgingLog, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, no_order, no_laboratorium, status, mapped_count, unmapped_count, detail_written,
       IFNULL(error_message,''), IFNULL(payload_json,''), created_at, updated_at
FROM lis_hasil_inbox WHERE id = ?`, id)
	var m model.BridgingLog
	if err := row.Scan(
		&m.ID, &m.NoOrder, &m.NoLaboratorium, &m.Status,
		&m.MappedCount, &m.UnmappedCount, &m.DetailWritten,
		&m.ErrorMessage, &m.PayloadJSON, &m.CreatedAt, &m.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}
