package handler

import (
	"net/http"
	"serverwithpostgres/connectdb"
	"serverwithpostgres/logic"

	_ "github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func ApiHandler(e *echo.Echo, udb connectdb.UserDB) {
	var UserHandler = userHandler{"UserHandler", logic.UserLogic{DB: udb.DBP,RoleMap: udb.RoleMap}}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	UserGroup := e.Group("/users")
	UserGroup.GET("/getusers", UserHandler.getUsers)
	UserGroup.GET("/getuser/:id", UserHandler.getUserbyId)
	UserGroup.POST("/createuser/", UserHandler.createUser)
	UserGroup.PUT("/updateuser/:id", UserHandler.updateUser)
	UserGroup.DELETE("/deleteuser/:id", UserHandler.deleteUser)
}
