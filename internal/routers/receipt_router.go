package routers

import (
	"EasySplit/internal/controllers"
	"EasySplit/internal/middleware"

	"github.com/labstack/echo/v4"
)

func SetupReceiptRoutes(e *echo.Echo) {
	e.POST("/upload", controllers.UploadImage, middleware.JWTMiddleware)
}
