package model

import "github.com/space-w-alker/chat-room-server/database"

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserList []User

func (user *User) Create() {
	insertDynStmt := `insert into "User"("username", "email", "password") values($1, $2, $3)`
	_, e := database.Db.Exec(insertDynStmt, user.Username, user.Email, user.Password)
	checkError(e)
}

func (user *User) Update() {
	updateStmt := `update "User" set "username"=$2, "email"=$3 "password"=$4 where "id"=$1`
	_, e := database.Db.Exec(updateStmt, user.Id, user.Username, user.Email, user.Password)
	checkError(e)
}

func (user *User) GetById() {
	getStmt := `select id, username, email, password from "User" where "id" = '$1'`
	row := database.Db.QueryRow(getStmt, user.Id)
	e := row.Scan(user.Id, user.Username, user.Email, user.Password)
	checkError(e)
}

func (userList *UserList) GetUsersWhere(user *map[string]interface{}) {
	row, e := database.Db.Query(`select id, username, email, password from "User"`)
	checkError(e)
	for row.Next() {
		var user User
		e := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
		checkError(e)
		*userList = append(*userList, user)
	}
}

func (userList *UserList) DeleteUser() {
	deleteStmt := `delete from "User" where id=$1`
	_, e := database.Db.Exec(deleteStmt, 1)
	checkError(e)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
