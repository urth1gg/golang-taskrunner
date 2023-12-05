package services

import (
	"caravagio-api-golang/internal/app/db"
	"context"
	"fmt"
)

type IVariablesService interface {
	GetVariables(ctx context.Context, headingID string) (db.Variables, error)
}

type VariablesService struct {
	db db.VariablesRepo
}

func NewVariablesService(db *db.DBVariablesRepo) *VariablesService {
	return &VariablesService{db: db}
}

func (s *VariablesService) GetVariables(ctx context.Context, headingID string) (db.Variables, error) {
	// get first row from variables table based on headingID

	variables, err := s.db.GetVariables(ctx, headingID)

	if err != nil {
		fmt.Println(err)
		return db.Variables{}, err
	}

	return variables, nil
}

func (s *VariablesService) CreateVariables(ctx context.Context, variables *db.Variables) (*db.Variables, error) {
	variables, err := s.db.CreateVariables(ctx, variables)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return variables, nil
}

func (s *VariablesService) UpdateVariables(ctx context.Context, variables *db.Variables) (int, error) {
	rowsAffected, err := s.db.UpdateVariables(ctx, variables)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return rowsAffected, nil
}
