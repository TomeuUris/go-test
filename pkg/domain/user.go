package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID    uint      `json:"-"` // Hide ID field from JSON
	UUID  uuid.UUID `gorm:"type:char(36);unique_index" json:"uuid"`
	Name  string    `json:"name" binding:"required"`
	Email string    `json:"email" gorm:"type:varchar(100);unique_index" binding:"required"`
}

type UserResponse struct {
	UUID  uuid.UUID `json:"uuid"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

// Convert User to UserResponse
func (u User) ToUserResponse() UserResponse {
	return UserResponse{
		UUID:  u.UUID,
		Name:  u.Name,
		Email: u.Email,
	}
}
