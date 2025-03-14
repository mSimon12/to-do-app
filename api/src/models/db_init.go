package models

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Database interface for mocking
type Database interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Close()
}

var dbConn Database

//go:embed schema.sql
var dbInitQuery string

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

func getDatabaseConnection() Database {
	if dbConn != nil {
		return dbConn
	}

	// conn, err := pgx.Connect(context.Background(), getDatabaseUrl())
	conn, err := pgxpool.New(context.Background(), getDatabaseUrl())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}

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
