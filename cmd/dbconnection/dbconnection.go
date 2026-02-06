package dbconnection

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func GetPostgressConnetion() (*sql.DB, error) {
	db, err := sql.Open(getEnvValue("DB_PROTOCOL"), getConnectionString())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	return db, err
}

func getConnectionString() string {
	user := getEnvValue("DB_USER")
	password := getEnvValue("DB_PASSWORD")
	host := getEnvValue("DB_HOST")
	port := getEnvValue("DB_PORT")
	dbname := getEnvValue("DB_NAME")
	options := getEnvValue("DB_OPTIONS")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s", host, port, user, password, dbname, options)
}

func getEnvValue(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is not set in config.", key)
	}

	return value
}
