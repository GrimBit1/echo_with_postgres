package handler

import (
	"io"
	"net/http"
	"serverwithpostgres/connectdb"
	"serverwithpostgres/logic"
	"serverwithpostgres/model"
	"strconv"

	"github.com/labstack/echo/v4"
)

var Logic = logic.UserLogic{}

type userHandler struct {
	Name string
}

// Route 1 get All users
func (userHandler) getUsers(c echo.Context, userdb connectdb.UserDB) error {

	var id = c.QueryParams()
	// fmt.Println(id)
	// If user hasn't given any query then give all users
	if len(id) == 0 {
		Users, err := Logic.GetAllUsers(userdb.DBP)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.Error{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, Users)

	} else {
		Users, err := Logic.GetUsersbyQuery(id, userdb.DBP)
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
		}
		if len(Users) == 0 {
			return c.JSON(http.StatusNotFound, model.Error{Message: "Not Found Anything..."})

		}
		return c.JSON(http.StatusOK, Users)
	}
}

// Route 2 get User by id
func (userHandler) getUserbyId(c echo.Context, userdb connectdb.UserDB) error {

	num, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusNotFound, model.Error{Message: "id should be valid"})

	}
	user, err := Logic.GetUser(num, userdb.DBP)

	if err != nil {
		if err.Error() == "Not Found" {

			return c.JSON(http.StatusNotFound, model.Error{Message: err.Error()})

		} else {
			return c.JSON(http.StatusNotFound, model.Error{Message: err.Error()})

		}
	}
	return c.JSON(http.StatusOK, user)

}


// Route 3 create user
func (userHandler) createUser(c echo.Context, userdb connectdb.UserDB) error {

	data, err := io.ReadAll(c.Request().Body)

	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
	}

	message, err := Logic.CreateUser(data, userdb.DBP)

	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})

	}
	// Logic.CreateUser()

	return c.JSON(http.StatusOK, message)

}

// Route 4 Update user
func (userHandler) updateUser(c echo.Context, userdb connectdb.UserDB) error {
	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
	}
	num, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
	}
	message, err := Logic.UpdateUser(data, num, userdb.DBP)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, message)
}

// Route 5 Delete user
func (userHandler) deleteUser(c echo.Context, userdb connectdb.UserDB) error {
	num, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusNotFound, model.Error{Message: "id should be valid"})
	}
	message, err := Logic.DeleteUser(num, userdb.DBP)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.Error{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, message)
}
