package dborm

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB

func SetDB(gdb *gorm.DB) {
	db = gdb
}

type table interface {
	TableName() string
}

type Common struct {
	t  table
	db *gorm.DB
}

func New(t table, tx ...*gorm.DB) Common {
	_db := db
	if len(tx) > 0 {
		_db = tx[0]
	}
	return Common{t, _db}
}

func (c Common) First(where []Where, order string, target interface{}) error {
	return ToWhere(c.db.Table(c.t.TableName()), where).Order(order).First(target).Error
}

func (c Common) Find(where []Where, order string, target interface{}) error {
	return ToWhere(c.db.Table(c.t.TableName()), where).Order(order).Find(target).Error
}

func (c Common) Insert(t table) error {
	return c.db.Create(t).Error
}

func (c Common) Save(t table) error {
	return c.db.Save(t).Error
}

func (c Common) Update(where []Where, data map[string]interface{}) (int64, error) {
	obj := reflect.TypeOf(c.t)
	if _, exists := obj.FieldByName("UpdatedAt"); exists {
		if data["updated_at"] == nil && reflect.ValueOf(data["updated_at"]).IsZero() {
			data["updated_at"] = time.Now().Unix()
		}
	}
	res := ToWhere(c.db.Table(c.t.TableName()), where).UpdateColumns(data)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (c Common) Delete(where []Where) (int64, error) {
	res := ToWhere(c.db.Table(c.t.TableName()), where).Delete(c.t)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (c Common) SoftDelete(where []Where) (int64, error) {
	timestamp := time.Now().Unix()
	return c.Update(where, map[string]interface{}{
		"deleted_at": timestamp,
	})
}

func (c Common) ListPageWithCount(where []Where, offset, limit int, order string, target interface{}) (count int64, err error) {
	query := ToWhere(c.db.Table(c.t.TableName()), where)
	err = query.Count(&count).Error
	if count > 0 {
		err = query.Offset(offset).Limit(limit).Order(order).Find(target).Error
	}
	return
}

type Where struct {
	Field string
	Op    string
	Value interface{}
}

func ToWhere(query *gorm.DB, where []Where) *gorm.DB {
	if where == nil {
		return query
	}
	for _, w := range where {
		query = w.toWhere(query)
	}
	return query
}

func (w Where) toWhere(query *gorm.DB) *gorm.DB {
	switch strings.ToUpper(w.Op) {
	case "IN":
	case "NOT IN":
		// uint8
		if uv, ok := w.Value.([]uint8); ok {
			nv := make([]int, 0)
			for _, v := range uv {
				nv = append(nv, int(v))
			}
			query = query.Where(fmt.Sprintf("%s %s (?)", w.Field, w.Op), nv)
		} else {
			query = query.Where(fmt.Sprintf("%s %s (?)", w.Field, w.Op), w.Value)
		}
	case "LIKE":
	case "NOT LIKE":
		if strings.HasPrefix(w.Value.(string), "%") || strings.HasSuffix(w.Value.(string), "%") {
			query = query.Where(fmt.Sprintf("%s %s ?", w.Field, w.Op), w.Value.(string))
		} else {
			query = query.Where(fmt.Sprintf("%s %s ?", w.Field, w.Op), "%"+w.Value.(string)+"%")
		}
	case "BETWEEN":
		v := w.Value.([]interface{})
		query = query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", w.Field), v[0], v[1])
	case "FIND_IN_SET":
		query = query.Where(fmt.Sprintf("FIND_IN_SET(?, %s)", w.Field), w.Value)
	case "RAW":
		query = query.Where(w.Field, w.Value.([]interface{})...)
	default:
		query = query.Where(fmt.Sprintf("%s %s ?", w.Field, w.Op), w.Value)
	}

	return query
}
