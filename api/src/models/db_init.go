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
	var env_exist bool
	var dbUser string
	var dbPwd string
	var dbHost string
	var database string

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	if dbUser, env_exist = os.LookupEnv("POSTGRES_USER"); !env_exist {
		log.Fatalf("Missing POSTGRES_USER env variable")
	}

	if dbPwd, env_exist = os.LookupEnv("POSTGRES_PASSWORD"); !env_exist {
		log.Fatalf("Missing POSTGRES_PASSWORD env variable")
	}

	if dbHost, env_exist = os.LookupEnv("POSTGRES_HOST"); !env_exist {
		log.Fatalf("Missing POSTGRES_HOST env variable")
	}

	if database, env_exist = os.LookupEnv("POSTGRES_DB"); !env_exist {
		log.Fatalf("Missing POSTGRES_DB env variable")
	}

	urlDatabase := fmt.Sprintf("%v://%v:%v@db:5432/%v", dbHost, dbUser, dbPwd, database)

	return urlDatabase
}
