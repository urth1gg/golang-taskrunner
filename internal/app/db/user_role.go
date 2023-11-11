package db

import "database/sql"

const (
	RoleNewUser    = "1"
	RoleSeniorUser = "2"
	RoleAdmin      = "3"
)

type DBUserRoleRepo struct {
	db *sql.DB
}

func NewDBUserRoleRepo(db *sql.DB) *DBUserRoleRepo {
	return &DBUserRoleRepo{db: db}
}

func (r *DBUserRoleRepo) GetUserRole(userID string) (string, error) {
	var role string
	query := "SELECT role FROM user_role WHERE user_id = ?"
	err := r.db.QueryRow(query, userID).Scan(&role)

	if err != nil {
		return "", err
	}

	return role, nil
}
