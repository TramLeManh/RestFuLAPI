package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Avatar      string    `json:"avatar"`
	PhoneNumber string    `json:"phone_number"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Country     string    `json:"country"`
	City        string    `json:"city"`
	StreetName  string    `json:"street_name"`
	StreetAddr  string    `json:"street_address"`
}

var db, _ = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})

func CreateUser(c *gin.Context) {
	var user User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		user, err = getRandomUser()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data or failed to generate random user"})
			return
		}
	}

	// Validate fields
	if user.Email == "" || user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and Username are required"})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func getRandomUser() (User, error) {
	resp, err := http.Get("https://random-data-api.com/api/v2/users")
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()
	var randUser struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		Avatar      string `json:"avatar"`
		PhoneNumber string `json:"phone_number"`
		DateOfBirth string `json:"date_of_birth"`
		Address     struct {
			Country       string `json:"country"`
			City          string `json:"city"`
			StreetName    string `json:"street_name"`
			StreetAddress string `json:"street_address"`
		} `json:"address"`
	}
	json.NewDecoder(resp.Body).Decode(&randUser)
	dob, _ := time.Parse("2006-01-02", randUser.DateOfBirth)
	user := User{
		FirstName:   randUser.FirstName,
		LastName:    randUser.LastName,
		Username:    randUser.Username,
		Email:       randUser.Email,
		Avatar:      randUser.Avatar,
		PhoneNumber: randUser.PhoneNumber,
		DateOfBirth: dob,
		Country:     randUser.Address.Country,
		City:        randUser.Address.City,
		StreetName:  randUser.Address.StreetName,
		StreetAddr:  randUser.Address.StreetAddress,
	}
	return user, nil
}
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Model(&user).Updates(input)
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
func ListUsers(c *gin.Context) {
	var users []User
	query := db.Model(&User{})
	// Filtering
	if v := c.Query("username"); v != "" {
		query = query.Where("username = ?", v)
	}
	if v := c.Query("first_name"); v != "" {
		query = query.Where("first_name = ?", v)
	}
	if v := c.Query("last_name"); v != "" {
		query = query.Where("last_name = ?", v)
	}
	// Sorting
	sort := c.Query("sort")
	if sort != "" {
		query = query.Order(sort)
	}
	query.Find(&users)
	c.JSON(http.StatusOK, users)
}
