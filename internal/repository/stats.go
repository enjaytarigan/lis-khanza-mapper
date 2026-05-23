package repository

import (
	"context"
	"database/sql"

	"lis-khanza-mapper/internal/model"
)

func Dashboard(ctx context.Context, db *sql.DB) (model.DashboardStats, error) {
	var s model.DashboardStats
	_ = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM lis_tests WHERE status='aktif'`).Scan(&s.LisTestsActive)
	_ = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM lis_mapping_tests WHERE status='aktif'`).Scan(&s.MappingsActive)
	_ = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM template_laboratorium`).Scan(&s.TemplateCount)
	return s, nil
}
