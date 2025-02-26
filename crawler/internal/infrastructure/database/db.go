package database

import (
	"crawler/internal/infrastructure/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db *sql.DB
}

func NewDB(cfg *config.Config) (*DB, error) {
	db, err := sql.Open("mysql", cfg.GetMySQLDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	db.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)       // 最大连接数
	db.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)       // 最大空闲连接数
	db.SetConnMaxLifetime(cfg.MySQL.ConnMaxLifetime) // 连接最大存活时间

	return &DB{db: db}, nil
}

func (d *DB) GetDB() *sql.DB {
	return d.db
}

func (d *DB) Close() error {
	return d.db.Close()
}
