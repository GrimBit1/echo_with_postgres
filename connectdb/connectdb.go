package connectdb

import (
	_ "database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserDB struct {
	Name    string
	DBP     *sqlx.DB
	RoleMap map[int]string
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
func (u *UserDB) FillRoleMap() error {
	res, err := u.DBP.Query(`Select * from roles`)
	if err != nil {
		return err
	}

	for res.Next() {
		var (
			id   int
			role string
		)
		res.Scan(&id, &role)
		u.RoleMap[id] = role
	}
	// fmt.Println(u.RoleMap)
	return nil
}
