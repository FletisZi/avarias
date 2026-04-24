package router

import (
	"time"

	"camsystem/handlers"
	"camsystem/schemas"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, manager *schemas.StreamManager) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // depois pode restringir
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := router.Group("/api/v1")

	v1.POST("/new-camera", handlers.SaveCamera)
	v1.POST("/get-camera", handlers.GetCamera(manager))
	v1.POST("/start-camera", handlers.StartCamera(manager))
	v1.POST("/stop-camera", handlers.StopCamera(manager))
	v1.GET("/cameras", handlers.GetCameras)
}
