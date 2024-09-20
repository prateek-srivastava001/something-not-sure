package routers

import (
	"EasySplit/internal/controllers"
	"EasySplit/internal/middleware"

	"github.com/labstack/echo/v4"
)

func SetupFriendRoutes(e *echo.Echo) {
	e.POST("/friend/add", controllers.SendFriendRequest, middleware.JWTMiddleware)
	e.POST("/friend/confirm", controllers.ConfirmFriendRequest, middleware.JWTMiddleware)
	e.GET("/friend/all", controllers.GetAllFriends, middleware.JWTMiddleware)
	e.GET("/friend/:email", controllers.GetFriendProfile, middleware.JWTMiddleware)
	e.DELETE("/friend/:email", controllers.RemoveFriend, middleware.JWTMiddleware)
	e.GET("/friend/requests/pending", controllers.GetPendingFriendRequests, middleware.JWTMiddleware)
}
