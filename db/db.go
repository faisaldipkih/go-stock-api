package db

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
    var err error
    connStr := "user=postgres dbname=case_stock sslmode=disable password=admin123"
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal(err)
    }

    log.Println("Connected to the database successfully.")
}
