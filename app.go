package main

import (
	"integra-api/database"
	"integra-api/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database.ConnectDB()

	app := echo.New()

	app.GET("/users", routes.GetAllUsers)
	app.POST("/users", routes.CreateUser)
	app.PUT("/users/:id", routes.UpdateUser)
	app.DELETE("/users/:id", routes.DeleteUser)

	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	app.Logger.Fatal(app.Start(":8081"))
}
