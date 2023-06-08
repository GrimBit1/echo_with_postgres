package handler

import (
	"io"
	"net/http"
	"serverwithpostgres/logic"
	"serverwithpostgres/model"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	Name string
	db   *sqlx.DB
}

// Route 1 get All users
func (u *userHandler) getUsers(c echo.Context) error {
	var Logic = logic.UserLogic{DB: u.db}
	var id = c.QueryParams()
	var pageNo int64
	var err error
	var pageSize int64

	//If user has given pageSize query
	if id.Has("pageSize") {

		//If user has given pageno query if not then
		if !id.Has("page") {
			pageNo = 1
		}
		pageSize, err = strconv.ParseInt(id.Get("pageSize"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.Error{Message: "pageSize parameter is not valid"})
		}
		id.Del("pageSize")
	}
	if id.Has("page") {
		pageNo, err = strconv.ParseInt(id.Get("page"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.Error{Message: "page parameter is not valid"})
		}
		id.Del("page")
	}
	// fmt.Println(id)

	// If user hasn't given any query then give all users
	if len(id) == 0 {
		Users, err := Logic.GetAllUsers(pageNo, pageSize)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.Error{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, Users)

	} else {
		Users, err := Logic.GetUsersbyQuery(id)
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
func (u *userHandler) getUserbyId(c echo.Context) error {
	var Logic = logic.UserLogic{DB: u.db}

	num, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusNotFound, model.Error{Message: "id should be valid"})

	}
	user, err := Logic.GetUser(num)

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
func (u *userHandler) createUser(c echo.Context) error {
	var Logic = logic.UserLogic{DB: u.db}

	data, err := io.ReadAll(c.Request().Body)

	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
	}

	message, err := Logic.CreateUser(data)

	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})

	}

	// Logic.CreateUser()

	return c.JSON(http.StatusOK, message)

}

// Route 4 Update user
func (u *userHandler) updateUser(c echo.Context) error {
	var Logic = logic.UserLogic{DB: u.db}

	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
	}
	num, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
	}
	message, err := Logic.UpdateUser(data, num)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, message)
}

// Route 5 Delete user
func (u *userHandler) deleteUser(c echo.Context) error {
	var Logic = logic.UserLogic{DB: u.db}

	num, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusNotFound, model.Error{Message: "id should be valid"})
	}
	message, err := Logic.DeleteUser(num)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.Error{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, message)
}
