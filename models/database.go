package models

import (
        "database/sql"
        _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB() error {
        var err error
        db, err = sql.Open("sqlite3", "./cyclesync.db")
        if err != nil {
                return err
        }

        err = db.Ping()
        if err != nil {
                return err
        }

        return nil
}

// CloseDB closes the database connection
func CloseDB() {
        if db != nil {
                db.Close()
        }
}

// CreateTables creates the necessary tables if they don't exist
func CreateTables() error {
        // Create users table
        query := `
        CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                username TEXT NOT NULL UNIQUE,
                email TEXT NOT NULL UNIQUE,
                password TEXT NOT NULL,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`

        _, err := db.Exec(query)
        if err != nil {
                return err
        }

        // Create posts table
        query = `
        CREATE TABLE IF NOT EXISTS posts (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                user_id INTEGER NOT NULL,
                title TEXT NOT NULL,
                content TEXT NOT NULL,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY (user_id) REFERENCES users(id)
        );`

        _, err = db.Exec(query)
        return err
}
