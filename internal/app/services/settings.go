package services

import (
	"caravagio-api-golang/internal/app/db"
	"context"
	"fmt"
)

type SettingsService struct {
	db db.SettingsRepo
}

func NewSettingsService(db db.SettingsRepo) *SettingsService {
	return &SettingsService{db: db}
}

func (s *SettingsService) GetSetting(ctx context.Context, userID string) (*db.Settings, error) {
	setting, err := s.db.GetSetting(ctx, userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return setting, nil
}
