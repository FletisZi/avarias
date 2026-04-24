package router

import (
	"camsystem/schemas"

	"github.com/gin-gonic/gin"
)

func Initialize(manager *schemas.StreamManager) {
	router := gin.Default()

	InitializeRoutes(router, manager)

	router.Run(":8080")
}
