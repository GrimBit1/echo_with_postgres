package model

/*
Schema

	id BIGSERIAL PRIMARY KEY ,
	first_name VARCHAR(150),
	last_name VARCHAR(150),
	role JSON,
	title VARCHAR(150)
*/
type User struct {
	ID        int64    `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Role      []string `json:"role"`
	JobTitle  string   `json:"title"`
}
type Error struct {
	Message string `json:"message"`
}

var UserSchema = `
CREATE TABLE userswithjob(id BIGSERIAL PRIMARY KEY ,
	first_name VARCHAR(150),
last_name VARCHAR(150),
role JSON,
title VARCHAR(150))
`
