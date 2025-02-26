package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

type Config struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
}

func InitDB(conf *Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 设置连接池
	DB.SetMaxIdleConns(10)
	DB.SetMaxOpenConns(100)
	DB.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("ping数据库失败: %v", err)
	}

	log.Println("数据库连接成功")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
