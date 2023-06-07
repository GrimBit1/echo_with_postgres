package connectdb

import (
	_ "database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserDB struct {
	Name string
	DBP *sqlx.DB
}

func (u *UserDB) ConnectDB() error {
	connStr := "user=nlab-7 dbname=nlab-7 password=nlab"
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return err
	}
	u.DBP = db
	fmt.Println("Connected")
	return nil
}

func (u *UserDB) CloseDB() {
	u.DBP.Close()
}
