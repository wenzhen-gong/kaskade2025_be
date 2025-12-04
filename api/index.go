package handler

import (
	"kaskade_backend/db"
	"kaskade_backend/routes"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

// init() 只在 serverless 冷启动执行一次
func init() {
	// 初始化 DB
	DB := db.Init()

	// Setup Gin
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // ⚠️ production 要改成你的前端 domain
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// 注册路由
	routes.SetupRoutes(r, DB)

	router = r
}

// Vercel 会调用这个
func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
