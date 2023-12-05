package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Variables struct {
	ID             int    `sql:"id" json:"id"`
	H1Title        string `sql:"h1_title" json:"h1_title"`
	H2Title        string `sql:"h2_title" json:"h2_title"`
	AllHeaders     string `sql:"all_headers" json:"all_headers"` // Text field, but consider using string in Go
	CurrentHeader  string `sql:"current_header" json:"current_header"`
	PreviousHeader string `sql:"previous_header" json:"previous_header"`
	NextHeader     string `sql:"next_header" json:"next_header"`
	Keywords       string `sql:"keywords" json:"keywords"` // Text field, but consider using string in Go
	ParentHeader   string `sql:"parent_header" json:"parent_header"`
	MoreInfo       string `sql:"more_info" json:"more_info"`             // Text field, but consider using string in Go
	AdditionalInfo string `sql:"additional_info" json:"additional_info"` // Text field, but consider using string in Go
	MaxLength      int    `sql:"max_length" json:"max_length"`
	HeadingID      string `sql:"heading_id" json:"heading_id"`
}

type VariablesRepo interface {
	GetVariables(ctx context.Context, headingID string) (Variables, error)
	CreateVariables(ctx context.Context, variables *Variables) (*Variables, error)
	UpdateVariables(ctx context.Context, variables *Variables) (int, error)
}

type DBVariablesRepo struct {
	db *sql.DB
}

func (s *DBVariablesRepo) GetVariables(ctx context.Context, headingID string) (Variables, error) {
	// get fist row from variables table

	var variable Variables

	err := s.db.QueryRowContext(ctx, "SELECT id, h1_title, h2_title, all_headers, current_header, previous_header, next_header, keywords, parent_header, more_info, additional_info, max_length, heading_id FROM variables WHERE heading_id = ?", headingID).Scan(
		&variable.ID,
		&variable.H1Title,
		&variable.H2Title,
		&variable.AllHeaders,
		&variable.CurrentHeader,
		&variable.PreviousHeader,
		&variable.NextHeader,
		&variable.Keywords,
		&variable.ParentHeader,
		&variable.MoreInfo,
		&variable.AdditionalInfo,
		&variable.MaxLength,
		&variable.HeadingID,
	)

	if err != nil {
		return variable, err
	}

	return variable, nil

}

func (s *DBVariablesRepo) UpdateVariables(ctx context.Context, variables *Variables) (int, error) {
	v := reflect.ValueOf(variables).Elem()
	t := v.Type()

	setClauses := []string{}
	args := []interface{}{}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		// if !field.IsValid() || field.IsZero() {
		// 	continue // Skip invalid fields or zero value fields
		// }

		dbField := t.Field(i).Tag.Get("sql")
		if dbField == "" {
			dbField = strings.ToLower(t.Field(i).Name) // Fallback to field name if sql tag is missing
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = ?", dbField))
		args = append(args, field.Interface())
	}

	if len(setClauses) == 0 {
		return 0, errors.New("no fields to update")
	}

	setClause := strings.Join(setClauses, ", ")
	query := fmt.Sprintf("UPDATE variables SET %s WHERE heading_id = ?", setClause)
	args = append(args, variables.HeadingID)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, errors.New("no rows updated")
	}

	return int(rowsAffected), nil
}

func (s *DBVariablesRepo) CreateVariables(ctx context.Context, variables *Variables) (*Variables, error) {
	// insert into variables table
	query := "INSERT INTO variables (h1_title, h2_title, all_headers, current_header, previous_header, next_header, keywords, parent_header, more_info, additional_info, max_length, heading_id) VALUES (?, ?,? ,? ,? ,? ,? ,? ,? ,? ,? ,?)"

	result, err := s.db.ExecContext(ctx, query,
		variables.H1Title,
		variables.H2Title,
		variables.AllHeaders,
		variables.CurrentHeader,
		variables.PreviousHeader,
		variables.NextHeader,
		variables.Keywords,
		variables.ParentHeader,
		variables.MoreInfo,
		variables.AdditionalInfo,
		variables.MaxLength,
		variables.HeadingID,
	)

	if err != nil {
		return variables, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return variables, err
	}

	variables.ID = int(id)

	return variables, nil
}

func NewDBVariablesRepo(db *sql.DB) *DBVariablesRepo {
	return &DBVariablesRepo{db: db}
}
