package controllers

import (
	"EasySplit/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SendFriendRequest(c echo.Context) error {
	var requestBody struct {
		ReceiverEmail string `json:"receiver_email"`
	}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
			"status":  "fail",
		})
	}

	senderEmail := c.Get("user_email").(string)
	receiverEmail := requestBody.ReceiverEmail

	if senderEmail == "" || receiverEmail == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "sender_email and receiver_email are required",
			"status":  "fail",
		})
	}

	if senderEmail == receiverEmail {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Cannot send friend request to yourself",
			"status":  "fail",
		})
	}

	isFriend, err := services.IsFriend(senderEmail, receiverEmail)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}
	if isFriend {
		return c.JSON(http.StatusConflict, map[string]string{
			"message": "You are already friends",
			"status":  "fail",
		})
	}

	if err := services.CreateFriendRequest(senderEmail, receiverEmail); err != nil {
		if err.Error() == "friend request already sent" {
			return c.JSON(http.StatusConflict, map[string]string{
				"message": "Friend request already sent",
				"status":  "fail",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Friend request sent",
		"status":  "success",
	})
}

func ConfirmFriendRequest(c echo.Context) error {
	var requestBody struct {
		SenderEmail string `json:"sender_email"`
	}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
			"status":  "fail",
		})
	}

	senderEmail := requestBody.SenderEmail
	receiverEmail := c.Get("user_email").(string)

	if senderEmail == "" || receiverEmail == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "sender_email and receiver_email are required",
			"status":  "fail",
		})
	}

	if senderEmail == receiverEmail {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid friend request",
			"status":  "fail",
		})
	}

	if err := services.AcceptFriendRequest(senderEmail, receiverEmail); err != nil {
		if err.Error() == "friend request not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "Friend request not found",
				"status":  "fail",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Friend request accepted",
		"status":  "success",
	})
}

func GetAllFriends(c echo.Context) error {
	email := c.Get("user_email").(string)
	friends, err := services.GetFriendsByEmail(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"friends": friends,
		"status":  "success",
	})
}

func GetPendingFriendRequests(c echo.Context) error {
	email := c.Get("user_email").(string)

	requests, err := services.GetPendingFriendRequests(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"requests": requests,
		"status":   "success",
	})
}

func GetFriendProfile(c echo.Context) error {
	friendEmail := c.Param("email")
	userEmail := c.Get("user_email").(string)

	friend, err := services.GetFriendByEmail(userEmail, friendEmail)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Friend not found",
			"status":  "fail",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"friend": friend,
		"status": "success",
	})
}

func RemoveFriend(c echo.Context) error {
	friendEmail := c.Param("id")
	email := c.Get("user_email").(string)

	if err := services.RemoveFriend(email, friendEmail); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Friend removed",
		"status":  "success",
	})
}
