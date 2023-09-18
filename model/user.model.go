package model

import "github.com/space-w-alker/chat-room-server/database"

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserList []User

func (user *User) Create(){
	insertDynStmt := `insert into "user"("id", "username", "email", "password") values($1, $2, $3, $4)`
	_, e := database.Db.Exec(insertDynStmt, user.Id, user.Username, user.Email, user.Password)
  CheckError(e)
}

func (user *User) Update(){
	updateStmt := `update "user" set "username"=$2, "email"=$3 "password"=$4 where "id"=$1`
	_, e := database.Db.Exec(updateStmt, user.Id, user.Username, user.Email, user.Password)
  CheckError(e)
}

func (user *User) GetById(){
	updateStmt := `select * from "user" where "id" = '$1'`
	row := database.Db.QueryRow(updateStmt, user.Id)
	e := row.Scan(user.Id, user.Username, user.Email, user.Password)
	CheckError(e)
}

//func (userList *UserList) GetUsersWhere (user *User){
//	row, e := database.Db.Query(`select * from "user"`)
//}

func CheckError(err error) {
  if err != nil {
    panic(err)
  }
}