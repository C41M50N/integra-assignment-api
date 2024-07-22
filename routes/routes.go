package routes

import (
	"fmt"
	"integra-api/database"
	"integra-api/models"
	"integra-api/utils"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetAllUsers(c echo.Context) error {
	users := []models.User{}
	database.DB.Find(&users)
	return c.JSON(http.StatusOK, users)
}

func CreateUser(c echo.Context) error {
	user := new(models.CreateUserInput)
	if err := c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	log.Printf("new user created: %#v", user)

	user_name := fmt.Sprintf("%s%s%s", strings.ToLower(user.FirstName), strings.ToLower(user.LastName), utils.GenerateUserNameSuffix())
	email := fmt.Sprintf("%s@integra.com", user_name)

	database.DB.Create(&models.User{
		UserName:   user_name,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      email,
		UserStatus: "A",
		Department: user.Department,
	})

	return c.String(http.StatusOK, "created")
}

func UpdateUser(c echo.Context) error {
	rawId := c.Param("id")
	id, err := strconv.Atoi(rawId)
	if err != nil && id > 0 {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user := new(models.UpdateUserInput)
	if err := c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	log.Printf("user to update: %#v", user)

	database.DB.Exec("UPDATE users SET first_name = ?, last_name = ?, email = ?, user_status = ?, department = ? WHERE id = ?", user.FirstName, user.LastName, user.Email, user.UserStatus, user.Department, user.Id)

	return c.String(http.StatusOK, "updated")
}

func DeleteUser(c echo.Context) error {
	rawId := c.Param("id")
	id, err := strconv.Atoi(rawId)
	if err != nil && id > 0 {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user := models.User{Id: id}
	database.DB.Delete(&user)

	return c.String(http.StatusOK, "deleted")
}
