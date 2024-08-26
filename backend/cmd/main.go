package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/mlemb/wallet/backend/internal/db"
	"github.com/mlemb/wallet/backend/internal/endpoint"
	_ "modernc.org/sqlite"
)

func main() {
	ctx := context.Background()

	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer conn.Close()

	log.Println("connected to database")

	db.SetDB(conn)

	if err := db.Migrate(ctx); err != nil {
		log.Fatalf("error migrating database: %v", err)
	}

	log.Println("database migrated")

	e := echo.New()

	v1 := e.Group("/api/v1")

	endpoint.NewTransferEndpoint().GroupRoutes(v1)

	log.Println("starting server")

	log.Fatalf("server error: %v", e.Start(":8080"))
}
