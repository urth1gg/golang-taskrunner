package db

import (
	"context"
	"database/sql"
	"time"
	"fmt"
	"caravagio-api-golang/internal/app/models"
)

type APIKey struct {
	Key       string
	Expiration models.NullTime
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
	err := r.db.QueryRowContext(ctx, "SELECT `key`, expiration FROM api_keys WHERE `key` = ? AND expiration > ?", key, currentTime).Scan(
		&apiKey.Key,
		&apiKey.Expiration,
	)

	if err != nil {
		fmt.Println(err)
	}
	return &apiKey, err
}
