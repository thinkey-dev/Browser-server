package mysql

import (
	"database/sql"
	"log/log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func (c *Config) NewMysql() (*sql.DB, error) {
	db, err := c.open()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (c *Config) open() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", c.DSN)
	if err != nil {
		log.Error("sql.Open() error (%v)", err)
		return nil, err
	}
	db.SetMaxOpenConns(c.Active)
	db.SetMaxIdleConns(c.Idle)
	db.SetConnMaxLifetime(time.Hour)
	return db, nil
}
