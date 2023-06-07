package handler

import (
	"net/http"
	"serverwithpostgres/connectdb"

	_ "github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

var UserHandler = userHandler{"UserHandler"}

func ApiHandler(e *echo.Echo, userDB connectdb.UserDB) {
	err := userDB.ConnectDB()
	if err != nil {
		userDB.CloseDB()
	}
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	UserGroup := e.Group("/users")
	UserGroup.GET("/getusers", func(c echo.Context) error {
		// if c.QueryParams()
		err := UserHandler.getUsers(c, userDB)
		return err
	})
	UserGroup.GET("/getuser/:id",
		func(c echo.Context) error {
			err := UserHandler.getUserbyId(c, userDB)
			return err
		})
	
	UserGroup.POST("/createuser/", func(c echo.Context) error {
		err := UserHandler.createUser(c, userDB)
		return err
	})
	UserGroup.PUT("/updateuser/:id", func(c echo.Context) error {
		err := UserHandler.updateUser(c, userDB)
		return err
	})
	UserGroup.DELETE("/deleteuser/:id", func(c echo.Context) error {
		err := UserHandler.deleteUser(c, userDB)
		return err
	})

}
