package router

import (
	"github.com/gin-gonic/gin"
	"github.com/renatopnasc/made2share-api/internal/middleware"
)

func Initialize() {
	r := gin.Default()
	r.Use(middleware.Cors())

	initializeRoutes(r)

	r.Run(":8080")
}
