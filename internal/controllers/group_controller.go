package controllers

import (
	"EasySplit/internal/services"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func CreateGroup(c echo.Context) error {
	var requestBody struct {
		Name string `json:"name"`
	}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
			"status":  "fail",
		})
	}

	creatorEmail := c.Get("user_email").(string)
	name := requestBody.Name

	groupID, err := services.CreateGroup(creatorEmail, name)
	if err != nil {
		if err.Error() == "group with this name already exists" {
			return c.JSON(http.StatusConflict, map[string]string{
				"message": "Group with this name already exists",
				"status":  "fail",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"group_id": groupID,
		"status":   "success",
	})
}

func AddUserToGroup(c echo.Context) error {
	var requestBody struct {
		GroupID   string `json:"group_id"`
		UserEmail string `json:"user_email"`
	}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
			"status":  "fail",
		})
	}

	groupID, err := uuid.Parse(requestBody.GroupID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid group ID",
			"status":  "fail",
		})
	}

	creatorEmail := c.Get("user_email").(string)
	userEmail := requestBody.UserEmail

	if err := services.AddUserToGroup(groupID, userEmail, creatorEmail); err != nil {
		if err.Error() == "user already in group" {
			return c.JSON(http.StatusConflict, map[string]string{
				"message": "User already in group",
				"status":  "fail",
			})
		}
		if err.Error() == "group does not exist" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "Group not found",
				"status":  "fail",
			})
		}
		if err.Error() == "user does not exist" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "User not found",
				"status":  "fail",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User added to group",
		"status":  "success",
	})
}

func GetGroupByID(c echo.Context) error {
	groupID := c.Param("id")
	parsedGroupID, err := uuid.Parse(groupID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid group ID", "status": "fail",
		})
	}

	group, err := services.GetGroupByID(parsedGroupID)
	if err != nil {
		if err.Error() == "group not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "Group not found",
				"status":  "fail",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"group":  group,
		"status": "success",
	})
}

func GetAllGroups(c echo.Context) error {
	userEmail := c.Get("user_email").(string)
	groups, err := services.GetAllGroups(userEmail)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"groups": groups,
		"status": "success",
	})
}

func DeleteGroup(c echo.Context) error {
	groupID := c.Param("id")
	parsedGroupID, err := uuid.Parse(groupID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid group ID", "status": "fail",
		})
	}

	creatorEmail := c.Get("user_email").(string)

	if err := services.DeleteGroup(parsedGroupID, creatorEmail); err != nil {
		if err.Error() == "user is not the creator of the group" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"message": "User is not the creator of the group",
				"status":  "fail",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Group deleted successfully",
		"status":  "success",
	})
}
