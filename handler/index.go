package handler

import (
	"net/http"

	_ "github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

var userHandler = UserHandler{"UserHandler"}

func ApiHandler(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	UserGroup := e.Group("/users")
	UserGroup.GET("/getusers", userHandler.getUsers)
	UserGroup.GET("/getuser/:id", userHandler.getUserbyId)
	UserGroup.GET("/getuser", userHandler.getUsersbyQuery)
	UserGroup.POST("/createuser/", userHandler.createUser)
}
