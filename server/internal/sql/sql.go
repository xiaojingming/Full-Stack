package sql

import (
	"context"
	"embed"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/xiaojingming/Full-Stack/server/internal/db"
)

//go:embed *.query.sql *.exec.sql
var queriesFS embed.FS

var queries = make(map[string]string)

func init() {
	entries, err := queriesFS.ReadDir(".")
	if err != nil {
		panic(fmt.Sprintf("failed to read embedded sql directory: %v", err))
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := queriesFS.ReadFile(entry.Name())
		if err != nil {
			panic(fmt.Sprintf("failed to read embedded sql file %s: %v", entry.Name(), err))
		}
		queries[entry.Name()] = strings.TrimSpace(string(data))
	}
}

func Query(ctx context.Context, name string, args ...any) ([]map[string]any, error) {
	sql, ok := queries[name+".query.sql"]
	if !ok {
		return nil, fmt.Errorf("query not found: %s", name)
	}

	rows, err := db.Pool().Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := rows.FieldDescriptions()
	var results []map[string]any

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		row := make(map[string]any, len(columns))
		for i, col := range columns {
			row[col.Name] = values[i]
		}
		results = append(results, row)
	}

	return results, rows.Err()
}

func Exec(ctx context.Context, name string, args ...any) (int64, error) {
	sql, ok := queries[name+".exec.sql"]
	if !ok {
		return 0, fmt.Errorf("exec not found: %s", name)
	}

	tag, err := db.Pool().Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func QueryRow(ctx context.Context, name string, args ...any) (pgx.Row, error) {
	sql, ok := queries[name+".query.sql"]
	if !ok {
		return nil, fmt.Errorf("query not found: %s", name)
	}

	return db.Pool().QueryRow(ctx, sql, args...), nil
}
