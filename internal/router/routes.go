package router

import (
	"github.com/gin-gonic/gin"
	"github.com/renatopnasc/made2share-api/internal/handler"
	"github.com/renatopnasc/made2share-api/internal/middleware"
)

const (
	baseURL = "/api/v1"
)

func initializeRoutes(r *gin.Engine) {

	v1 := r.Group(baseURL)
	{
		v1.GET("/login", handler.LoginHandler)
		v1.GET("/callback", handler.CallbackHandler)

		authRoutes := v1.Group("/")
		authRoutes.Use(middleware.VerifyAuthentication())
		{
			authRoutes.GET("/me", handler.MeHandler)
			authRoutes.POST("/playlists", handler.CreatePlaylistHandler)
		}
	}
}
