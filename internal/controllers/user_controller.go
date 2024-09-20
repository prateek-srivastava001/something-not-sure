package controllers

import (
	"EasySplit/internal/models"
	"EasySplit/internal/services"
	"database/sql"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(c echo.Context) error {
	userEmail := c.Get("user_email").(string)

	user, err := services.FindUserByEmail(userEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "user not found",
				"status":  "fail",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":   user,
		"status": "success",
	})
}

func UpdateUser(c echo.Context) error {
	var payload models.UpdateUserRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
			"status":  "fail",
		})
	}

	userEmail := c.Get("user_email").(string)

	user, err := services.FindUserByEmail(userEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "user not found",
				"status":  "fail",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	// Update fields if provided
	if payload.FirstName != "" {
		user.FirstName = payload.FirstName
	}
	if payload.LastName != "" {
		user.LastName = payload.LastName
	}
	if payload.Username != "" {
		user.Username = payload.Username
	}
	if payload.Email != "" {
		user.Email = strings.ToLower(payload.Email)
	}
	if payload.Phone != "" {
		user.Phone = payload.Phone
	}
	if payload.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
				"status":  "error",
			})
		}
		user.Password = string(hashedPassword)
	}

	if err := services.UpdateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "user updated successfully",
		"status":  "success",
	})
}
