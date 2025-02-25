package models

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

//go:embed schema.sql
var dbInitQuery string

func getDatabaseUrl() string {
	err := godotenv.Load("./models/database/.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	dbUser := os.Getenv("POSTGRES_USER")
	dbPwd := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	database := os.Getenv("POSTGRES_DB")

	urlDatabase := fmt.Sprintf("%v://%v:%v@localhost:5432/%v", dbHost, dbUser, dbPwd, database)

	return urlDatabase
}

func InitDatabase() {
	conn, err := pgx.Connect(context.Background(), getDatabaseUrl())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if _, err := conn.Exec(context.Background(), dbInitQuery); err != nil {
		panic(err)
	}

	defer conn.Close(context.Background())
}
