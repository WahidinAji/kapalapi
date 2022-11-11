package pkg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func Migrate(ctx context.Context, db *pgx.Conn) (err error) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS vessel (
			id bigserial PRIMARY KEY,
			width serial NULL,
			length serial NULL,
			depth serial NULL,
			flag varchar(255) NULL,
			call_sign varchar(255) NULL,
			type serial NULL,
			name VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS transponder (
			id bigserial PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS vessel (
			id bigserial PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS licenses (
			id bigserial PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS engines (
			id bigserial PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS fishing_capacities (
			id bigserial PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS owner_operators (
			id bigserial PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)`,
	}

	for _, query := range queries {
		_, err = db.Exec(ctx, query)
		if err != nil {
			err = fmt.Errorf("migration error: %s", err.Error())
			return
		}
	}
	return
}
