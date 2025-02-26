package database

import (
	"database/sql"
	"gateService/internal/infrastructure/config"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db *sql.DB
}

func NewDB(cfg *config.Config) *DB {
	db, err := sql.Open("mysql", cfg.GetMySQLDSN())
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MySQL.ConnMaxLifetime)
	return &DB{db: db}
}

func (d *DB) GetDB() *sql.DB {
	return d.db
}

func (d *DB) Close() {
	d.db.Close()
}
