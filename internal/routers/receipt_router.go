package routers

import (
	"EasySplit/internal/controllers"
	"EasySplit/internal/middleware"

	"github.com/labstack/echo/v4"
)

func SetupReceiptRoutes(e *echo.Echo) {
	e.POST("/upload/image", controllers.UploadImage, middleware.JWTMiddleware)
	e.POST("/upload/audio", controllers.UploadAudio, middleware.JWTMiddleware)
}
