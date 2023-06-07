package main

import (
	"fmt"
	"serverwithpostgres/connectdb"
	"serverwithpostgres/handler"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

var userDB = connectdb.UserDB{Name: "Hi", DBP: nil}

func main() {
	e := echo.New()

	handler.ApiHandler(e, userDB)

	e.Logger.Fatal(e.Start(":1323"))

	defer func() {
		fmt.Println("DB disconnected")
		userDB.CloseDB()
	}()
}
