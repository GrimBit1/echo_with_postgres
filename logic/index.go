package logic

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	// "fmt"

	"serverwithpostgres/model"

	"github.com/jmoiron/sqlx"
)

type UserLogic struct {
	DB      *sqlx.DB
	RoleMap map[int]string
}

// Route 1 Done Get All Users
func (u *UserLogic) GetAllUsers(pageNo int64, pageSize int64) ([]model.User, error) {
	var pageQuery string
	if pageSize == 0 {
		pageSize = 10
	}
	if pageNo > 0 {
		pageQuery = `Limit ` + strconv.FormatInt(pageSize, 10) + ` OFFSET ` + strconv.FormatInt((pageNo*pageSize)-pageSize, 10)
	}
	mainQuery := "Select * from userswithjob " + " Order by id asc " + pageQuery
	// fmt.Println(mainQuery)
	res, err := u.DB.Query(mainQuery)

	if err != nil {
		return []model.User{}, err
	}

	var Users = []model.User{}
	for res.Next() {

		var newUser, _ = u.makeUserFromDB(nil, res)

		Users = append(Users, newUser)

	}
	return Users, nil
}

func (u *UserLogic) GetUser(id int64) (model.User, error) {

	res := u.DB.QueryRow(`select * from userswithjob where id =$1`, id)

	// user, err := makeUserFromDB(res, nil)

	user, err := u.makeUserFromDB(res, nil)

	if err != nil {
		if err.Error() == `sql: no rows in result set` {

			return model.User{}, errors.New("not found")

		} else {
			return model.User{}, err
		}
	}
	return user, nil
}

// Function to create user
func (u *UserLogic) CreateUser(data []byte) (model.Error, error) {

	var user = model.User{}

	err := json.Unmarshal(data, &user)
	if err != nil {
		return model.Error{}, err
	}
	// fmt.Printf("%#v\n", user)

	if len(user.FirstName) == 0 {

		return model.Error{}, errors.New("firstname and lastname cannot be empty")

	}
	if len(user.LastName) == 0 {

		return model.Error{}, errors.New("firstname and lastname cannot be empty")

	}
	temprole, err := json.Marshal(user.Role)

	if err != nil {
		return model.Error{}, err
	}

	// fmt.Println(string(temprole))

	res := u.DB.MustExec("INSERT INTO userswithjob(first_name,last_name,role,title) VALUES ($1,$2,$3,$4)", user.FirstName, user.LastName, temprole, user.JobTitle)

	rowChanged, err := res.RowsAffected()

	if err != nil {
		return model.Error{}, err
	}

	// fmt.Println("id", rowChanged)

	if rowChanged > 0 {

		return model.Error{Message: "User Created Successfully"}, nil
	}
	// var id int64

	// fmt.Println(newUser)
	return model.Error{}, errors.New("Something went wrong")

}

// Function to update user in db
func (u *UserLogic) UpdateUser(data []byte, id int64) (model.Error, error) {
	// Check the db that the id exists or not
	var res = u.DB.QueryRow("Select * from userswithjob where id = $1", id)
	oldUser, err := u.makeUserFromDB(res, nil)
	if err != nil {
		if err.Error() == `sql: no rows in result set` {

			return model.Error{}, errors.New("Not Found")

		} else {
			return model.Error{}, err
		}
	}

	//Create user template to update old one
	var newUser = model.UpdateUser{}

	err = json.Unmarshal(data, &newUser)
	// fmt.Println(newUser)
	if err != nil {
		return model.Error{}, errors.New("Invalid Json Values")
	}

	//Check if newUser has provided the firstname and lastname
	if len(newUser.FirstName) != 0 {

		oldUser.FirstName = newUser.FirstName

	}
	if len(newUser.LastName) != 0 {

		oldUser.LastName = newUser.LastName
	}
	// fmt.Println(newUser.AddRole)
	if newUser.AddRole != 0 {
		index := GiveIndex(oldUser.Role, u.RoleMap[newUser.AddRole])
		fmt.Println(index)
		if index > -1 {
			return model.Error{}, errors.New("The Given role is already available in the user")
		}
		oldUser.Role = append(oldUser.Role, u.RoleMap[newUser.AddRole])
	}
	if newUser.RemoveRole != 0 {
		index := GiveIndex(oldUser.Role, u.RoleMap[newUser.AddRole])
		if index < 0 {
			return model.Error{}, errors.New("The Given role is not available in the users")
		}
		// fmt.Println(index)
		oldUser.Role = RemoveIndex(oldUser.Role, index)
	}
	if len(newUser.JobTitle) != 0 {
		oldUser.JobTitle = newUser.JobTitle
	}
	// fmt.Println(oldUser)
	// Update the user in db
	idRole := u.GiveIdFromRole(oldUser.Role)
	fmt.Println(idRole)
	var temprole, err1 = json.Marshal(idRole)
	if err1 != nil {
		return model.Error{}, errors.New("Invalid Json Values")
	}

	result := u.DB.MustExec(`Update userswithjob set first_name = $1 ,last_name=$2 ,role=$3 ,title=$4 where id = $5`, oldUser.FirstName, oldUser.LastName, temprole, oldUser.JobTitle, id)
	var rowChanged, err2 = result.RowsAffected()
	if err2 != nil {
		return model.Error{}, err
	}
	if rowChanged > 0 {

		return model.Error{Message: "Updated Successfully"}, nil
	}
	return model.Error{}, errors.New("Some Error Occured")

}

// Function to delete user from db
func (u *UserLogic) DeleteUser(id int64) (model.Error, error) {
	res := u.DB.MustExec(`Delete from userswithjob where id = $1`, id)
	// fmt.Println(res)
	rowChanged, err := res.RowsAffected()

	if err != nil {
		return model.Error{}, errors.New("Something went wrong")

	}
	if rowChanged > 0 {

		return model.Error{Message: "User Deleted Successfully"}, nil
	}
	return model.Error{}, errors.New("Bad Request")

}

func (u *UserLogic) GetUsersbyQuery(roleArr []int, queryStrForRole string, id url.Values) ([]model.User, error) {

	// fmt.Println(queryStr)
	// Making a main query string using all the query parameter
	// fmt.Println(id)
	if len(roleArr) == 1 {
		queryStrForRole = "role::varchar ilike '%" + strconv.Itoa(roleArr[0]) + "%'"
	}
	var queryStr string
	if len(id) > 0 {
		if len(queryStrForRole) != 0 {
			queryStrForRole = " AND " + queryStrForRole
		}
		querys := UrlValuesToString(id, " iLike ", "'", "%")
		queryStr = JoinArray(querys, " AND ")
	}
	queryStrMain := "Select * from userswithjob where " + queryStr + queryStrForRole + " order by id asc"
	// fmt.Println(queryStrMain)
	res, err := u.DB.Query(queryStrMain)

	// If got error from the db then push it as error
	if err != nil {
		return nil, errors.New("Query Parameter is wrong" + strings.Split(err.Error(), ":")[1])
	}
	var Users = []model.User{}

	for res.Next() {
		var newUser, err = u.makeUserFromDB(nil, res)
		if err != nil {
			return []model.User{}, err
		}
		Users = append(Users, newUser)
	}
	if err != nil {
		return []model.User{}, err
	}

	return Users, nil
}

// Function to make user from db
func (u *UserLogic) makeUserFromDB(res *sql.Row, ress *sql.Rows) (model.User, error) {

	if res != nil {
		var (
			id       int64
			f_name   string
			l_name   string
			temprole string
			title    string
		)

		err := res.Scan(&id,
			&f_name,
			&l_name,
			&temprole,
			&title)
		// fmt.Println(id,
		// f_name,
		// l_name,
		// temprole,
		// title)

		if err != nil {
			return model.User{}, err
		}
		var role []int
		err1 := json.Unmarshal([]byte(temprole), &role)

		if err1 != nil {
			return model.User{}, err1
		}

		roleArr := u.GiveRoleFromId(role)
		var newUser = model.User{ID: id, FirstName: f_name, LastName: l_name, Role: roleArr, JobTitle: title}
		return newUser, nil
	} else {
		var (
			id       int64
			f_name   string
			l_name   string
			temprole string
			title    string
		)

		err := (ress.Scan(&id,
			&f_name,
			&l_name,
			&temprole,
			&title))

		if err != nil {
			return model.User{}, err
		}

		// fmt.Println(id,
		// 	f_name,
		// 	l_name,
		// 	temprole,
		// 	title)

		var role []int
		err = json.Unmarshal([]byte(temprole), &role)
		if err != nil {
			return model.User{}, err
		}
		roleArr := u.GiveRoleFromId(role)

		var newUser = model.User{ID: id, FirstName: f_name, LastName: l_name, Role: roleArr, JobTitle: title}
		return newUser, nil
	}

}

// First parameter is for slice and second is for join string
func JoinArray(slc []string, str string) string {
	var temp string = slc[0]

	for i := 1; i < len(slc); i++ {
		v := slc[i]
		temp += str + v
	}
	return temp
}

// First parameter is for slice and second is for join map values and third if we want to give some focus with additional 4th and 5th string  on map values
func UrlValuesToString(slc url.Values, str string, str2 string, str3 string) []string {
	var temp []string
	for i := range slc {
		tempstr := i + str + string(str2+slc[i][0]+str3+str2)
		temp = append(temp, tempstr)
	}
	return temp
}

func GiveIndex(slc []string, integer string) int {
	for i := range slc {
		// fmt.Println(i)
		if slc[i] == integer {
			return i
		}
	}
	return -1
}
func GiveKey(slc map[int]string, integer string) int {
	for i := range slc {
		// fmt.Println(i)
		if slc[i] == integer {
			return i
		}
	}
	return -1
}
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
func StringArrtoIntArr(slc []string) ([]int, error) {
	tempArr := []int{}
	for _, v := range slc {
		strtoint, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		tempArr = append(tempArr, strtoint)
	}
	return tempArr, nil
}

func (u *UserLogic) GiveRoleFromId(intarr []int) []string {
	roleArr := []string{}
	for _, v := range intarr {
		roleArr = append(roleArr, u.RoleMap[v])
	}
	return roleArr
}
func (u *UserLogic) GiveIdFromRole(intarr []string) []int {
	idArr := []int{}
	for _, v := range intarr {
		index := GiveKey(u.RoleMap, v)
		idArr = append(idArr, index)
	}
	return idArr
}
