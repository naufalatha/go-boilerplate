package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/naufalatha/go-boilerplate/config"
)

type Database struct {
	//
}

func InitDatabase(db *sql.DB) *Database {
	return &Database{}
}

// NewConnection Init new database connection, should be called only once per application
func NewConnection(config *config.Configuration) *sql.DB {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s&connect_timeout=%d", config.DbUsername, config.DbPassword, config.DbHost, config.DbPort, config.DbName, config.DbSSLMode, config.DbTimeout)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	// Ping to make sure database is online
	if err = db.Ping(); err != nil {
		panic(err)
	}
	return db
}

// Shutdown close database connection
func Shutdown(d *sql.DB) {
	err := d.Close()
	if err != nil {
		panic(err)
	}
}

// CheckConnection check database connection
func CheckConnection(db *sql.DB) error {
	return db.Ping()
}
