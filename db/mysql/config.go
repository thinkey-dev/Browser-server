package mysql

import (
	"database/sql"
	"strings"
	"time"

	"github.com/go-ini/ini"
)

type Config struct {
	DSN         string
	Active      int       // pool
	Idle        int       // pool
	IdleTimeout time.Time // connect max life time.
}

type MySql struct {
	Master Config
}

const (
	x = "_micro_auth_" // 密码加密参数
)

var (
	alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	userQ    = map[string]string{
		"update": "UPDATE %s.%s set password=?, salt=?, email=?, userName=?, company=? where id=?",
	}
	st = map[string]*sql.Stmt{}
)
var MySQL MySql

func init() {
	cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		panic(err)
	}
	userName := cfg.Section("mysql").Key("user").String()
	host := cfg.Section("mysql").Key("host").String()
	port := cfg.Section("mysql").Key("port").String()
	password := cfg.Section("mysql").Key("password").String()
	dbname := cfg.Section("mysql").Key("dbname").String()
	active := cfg.Section("mysql").Key("active").MustInt()
	idle := cfg.Section("mysql").Key("idle").MustInt()
	path := strings.Join([]string{userName, ":", password, "@tcp(", host, ":", port, ")/", dbname, "?charset=utf8"}, "")
	MySQL.Master = Config{DSN: path, Active: active, Idle: idle, IdleTimeout: time.Time{}}
}
