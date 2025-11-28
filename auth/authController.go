package auth

import (
	"kaskade_backend/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT secret 从环境变量读取
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// 🔹 生成 JWT Token
func CreateJWT(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "missing user info"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.(models.User).ID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // 1天后过期
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	} else {
		c.SetCookie("jwt", tokenString, 3600, "/", "localhost", true, true)
		c.Next()
	}
}

// 🔹 验证 JWT Token
func AuthRequired(c *gin.Context) {
	// 用authorization header的情况，这样解析token
	// authHeader := c.GetHeader("Authorization")
	// if authHeader == "" {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
	// 	c.Abort()
	// 	return
	// }

	// 通常格式为 "Bearer <token>"
	// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	// tokenString = strings.TrimSpace(tokenString)
	tokenString, err := c.Cookie("jwt")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "token extraction failed"})
		return
	}

	// 解析 token 不清楚某次gpt chat为什么会提供这种方式
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	// 确保签名算法是预期的
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, jwt.ErrTokenMalformed
	// 	}
	// 	return jwtSecret, nil
	// })

	// if err != nil || !token.Valid {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
	// 	c.Abort()
	// 	return
	// }
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return
	}
	// 提取用户ID
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID, ok := claims["userID"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}
		c.Set("userID", userID)
		c.Next()
	}

}
