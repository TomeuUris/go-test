package main

import (
	"net/http"
	"os"

	_ "github.com/TomeuUris/go-test/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID    uint      `json:"-"` // Hide ID field from JSON
	UUID  uuid.UUID `gorm:"type:char(36);unique_index" json:"uuid"`
	Name  string    `json:"name" binding:"required"`
	Email string    `json:"email" gorm:"type:varchar(100);unique_index" binding:"required"`
}

var db *gorm.DB

// @Summary List users
// @Description get all Users
// @Accept  json
// @Produce  json
// @Success 200 {array} User
// @Router /users [get]
func GetUsers(c *gin.Context) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, users)
	}
}

// @Summary Create user
// @Description Create a new User
// @Accept  json
// @Produce  json
// @Param   user     body    User     true        "User info"
// @Success 200 {object} User
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.UUID = uuid.New() // Generate a new UUID for the user

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Get user
// @Description get User by UUID
// @Accept  json
// @Produce  json
// @Param   uuid     path    string     true        "User UUID"
// @Success 200 {object} User
// @Router /users/{uuid} [get]
func GetUser(c *gin.Context) {
	var user User
	uuid := c.Params.ByName("uuid")
	if err := db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, user)
	}
}

// @Summary Update user
// @Description Update an existing User
// @Accept  json
// @Produce  json
// @Param   id     path    int     true        "User ID"
// @Param   name     body    string     true        "User name"
// @Param   email    body    string     true        "User email"
// @Success 200 {object} User
// @Router /users/{id} [patch]
func UpdateUser(c *gin.Context) {
	var user User
	uuid := c.Params.ByName("uuid")
	if err := db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		c.AbortWithStatus(404)
	} else {
		c.BindJSON(&user)
		db.Save(&user)
		c.JSON(200, user)
	}
}

// @Summary Delete user
// @Description Delete an existing User
// @Accept  json
// @Produce  json
// @Param   id     path    int     true        "User ID"
// @Success 200 {object} map[string]string
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	uuid := c.Params.ByName("uuid")
	var user User
	d := db.Where("uuid = ?", uuid).Delete(&user)
	if d.Error != nil {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, gin.H{"id " + uuid: "is deleted"})
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/users", GetUsers)
	r.GET("/users/:id", GetUser)
	r.POST("/users", CreateUser)
	r.PATCH("/users/:id", UpdateUser)
	r.DELETE("/users/:id", DeleteUser)

	if os.Getenv("ENV") != "prod" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return r
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server for a swagger example.
// @host localhost:8080
// @BasePath /
func main() {
	// NOTE: replace 'test.db' with your actual db path
	var err error
	db, err = gorm.Open(sqlite.Open("/home/appuser/test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	defer func() {
		dbSQL, err := db.DB()
		if err != nil {
			panic("failed to get database connection")
		}
		dbSQL.Close()
	}()

	db.AutoMigrate(&User{})

	db.Migrator().AlterColumn(&User{}, "ID")

	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

// func main() {
//     db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
//     db.AutoMigrate(&User{})

//     userService := NewUserService(db)
//     userHandler := NewUserHandler(userService)

//     r := gin.Default()
//     r.GET("/users/:uuid", userHandler.GetUser)
//     // Register other handlers...

//     r.Run()
// }
