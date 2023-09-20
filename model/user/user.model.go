package user

import (
	"time"

	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/space-w-alker/chat-room-server/database"
	"github.com/space-w-alker/chat-room-server/model/generic"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

type User struct {
	Id        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
}

type CreateUserDTO struct {
	Username string `json:"username" binding:"required,min=1,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

type GetOrUpdateUserDTO struct {
	Username string `json:"username" form:"username" binding:"omitempty,min=1,max=50"`
	Email    string `json:"email" form:"email" binding:"omitempty,email"`
	Password string `json:"password" form:"password" binding:"omitempty,min=6,max=50"`
}

type UserList []User

func Create(user *CreateUserDTO) (string, error) {
	insertDynStmt := `insert into "User"("id", "username", "email", "password") values($1, $2, $3, $4)`
	id := uuid.New()
	_, e := database.Db.Exec(insertDynStmt, id, user.Username, user.Email, user.Password)
	return id.String(), e
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

func GetWhere(getArgs *GetOrUpdateUserDTO, opts *generic.PaginationArgs) (userList UserList, e error) {
	db := database.DB
	_where := map[string]interface{}{}
	where := goqu.Ex{}
	structs.FillMap(getArgs, _where)
	for _, field := range structs.Fields(getArgs) {
		if _where[field.Name()] != "" {
			where[field.Tag("json")] = _where[field.Name()]
		}
	}
	offset := (opts.Page - 1) * opts.Limit
	query := db.From("User").Order(goqu.C("createdAt").Desc().NullsLast()).Offset(offset).Limit(opts.Limit).Where(where)
	if err := query.ScanStructs(&userList); err != nil {
		return nil, err
	}
	return userList, nil
}

func Delete(id *string) error {
	deleteStmt := `delete from "User" where id=$1`
	_, e := database.Db.Exec(deleteStmt, id)
	return e
}
