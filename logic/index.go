package logic

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	// "fmt"

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

		} else {
			return model.User{}, err
		}
	}
	return user, nil
}

// Function to create user
func CreateUser(data []byte) (model.Error, error) {

	var user = model.User{}

	checkerror.CheckError(json.Unmarshal(data, &user))

	// fmt.Printf("%#v\n", user)

	if len(user.FirstName) == 0 {

		return model.Error{}, errors.New("firstname and lastname cannot be empty")

	}
	if len(user.LastName) == 0 {

		return model.Error{}, errors.New("firstname and lastname cannot be empty")

	}
	temprole, err := json.Marshal(user.Role)

	checkerror.CheckError(err)

	// fmt.Println(string(temprole))

	res := connectdb.GiveDB().MustExec("INSERT INTO userswithjob(first_name,last_name,role,title) VALUES ($1,$2,$3,$4)", user.FirstName, user.LastName, temprole, user.JobTitle)

	var rowChanged, err2 = res.RowsAffected()

	checkerror.CheckError(err2)

	fmt.Println("id", rowChanged)

	// var id int64

	// fmt.Println(newUser)

	return model.Error{Message: "User Created Successfully"}, nil
}

// Function to update user in db
func UpdateUser(data []byte, id int64) (model.Error, error) {
	// Check the db that the id exists or not
	var res = connectdb.GiveDB().QueryRow("Select * from userswithjob where id = $1", id)
	oldUser, err := makeUserFromDB(res, nil)
	if err != nil {
		if err.Error() == `sql: no rows in result set` {

			return model.Error{}, errors.New("not found")

		} else {
			return model.Error{}, err
		}
	}

	//Create user tempalate to update old one
	var newUser = model.User{}

	checkerror.CheckError(json.Unmarshal(data, &newUser))

	//Check if newUser has provided the firstname and lastname
	if len(newUser.FirstName) != 0 {

		oldUser.FirstName = newUser.FirstName

	}
	if len(newUser.LastName) != 0 {

		oldUser.LastName = newUser.LastName
	}
	if len(newUser.Role) != 0 {
		oldUser.Role = append(oldUser.Role, newUser.Role...)
	}
	if len(newUser.JobTitle) != 0 {
		oldUser.JobTitle = newUser.JobTitle
	}
	fmt.Println(oldUser)
	// Update the user in db
	var temprole, err1 = json.Marshal(oldUser.Role)
	checkerror.CheckError(err1)
	result := connectdb.GiveDB().MustExec(`Update userswithjob set first_name = $1 ,last_name=$2 ,role=$3 ,title=$4 where id = $5`, oldUser.FirstName, oldUser.LastName, temprole, oldUser.JobTitle,id)
	var rowChanged, err2 = result.RowsAffected()
	checkerror.CheckError(err2)
	if rowChanged > 0 {

		return model.Error{Message: "Updated Successfully"}, nil
	}
	return model.Error{}, errors.New("Some Error Occured")

}

// Function to delete user from db
func DeleteUser(id int64) (model.Error, error) {
	res := connectdb.GiveDB().MustExec(`Delete from userswithjob where id = $1`, id)
	// fmt.Println(res)
	rowChanged, err := res.RowsAffected()
	fmt.Println(err)
	if rowChanged > 0 {

		return model.Error{Message: "User Deleted Successfully"}, nil
	}
	return model.Error{}, errors.New("Bad Request")

}

func GetUsersbyQuery()  {
	
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
