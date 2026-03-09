package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectDB opens a MySQL connection using environment variables.
func ConnectDB() *sql.DB {
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	pass := strings.TrimSpace(os.Getenv("DB_PASSWORD"))
	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	name := strings.TrimSpace(os.Getenv("DB_NAME"))

	if user == "" || name == "" {
		log.Println("DB_USER or DB_NAME not set — running without database")
		return nil
	}
	if host == "" {
		host = "localhost"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", user, pass, host, name)

	var db *sql.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("DB connect attempt %d failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Printf("Could not connect to MySQL: %v — running without database", err)
		return nil
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to MySQL")
	return db
}
