package chat

import (
	"time"

	"github.com/google/uuid"
	"github.com/space-w-alker/chat-room-server/model/generic"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

type Chat struct {
	Id         *string    `json:"id" form:"id" db:"id"`
	Text       *string    `json:"text" form:"text" db:"text"`
	Published  *bool      `json:"published" form:"published" db:"published"`
	AuthorId   *string    `json:"authorId" form:"authorId" db:"authorId"`
	ChatRoomId *string    `json:"chatRoomId" form:"chatRoomId" db:"chatRoomId"`
	CreatedAt  *time.Time `json:"createdAt" form:"createdAt" db:"createdAt"`
}

type ChatList []Chat

func Create(db *goqu.Database, m *Chat) (string, error) {
	id := uuid.New().String()
	mp := map[string]interface{}{}
	m.Id = &id
	generic.ToJsMap(*m, mp)
	_, e := db.Insert("Chat").Rows(mp).Executor().Exec()
	return id, e
}

func Update(db *goqu.Database, id *string, m *Chat) error {
	mp := map[string]interface{}{}
	generic.ToJsMap(*m, mp)
	_, e := db.Update("Chat").Set(mp).Where(goqu.Ex{"id": id}).Executor().Exec()
	return e
}

func GetById(db *goqu.Database, id *string) (m *Chat, e error) {
	m = &Chat{}
	_, e = db.From("Chat").Where(goqu.Ex{"id": *id}).ScanStruct(m)
	return m, e
}

func GetWhere(db *goqu.Database, getArgs *Chat, opts *generic.PaginationArgs) (mList ChatList, meta generic.PaginationMeta, e error) {
	where := goqu.Ex{}
	generic.ToJsMap(getArgs, where)
	offset := (opts.Page - 1) * opts.Limit
	query := db.From("Chat").Order(goqu.C("createdAt").Desc().NullsLast()).Offset(offset).Limit(opts.Limit).Where(where)
	countQuery := db.From("Chat").Where(where)
	if err := query.ScanStructs(&mList); err != nil {
		return nil, meta, err
	}
	total, err := countQuery.Count()
	if err != nil {
		return nil, meta, err
	}
	meta = generic.PaginationMeta{TotalItems: uint(total), Page: opts.Page, Limit: opts.Limit}
	return mList, meta, nil
}

func Delete(db *goqu.Database, id *string) error {
	_, e := db.Delete("Chat").Where(goqu.Ex{"id": id}).Executor().Exec()
	return e
}
