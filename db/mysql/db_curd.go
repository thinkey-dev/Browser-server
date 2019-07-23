package mysql

import "database/sql"

type ChainRepository struct {
	Db *sql.DB
}
