package main

import (
	connectdb "serverwithpostgres/connectdb"
	"serverwithpostgres/handler"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	e := echo.New()
	handler.ApiHandler(e)
	connectdb.ConnectDB()
	e.Logger.Fatal(e.Start(":1323"))
}
