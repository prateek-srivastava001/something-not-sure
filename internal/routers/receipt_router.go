package routers

import (
	"EasySplit/internal/controllers"

	"github.com/labstack/echo/v4"
)

func SetupReceiptRoutes(e *echo.Echo) {
	e.POST("/upload", controllers.UploadImage)
}
