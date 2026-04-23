package router

import (
	"time"

	"camsystem/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	
)

func InitializeRoutes(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // depois pode restringir
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))


	v1 := router.Group("/api/v1")

	 v1.POST("/save-camera", handlers.SaveCamera)
	 v1.GET("/cameras", handlers.GetCameras)
}