package application

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/TomeuUris/go-test/pkg/domain"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUser(uuid uuid.UUID) (domain.UserResponse, error) {
	var user domain.User
	if err := s.db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return domain.UserResponse{}, err
	}
	return user.ToUserResponse(), nil
}
