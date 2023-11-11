package services

import "caravagio-api-golang/internal/app/db"

type UserRoleService struct {
	db db.DBUserRoleRepo
}

func NewUserRoleService(db db.DBUserRoleRepo) *UserRoleService {
	return &UserRoleService{db: db}
}

func (s *UserRoleService) GetUserRole(userID string) (string, error) {

	role, err := s.db.GetUserRole(userID)

	if err != nil {
		return "", err
	}

	return role, nil
}
