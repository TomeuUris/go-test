package infrastructure

import (
	"github.com/TomeuUris/go-test/pkg/application"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *application.UserService
}

func NewUserHandler(userService *application.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	uuidStr := c.Params.ByName("uuid")
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}
	user, err := h.userService.GetUser(uuid)
	if err != nil {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, user)
	}
}

// Add other handlers for creating, updating, and deleting users...
