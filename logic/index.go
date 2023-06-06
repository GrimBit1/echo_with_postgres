package logic

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	checkerror "serverwithpostgres/checkError"
	"serverwithpostgres/connectdb"
	"serverwithpostgres/model"
)

var Users = []model.User{}
var UserPointer = &Users

// Route 1 Done Get All Users
func GetAllUsers() {
	res, err := connectdb.GiveDB().Query("Select * from userswithjob")
	checkerror.CheckError(err)
	var tempArr = []model.User{}
	for res.Next() {

		var newUser, err = makeUserFromDB(nil, res)
		checkerror.CheckError(err)
		tempArr = append(tempArr, newUser)
	}
	*UserPointer = tempArr
}

func GetUser(id int64) (model.User, error) {
	res := connectdb.GiveDB().QueryRow(`select * from userswithjob where id =$1`, id)
	// user, err := makeUserFromDB(res, nil)
	user, err := makeUserFromDB(res, nil)
	if err != nil {
		if err.Error() == `sql: no rows in result set` {
			return model.User{}, errors.New("not found")
		}else{
			return model.User{}, err
		}
	}
	return user, nil
}

// Function to create user
func CreateUser(data []byte) (model.Error, error) {
	var user = model.User{}
	checkerror.CheckError(json.Unmarshal(data, &user))
	fmt.Printf("%#v\n", user)
	if len(user.FirstName) == 0 {
		return model.Error{}, errors.New("firstname and lastname cannot be empty")
	}
	temprole, err := json.Marshal(user.Role)
	checkerror.CheckError(err)
	fmt.Println(string(temprole))
	res := connectdb.GiveDB().MustExec("INSERT INTO userswithjob(first_name,last_name,role,title) VALUES ($1,$2,$3,$4)", user.FirstName, user.LastName, temprole, user.JobTitle)
	var rowChanged, err2 = res.RowsAffected()
	checkerror.CheckError(err2)
	fmt.Println("id", rowChanged)
	// var id int64

	// fmt.Println(newUser)
	return model.Error{Message: "User Created Successfully"}, nil
}

// Function to make user from db
func makeUserFromDB(res *sql.Row, ress *sql.Rows) (model.User, error) {
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
		if err != nil {
			return model.User{}, err
		}
		var role []string
		err1 := json.Unmarshal([]byte(temprole), &role)
		if err1 != nil {
			return model.User{}, err1
		}
		var newUser = model.User{ID: id, FirstName: f_name, LastName: l_name, Role: role, JobTitle: title}
		return newUser, nil
	} else {
		var (
			id       int64
			f_name   string
			l_name   string
			temprole string
			title    string
		)
		checkerror.CheckError(ress.Scan(&id,
			&f_name,
			&l_name,
			&temprole,
			&title))
		// fmt.Println(id,
		// 	f_name,
		// 	l_name,
		// 	temprole,
		// 	title)
		var role []string
		checkerror.CheckError(json.Unmarshal([]byte(temprole), &role))
		var newUser = model.User{ID: id, FirstName: f_name, LastName: l_name, Role: role, JobTitle: title}
		return newUser, nil
	}

}
