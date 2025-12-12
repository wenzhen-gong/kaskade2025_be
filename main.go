package main

import (
	"kaskade_backend/db"
	"kaskade_backend/routes"

	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 尝试加载 .env（本地开发用），Cloud Run 会忽略这个错误
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Cloud Run 默认端口
	}
	// Connect to db
	DB := db.Init() // 初始化数据库 + 自动迁移
	// Create router using gin
	r := gin.Default()
	// ✅ 启用 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	routes.SetupRoutes(r, DB)
	// Start server
	r.Run(":" + port)
	log.Println("Server is running")
}
