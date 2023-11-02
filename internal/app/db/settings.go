package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Settings struct {
	SettingID string
	APIKey    sql.NullString
	UserID    sql.NullString
	PromptID  sql.NullString
}

type SettingsRepo interface {
	GetSetting(ctx context.Context, settingID string) (*Settings, error)
}

type DBSettingsRepo struct {
	db *sql.DB
}

func NewDBSettingsRepo(db *sql.DB) *DBSettingsRepo {
	return &DBSettingsRepo{db: db}
}

func (r *DBSettingsRepo) GetSetting(ctx context.Context, settingID string) (*Settings, error) {
	var setting Settings
	query := "SELECT setting_id, api_key, user_id, prompt_id FROM settings WHERE user_id = ?"
	err := r.db.QueryRowContext(ctx, query, settingID).Scan(
		&setting.SettingID,
		&setting.APIKey,
		&setting.UserID,
		&setting.PromptID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &setting, nil
}
