package pkg

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func Mig(ctx context.Context, db *pgx.Conn) error {
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
			uuid VARCHAR(40) NOT NULL,
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

			preferred_image VARCHAR(255) NULL,

			created_by VARCHAR(40) NOT NULL,
			created_at timestamp default current_timestamp,
			updated_at timestamp default current_timestamp,
			CONSTRAINT fk_user_keys foreign key (user_key_id) REFERENCES user_keys(id)
		)`,
		`COMMENT on column vessel.created_by is  'generated from uuid on user_keys table'`,
		`CREATE UNIQUE INDEX IF NOT EXISTS uuid_idx ON vessel (uuid)`,
		`CREATE INDEX IF NOT EXISTS created_by_idx ON vessel (created_by)`,

		`CREATE TABLE IF NOT EXISTS transponder (
			id bigserial PRIMARY KEY,
			vessel_uuid VARCHAR(40) NOT NULL,
			install_date bigint NULL,
			install_company VARCHAR(255) NULL,
			installer_name VARCHAR(255) NULL,
			installion_port VARCHAR(255) NULL,
			install_latitude double precision NULL,
			install_longtitude double precision NULL,
			install_height double precision NULL,
			vendor_type VARCHAR(255) NULL,
			transceiver_manufacturer VARCHAR(255) NULL,
			serial_number VARCHAR(510) NOT NULL
		)`,
		// CONSTRAINT fk_vessel foreign key (vessel_uuid) REFERENCES vessel(uuid)
		`CREATE INDEX IF NOT EXISTS fk_transponder_uuid ON transponder (vessel_uuid)`,

		`CREATE TABLE IF NOT EXISTS licenses (
			id bigserial PRIMARY KEY,
			vessel_uuid VARCHAR(40) NOT NULL,
			type VARCHAR(255) NULL,
			expiry bigint NULL
		)`,
		// CONSTRAINT fk_vessel foreign key (vessel_uuid) REFERENCES vessel(uuid)
		`CREATE INDEX IF NOT EXISTS fk_licenses_uuid ON licenses (vessel_uuid)`,

		`CREATE TABLE IF NOT EXISTS engines (
			id bigserial PRIMARY KEY,
			vessel_uuid VARCHAR(40) NOT NULL,
			power double precision NULL,
			type VARCHAR(255) NULL
		)`,
		// CONSTRAINT fk_vessel foreign key (vessel_uuid) REFERENCES vessel(uuid)
		`CREATE INDEX IF NOT EXISTS fk_engines_uuid ON engines (vessel_uuid)`,
		`CREATE TABLE IF NOT EXISTS fishing_capacities (
			id bigserial PRIMARY KEY,
			vessel_uuid VARCHAR(40) NOT NULL,
			sub_type int NULL,
			group_seine_fishing bool NULL,
			main_gear VARCHAR(255) NULL,
			subsidiary_gear VARCHAR(255) NULL,
			freezer_snap bool NULL,
			freezer_ice bool NULL,
			freezer_seawater_refrigerated bool NULL,
			freezer_seawater_chilled bool NULL,
			freezer_blast_or_dry bool NULL,
			freezer_other bool NULL,
			freezer_hold_capacity double precision NULL
		)`,
		// CONSTRAINT fk_vessel foreign key (vessel_uuid) REFERENCES vessel(uuid)
		`CREATE INDEX IF NOT EXISTS fk_fishing_capacities_uuid ON fishing_capacities (vessel_uuid)`,
		`CREATE TABLE IF NOT EXISTS owner_operators (
			id bigserial PRIMARY KEY,
			vessel_uuid VARCHAR(40) NOT NULL,
			role VARCHAR(255) NULL,
			nationality VARCHAR(255) NULL,
			address VARCHAR(255) NULL,
			email VARCHAR(255) NULL,
			phone_number_1 VARCHAR(255) NULL,
			phone_number_2 VARCHAR(255) NULL,
			mobile_1 VARCHAR(255) NULL,
			mobile_2 VARCHAR(255) NULL,
			current bool NULL,
			preferred_image VARCHAR(255) NULL
		)`,
		// CONSTRAINT fk_vessel foreign key (vessel_uuid) REFERENCES vessel(uuid)
		`CREATE INDEX IF NOT EXISTS fk_owner_operators_uuid ON owner_operators (vessel_uuid)`,
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
