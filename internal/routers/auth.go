package routers

import (
	"EasySplit/internal/controllers"
	"EasySplit/internal/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.POST("/signup", controllers.CreateUser)
	e.POST("/login", controllers.Login)
	e.GET("/me", controllers.GetUser, middleware.JWTMiddleware)
	e.PATCH("/me", controllers.UpdateUser, middleware.JWTMiddleware)
}
