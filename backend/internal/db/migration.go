package db

import (
	"context"

	"github.com/huandu/go-sqlbuilder"
)

func Migrate(ctx context.Context) error {
	transfers := sqlbuilder.NewCreateTableBuilder().CreateTable("transfers").IfNotExists()
	transfers.Define(`"id"`, "INTEGER", "PRIMARY KEY", "AUTOINCREMENT")
	transfers.Define(`"type"`, "TEXT")
	transfers.Define(`"from"`, "TEXT")
	transfers.Define(`"to"`, "TEXT")
	transfers.Define(`"amount"`, "REAL")
	transfers.Define(`"created_at"`, "DATETIME", "DEFAULT CURRENT_TIMESTAMP")
	transfers.Define(`"time"`, "DATETIME")
	transfersQuery, transfersArgs := transfers.Build()

	_, err := db.ExecContext(ctx, transfersQuery, transfersArgs...)
	if err != nil {
		return err
	}

	return nil
}
