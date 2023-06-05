package connectdb

import (
	"database/sql"
	"fmt"
	checkerror "serverwithpostgres/checkError"
)

var db *sql.DB
var err error

func ConnectDB() {
	connStr := "user=nlab-7 dbname=nlab-7 password=nlab"
	db, err = sql.Open("postgres", connStr)
	checkerror.CheckError(err)
	fmt.Println("Connected")
}

func GiveDB() *sql.DB {
	return db
}
