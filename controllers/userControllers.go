package controllers

import (
	"kaskade_backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUsers(c *gin.Context, db *gorm.DB) {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context, db *gorm.DB) {
	var req models.RegisterRequest

	// 尝试解析请求体为 JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	var existingUser models.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "username exists"})
		return
	}
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "email exists"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}
	if err := db.Create(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func Login(c *gin.Context, db *gorm.DB) {
	var req models.LoginRequest
	var foundUser models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}
	result := db.Where("username = ?", req.Username).First(&foundUser)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "username does not exist"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(req.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
		return
	}
	c.Set("user", foundUser)
	c.Next()
}
func Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func DeleteUser(c *gin.Context, db *gorm.DB) {
	username := c.Param("username")
	if db.Where("username = ?", username).Delete(&models.User{}).RowsAffected != 1 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "username not found, deletion failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func UpdateUser(c *gin.Context, db *gorm.DB) {
	username := c.Param("username")

	var existingUser models.User
	if err := db.Where("username = ?", username).First(&existingUser).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updatedData struct {
		Email           string `json:"email"`
		Username        string `json:"username"`
		CurrentPassword string `json:"currentpassword"`
		NewPassword     string `json:"newpassword"`
	}

	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 只更新提供的字段
	if updatedData.Email != "" {
		existingUser.Email = updatedData.Email
	}
	if updatedData.Username != "" {
		existingUser.Username = updatedData.Username
	}
	if updatedData.CurrentPassword != "" && updatedData.NewPassword != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(updatedData.CurrentPassword)); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Current password is wrong"})
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(updatedData.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		existingUser.PasswordHash = string(hashed)
	}

	if err := db.Save(&existingUser).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
