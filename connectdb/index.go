package connectdb

import (
	_ "database/sql"
	"fmt"
	checkerror "serverwithpostgres/checkError"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var err error

func ConnectDB() {
	connStr := "user=nlab-7 dbname=nlab-7 password=nlab"
	db, err = sqlx.Open("postgres", connStr)
	checkerror.CheckError(err)
	fmt.Println("Connected")
}

func GiveDB() *sqlx.DB {
	return db
}
