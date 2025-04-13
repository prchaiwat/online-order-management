package db

import (
	"database/sql"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func ConnectDB(connection string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}
	maxOpenCons, _ := strconv.Atoi(getEnvWithDefault("DB_MAX_OPEN_CONNS", "20"))
	maxIdleConns, _ := strconv.Atoi(getEnvWithDefault("DB_MAX_IDLE_CONNS", "10"))
	connMaxLifetime, _ := strconv.Atoi(getEnvWithDefault("DB_CONN_MAX_LIFETIME", "60"))

	db.SetMaxOpenConns(maxOpenCons)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	return db, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
