package controllers

import (
	"EasySplit/internal/models"
	"EasySplit/internal/services"
	"EasySplit/internal/utils"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(ctx echo.Context) error {
	var payload models.LoginRequest

	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
			"status":  "fail",
		})
	}

	payload.Email = strings.ToLower(payload.Email)

	user, err := services.FindUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"message": "user does not exist",
				"status":  "fail",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return ctx.JSON(http.StatusConflict, map[string]string{
			"message": "Invalid password",
			"status":  "fail",
		})
	}

	accessToken, err := utils.CreateToken(utils.TokenPayload{
		Exp:   time.Minute * 500,
		Email: user.Email,
		Role:  user.Role,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	ctx.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400,
		Secure:   true,
	})

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":     "login successful",
		"accessToken": accessToken,
		"status":      "success",
	})
}

func CreateUser(ctx echo.Context) error {
	var payload models.CreateUserRequest

	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
			"status":  "fail",
		})
	}

	payload.Email = strings.ToLower(payload.Email)

	_, err := services.FindUserByEmail(payload.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	} else if err == nil {
		return ctx.JSON(http.StatusConflict, map[string]string{
			"message": "user already exists",
			"status":  "error",
		})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	user := models.NewUser(payload.FirstName, payload.LastName, payload.Username, payload.Email, payload.Phone, string(hashed), "user")

	if err := services.InsertUser(user); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "user creation was successful",
		"status":  "success",
	})
}
