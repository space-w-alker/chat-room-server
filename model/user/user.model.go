package user

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/space-w-alker/chat-room-server/database"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserDTO struct {
	Username string `json:"username" binding:"required,min=1,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

type GetOrUpdateUserDTO struct {
	Username string `json:"username" binding:"min=5,max=50"`
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"min=6,max=50"`
}

type UserList []User

func Create(user *CreateUserDTO) error {
	insertDynStmt := `insert into "User"("id", "username", "email", "password") values($1, $2, $3, $4)`
	_, e := database.Db.Exec(insertDynStmt, uuid.New(), user.Username, user.Email, user.Password)
	return e
}

func Update(id *string, user *GetOrUpdateUserDTO) error {
	updateStmt := `update "User" set "username"=$2, "email"=$3 "password"=$4 where "id"=$1`
	_, e := database.Db.Exec(updateStmt, id, user.Username, user.Email, user.Password)
	return e
}

func GetById(id *string) (user User, e error) {
	getStmt := `select id, username, email, password from "User" where "id" = '$1'`
	row := database.Db.QueryRow(getStmt, user.Id)
	e = row.Scan(user.Id, user.Username, user.Email, user.Password)
	return user, e
}

func GetWhere(getArgs *GetOrUpdateUserDTO) (userList UserList, e error) {
	keys := structs.Names(getArgs)
	values := structs.Values(getArgs)
	_keys := []string{}
	_values := []interface{}{}
	index := 1;
	for i, v := range keys {
		if values[i] != "" {
			_keys = append(_keys,fmt.Sprintf("%v=$%v", strings.ToLower(v), index))
			_values = append(_values, values[i])
			index++
		}
	}
	getStmt := fmt.Sprintf(`select id, username, email, password from "User" where %v`, strings.Join(_keys, " and "))
	fmt.Print(getStmt)
	row, e := database.Db.Query(getStmt, _values...)
	if e != nil {
		return nil, e
	}
	for row.Next() {
		var user User
		e = row.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
		if e != nil {
			return nil, e
		}
		userList = append(userList, user)
	}
	return userList, nil
}

func Delete(id *string) error {
	deleteStmt := `delete from "User" where id=$1`
	_, e := database.Db.Exec(deleteStmt, id)
	return e
}
