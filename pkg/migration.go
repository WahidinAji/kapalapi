package pkg

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func Migrate(ctx context.Context, db *pgx.Conn) error {
	err := db.Ping(ctx)
	if err != nil {
		return fmt.Errorf("connection lost to database : %s", err.Error())
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction error while try to migrating table : %s", err.Error())
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS user_keys(
			id bigserial PRIMARY KEY,
			uuid VARCHAR(40) NOT NULL,
			created_at timestamp default current_timestamp,
			updated_at timestamp default current_timestamp
		)
		`,
		`CREATE TABLE IF NOT EXISTS vessel (
			id bigserial PRIMARY KEY,
			user_key_id bigserial NOT NULL,
			name VARCHAR(255) NOT NULL,

			width bigint NULL,
			length bigint NULL,
			depth bigint NULL,
			flag VARCHAR(255) NULL,
			call_sign VARCHAR(255) NULL,
			type bigint NULL,
			imo VARCHAR(255) NULL,
			registration VARCHAR(255) NULL,
			mmsi VARCHAR(255) NULL,
			part_of_registration VARCHAR(255) NULL,
			external_marking VARCHAR(255) NULL,
			satellite_phone VARCHAR(255) NULL,
			dsc_number VARCHAR(255) NULL,
			max_crew bigint NULL,
			hull_material VARCHAR(255) NULL,
			stern_type VARCHAR(255) NULL,
			constructor VARCHAR(255) NULL,
			gross_tonnage double precision NULL,
			region_of_registration VARCHAR(255) NULL,

			transponder jsonb NULL,
			licenses jsonb NULL,
			engines jsonb NULL,
			fishing_capacity jsonb NULL,
			owner_operators jsonb NULL,

			preferred_image VARCHAR(255) NULL,

			created_by VARCHAR(40) NOT NULL,
			created_at timestamp default current_timestamp,
			updated_at timestamp default current_timestamp,
			CONSTRAINT fk_user_keys foreign key (user_key_id) REFERENCES user_keys(id)
		)`,
		`COMMENT on column vessel.created_by is  'generated from uuid on user_keys table'`,
	}

	for _, query := range queries {
		_, err = tx.Exec(ctx, query)
		if err != nil {
			errRollback := tx.Rollback(ctx)
			if errRollback != nil {
				return fmt.Errorf("rollback error while trying to migrate table : %s", errRollback.Error())
			}
			err = fmt.Errorf("creating table error: %s", err.Error())
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			return fmt.Errorf("rollback error while trying to commit migration : %s", errRollback.Error())
		}
		err = fmt.Errorf("commit error: %s", err.Error())
		return err
	}

	log.Printf("Migration completed!")
	return err
}
