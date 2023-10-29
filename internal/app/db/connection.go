package db

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"  // Import your database driver
)

type Connection struct {
    DB *sql.DB
}

func NewConnection(dataSourceName string) (*Connection, error) {
    db, err := sql.Open("mysql", dataSourceName)  // Replace "mysql" with your database type
    if err != nil {
        return nil, err
    }
    return &Connection{DB: db}, nil
}