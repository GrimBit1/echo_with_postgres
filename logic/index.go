package logic

import (
	"fmt"
	checkerror "serverwithpostgres/checkError"
	"serverwithpostgres/connectdb"
)

func GiveAllUsers() {
	res, err := connectdb.GiveDB().
	checkerror.CheckError(err)
	var (
		id     int
		f_name string
		l_name string
		role   []string
		title  string
	)
	res.Scan(&id,
		&f_name,
		&l_name,
		&role,
		&title)
	fmt.Println(id,
		f_name,
		l_name,
		role,
		title)
}
