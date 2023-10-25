package repository

import "database/sql"

type MigrationRepository struct {
	db *sql.DB
}

func NewMigrationRepository(DB *sql.DB) *MigrationRepository {
	return &MigrationRepository{
		DB,
	}
}

func (repo MigrationRepository) Migrate() error {
	sql := "CREATE TABLE IF NOT EXISTS sources (" +
		"id SERIAL PRIMARY KEY NOT NULL," +
		"url VARCHAR(255) NOT NULL," +
		"prefix VARCHAR(255) NOT NULL," +
		"selector VARCHAR(255) NOT NULL," +
		"name VARCHAR(255) NOT NULL" +
		")"

	_, err := repo.db.Exec(sql)
	return err
}
