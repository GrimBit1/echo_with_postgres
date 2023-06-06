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

// Route 1 get All users
func (UserHandler) getUsers(c echo.Context) error {

	logic.GetAllUsers()

	return c.JSON(http.StatusOK, logic.Users)

}

// Route 2 get User by id
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

// Route 3 get user by query
func (UserHandler) getUsersbyQuery(c echo.Context) error {

	var id = c.QueryParams()

	// var queryStr string
	// if id.Has() {
	// }
	fmt.Println(len(id))
	
	for i := range id {
		fmt.Println(i, ":", id[i][0])
	}
	logic.GetUsersbyQuery()
	return c.JSON(http.StatusOK, id)

}

// Route 4 create user
func (UserHandler) createUser(c echo.Context) error {

	data, err := io.ReadAll(c.Request().Body)

	checkerror.CheckError(err)

	message, err1 := logic.CreateUser(data)

	if err1 != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err1.Error()})

	}
	// logic.CreateUser()

	return c.JSON(http.StatusOK, message)

}

// Route 5 Update user
func (UserHandler) updateUser(c echo.Context) error {
	data, err := io.ReadAll(c.Request().Body)
	num, err1 := strconv.ParseInt(c.Param("id"), 10, 64)
	checkerror.CheckError(err, err1)
	fmt.Println(num, string(data))
	message, err2 := logic.UpdateUser(data, num)
	checkerror.CheckError(err2)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, message)
}

// Route 6 Delete user
func (UserHandler) deleteUser(c echo.Context) error {
	num, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusNotFound, model.Error{Message: "id should be valid"})
	}
	message, err := logic.DeleteUser(num)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.Error{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, message)

}
