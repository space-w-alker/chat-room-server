package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/space-w-alker/chat-room-server/model/generic"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

type User struct {
	Id        *string    `json:"id" form:"id" db:"id"`
	Username  *string    `json:"username" form:"username" db:"username"`
	Email     *string    `json:"email" form:"email" db:"email"`
	Password  *string    `json:"password" form:"password" db:"password"`
	CreatedAt *time.Time `json:"createdAt" form:"createdAt" db:"createdAt"`
}

type UserList []User

func Create(db *goqu.Database, u *User) (string, error) {
	id := uuid.New().String()
	m := map[string]interface{}{}
	u.Id = &id
	generic.ToJsMap(*u, m)
	_, e := db.Insert("User").Rows(m).Executor().Exec()
	return id, e
}

func Update(db *goqu.Database, id *string, u *User) error {
	m := map[string]interface{}{}
	generic.ToJsMap(*u, m)
	_, e := db.Update("User").Set(m).Where(goqu.Ex{"id": id}).Executor().Exec()
	return e
}

func GetById(db *goqu.Database, id *string) (u *User, e error) {
	u = &User{}
	_, e = db.From("User").Where(goqu.Ex{"id": *id}).ScanStruct(u)
	return u, e
}

func GetWhere(db *goqu.Database, getArgs *User, opts *generic.PaginationArgs) (userList UserList, meta generic.PaginationMeta, e error) {
	where := goqu.Ex{}
	generic.ToJsMap(getArgs, where)
	offset := (opts.Page - 1) * opts.Limit
	query := db.From("User").Order(goqu.C("createdAt").Desc().NullsLast()).Offset(offset).Limit(opts.Limit).Where(where)
	countQuery := db.From("User").Where(where)
	if err := query.ScanStructs(&userList); err != nil {
		return nil, meta, err
	}
	total, err := countQuery.Count()
	if err != nil {
		return nil, meta, err
	}
	meta = generic.PaginationMeta{TotalItems: uint(total), Page: opts.Page, Limit: opts.Limit}
	return userList, meta, nil
}

func Delete(db *goqu.Database, id *string) error {
	_,e := db.Delete("User").Where(goqu.Ex{"id":id}).Executor().Exec()
	return e
}
