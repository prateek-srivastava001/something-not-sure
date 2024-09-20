package routers

import (
	"EasySplit/internal/controllers"
	"EasySplit/internal/middleware"

	"github.com/labstack/echo/v4"
)

func SetupGroupRoutes(e *echo.Echo) {
	e.PUT("/group", controllers.CreateGroup, middleware.JWTMiddleware)
	e.POST("/group/:id/user", controllers.AddUserToGroup, middleware.JWTMiddleware)
	e.GET("/group/:id", controllers.GetGroupByID, middleware.JWTMiddleware)
	e.GET("/group/all", controllers.GetAllGroups, middleware.JWTMiddleware)
	e.DELETE("/group/:id", controllers.DeleteGroup, middleware.JWTMiddleware)
}
