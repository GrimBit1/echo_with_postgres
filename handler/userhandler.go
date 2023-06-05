package handler

import (
	"net/http"
	"serverwithpostgres/model"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Name string
}

var Users = []model.User{}

func (UserHandler) getUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, Users)
}
