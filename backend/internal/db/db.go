package db

import (
	"database/sql"

	"github.com/huandu/go-sqlbuilder"
)

func init() {
	sqlbuilder.DefaultFlavor = sqlbuilder.SQLite
}

func SetDB(sqlDB *sql.DB) {
	db = sqlDB
}

var db *sql.DB
