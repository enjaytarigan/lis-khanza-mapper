package db

import (
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("database ping: %w", err)
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	entries, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		data, err := migrationFS.ReadFile("migrations/" + name)
		if err != nil {
			return err
		}
		for _, stmt := range splitSQL(string(data)) {
			if _, err := db.Exec(stmt); err != nil {
				return fmt.Errorf("migrate %s: %w\nstmt: %s", name, err, stmt)
			}
		}
	}
	return nil
}

func splitSQL(s string) []string {
	var out []string
	var buf strings.Builder
	inSingleQuote := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\'' {
			buf.WriteByte(c)
			if inSingleQuote {
				if i+1 < len(s) && s[i+1] == '\'' {
					buf.WriteByte('\'')
					i++
				} else {
					inSingleQuote = false
				}
			} else {
				inSingleQuote = true
			}
			continue
		}
		if c == ';' && !inSingleQuote {
			if part := strings.TrimSpace(buf.String()); part != "" {
				out = append(out, part)
			}
			buf.Reset()
			continue
		}
		buf.WriteByte(c)
	}
	if part := strings.TrimSpace(buf.String()); part != "" {
		out = append(out, part)
	}
	return out
}
