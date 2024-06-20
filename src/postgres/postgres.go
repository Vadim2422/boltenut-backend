package postgres

import (
	"boltenut/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/tracelog"
	"os"
)

// ConnectPostgres инициализирует соединение с базой данных
func ConnectPostgres() *pgx.Conn {
	dbHost, dbPort, dbUser, dbPassword, dbName :=
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	connConfig, _ := pgx.ParseConfig(dsn)
	connConfig.Tracer = &tracelog.TraceLog{Logger: logger.Logger{}, LogLevel: tracelog.LogLevelDebug}
	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}
