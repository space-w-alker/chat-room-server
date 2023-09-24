package room_tag

import (
	"time"

	"github.com/google/uuid"
	"github.com/space-w-alker/chat-room-server/model/generic"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

type RoomTag struct {
	Id *string `json:"id" form:"id" db:"id"`
Name *string `json:"name" form:"name" db:"name"`
CreatedAt *time.Time `json:"createdAt" form:"createdAt" db:"createdAt"`
}

type RoomTagList []RoomTag

func Create(db *goqu.Database, m *RoomTag) (string, error) {
	id := uuid.New().String()
	mp := map[string]interface{}{}
	m.Id = &id
	generic.ToJsMap(*m, mp)
	_, e := db.Insert("RoomTag").Rows(mp).Executor().Exec()
	return id, e
}

func Update(db *goqu.Database, id *string, m *RoomTag) error {
	mp := map[string]interface{}{}
	generic.ToJsMap(*m, mp)
	_, e := db.Update("RoomTag").Set(mp).Where(goqu.Ex{"id": id}).Executor().Exec()
	return e
}

func GetById(db *goqu.Database, id *string) (m *RoomTag, e error) {
	m = &RoomTag{}
	_, e = db.From("RoomTag").Where(goqu.Ex{"id": *id}).ScanStruct(m)
	return m, e
}

func GetWhere(db *goqu.Database, getArgs *RoomTag, opts *generic.PaginationArgs) (mList RoomTagList, meta generic.PaginationMeta, e error) {
	where := goqu.Ex{}
	generic.ToJsMap(getArgs, where)
	offset := (opts.Page - 1) * opts.Limit
	query := db.From("RoomTag").Order(goqu.C("createdAt").Desc().NullsLast()).Offset(offset).Limit(opts.Limit).Where(where)
	countQuery := db.From("RoomTag").Where(where)
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
	_,e := db.Delete("RoomTag").Where(goqu.Ex{"id":id}).Executor().Exec()
	return e
}
