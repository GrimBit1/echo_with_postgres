package handler

import (
	"fmt"
	"io"
	"net/http"
	checkerror "serverwithpostgres/checkError"
	"serverwithpostgres/logic"
	"serverwithpostgres/model"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Name string
}

func (UserHandler) getUsers(c echo.Context) error {
	logic.GetAllUsers()
	return c.JSON(http.StatusOK, logic.Users)
}

func (UserHandler) getUserbyId(c echo.Context) error {
	num, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.Error{Message: "id should be valid"})
	}
	user, err := logic.GetUser(num)
	if err != nil {
		if err.Error() == "Not Found" {
			return c.JSON(http.StatusNotFound, model.Error{Message: err.Error()})

		} else {
			return c.JSON(http.StatusNotFound, model.Error{Message: err.Error()})
		}
	}
	return c.JSON(http.StatusOK, user)

}

func (UserHandler) getUsersbyQuery(c echo.Context) error {
	var id = c.QueryParams()
	// var queryStr string
	// if id.Has() {

	// }
	fmt.Println(id.Has("hi"))
	fmt.Println(id)
	return c.JSON(http.StatusOK, id)
}
func (UserHandler) createUser(c echo.Context) error {
	// var (
	// 	f_name string
	// 	l_name string
	// 	role []string
	// 	title string
	// )
	data, err := io.ReadAll(c.Request().Body)
	checkerror.CheckError(err)
	message, err1 := logic.CreateUser(data)
	if err1 != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err1.Error()})
	}
	// logic.CreateUser()
	return c.JSON(http.StatusOK, message)
}
