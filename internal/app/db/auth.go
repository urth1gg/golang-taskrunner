package db

import (
	"caravagio-api-golang/internal/app/models"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type APIKey struct {
	Key        string          `sql:"key"`
	Expiration models.NullTime `sql:"expiration"`
	UserID     string          `sql:"user_id"`
}

type AuthRepo interface {
	GetAPIKey(ctx context.Context, key string) (*APIKey, error)
}

type DBAuthRepo struct {
	db *sql.DB
}

func NewDBAuthRepo(db *sql.DB) *DBAuthRepo {
	return &DBAuthRepo{db: db}
}

func (r *DBAuthRepo) GetAPIKey(ctx context.Context, key string) (*APIKey, error) {
	var apiKey APIKey
	currentTime := time.Now().UTC()
	err := r.db.QueryRowContext(ctx, "SELECT `key`, expiration, user_id FROM api_keys WHERE `key` = ? AND expiration > ?", key, currentTime).Scan(
		&apiKey.Key,
		&apiKey.Expiration,
		&apiKey.UserID,
	)

	if err != nil {
		fmt.Println(err)
	}
	return &apiKey, err
}
