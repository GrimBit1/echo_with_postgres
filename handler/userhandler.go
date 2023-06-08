package handler

import (
	"io"
	"net/http"
	"serverwithpostgres/logic"
	"serverwithpostgres/model"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	Name      string
	db        *sqlx.DB
	userLogic logic.UserLogic
}

// Route 1 get All users
func (u *userHandler) getUsers(c echo.Context) error {

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
		Users, err := u.userLogic.GetAllUsers(pageNo, pageSize)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.Error{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, Users)

	} else {
		// If user has given the role in the query
		var queryStrForRole string
		if id.Has("role") {
			roleArr := strings.Split(id.Get("role"), ",")
			if len(roleArr) > 1 {
				return c.JSON(http.StatusBadRequest, model.Error{Message: "Only 1 value for role is available at this moment "})

			}
			if len(roleArr) == 1 {
				queryStrForRole = "role::varchar ilike '%" + roleArr[0] + "%'"
				id.Del("role")
			}
		}

		// If user has given other parameters in the query
		var queryStr string
		if len(id) != 0 {
			if len(queryStrForRole) != 0 {
				queryStrForRole = " AND " + queryStrForRole
			}
			querys := logic.UrlValuesToString(id, " iLike ", "'", "%")
			queryStr = logic.JoinArray(querys, " AND ")
		}
		Users, err := u.userLogic.GetUsersbyQuery(queryStr, queryStrForRole)
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

	num, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusNotFound, model.Error{Message: "id should be valid"})

	}
	user, err := u.userLogic.GetUser(num)

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

	data, err := io.ReadAll(c.Request().Body)

	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
	}

	message, err := u.userLogic.CreateUser(data)

	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})

	}

	// u.userLogic.CreateUser()

	return c.JSON(http.StatusOK, message)

}

// Route 4 Update user
func (u *userHandler) updateUser(c echo.Context) error {

	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
	}
	num, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
	}
	message, err := u.userLogic.UpdateUser(data, num)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, message)
}

// Route 5 Delete user
func (u *userHandler) deleteUser(c echo.Context) error {

	num, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusNotFound, model.Error{Message: "id should be valid"})
	}
	message, err := u.userLogic.DeleteUser(num)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.Error{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, message)
}
