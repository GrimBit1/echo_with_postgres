package main

import (
	"fmt"
	"serverwithpostgres/connectdb"
	"serverwithpostgres/handler"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	userDB := connectdb.UserDB{Name: "Hi", DBP: nil, RoleMap: map[int]string{}}
	e := echo.New()

	err := userDB.ConnectDB()
	if err != nil {
		userDB.CloseDB()
	}
	err = userDB.FillRoleMap()
	if err != nil {
		panic(err)
	}
	handler.ApiHandler(e, userDB)

	e.Logger.Fatal(e.Start(":1323"))

	defer func() {
		fmt.Println("DB disconnected")
		userDB.CloseDB()
	}()
}
