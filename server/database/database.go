package database

import (
    "context"
    "database/sql"
    "fmt"
    "os"
    "path"
    "time"

    _ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
    dir, err := os.Getwd()
    if err != nil {
        return fmt.Errorf("failed to get working directory: %w", err)
    }

    dbFile := path.Join(dir, "database", "database.db")

    if _, err := os.Stat(dbFile); os.IsNotExist(err) {
        if err := os.MkdirAll(path.Dir(dbFile), os.ModePerm); err != nil {
            return fmt.Errorf("failed to create database directory: %w", err)
        }

        file, err := os.Create(dbFile)
        if err != nil {
            return fmt.Errorf("failed to create database file: %w", err)
        }
        file.Close()
    }

    DB, err = sql.Open("sqlite", dbFile)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    if err := createTables(); err != nil {
        return fmt.Errorf("failed to create tables: %w", err)
    }

    return nil
}

func createTables() error {
    sqlStmt := `
        CREATE TABLE IF NOT EXISTS history_usd_brl (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            bid TEXT NOT NULL,
            timestamp INTEGER NOT NULL
        );
    `

    if _, err := DB.Exec(sqlStmt); err != nil {
        return fmt.Errorf("failed to execute table creation statement: %w", err)
    }

    return nil
}

func InsertNewExchangeRate(ctx context.Context, bid string) error {
    stmt, err := DB.PrepareContext(ctx, "INSERT INTO history_usd_brl(bid, timestamp) VALUES(?, ?)")
    if err != nil {
        return fmt.Errorf("failed to prepare insert statement: %w", err)
    }
    defer stmt.Close()

    if _, err := stmt.Exec(bid, time.Now().Unix()); err != nil {
        return fmt.Errorf("failed to execute insert statement: %w", err)
    }

    return nil
}