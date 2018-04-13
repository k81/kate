// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package orm provide ORM for MySQL/PostgreSQL/sqlite
// Simple Usage
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/k81/kate/orm"
//		_ "github.com/go-sql-driver/mysql" // import your used driver
//	)
//
//	// Model Struct
//	type User struct {
//		Id   int    `orm:"auto"`
//		Name string `orm:"size(100)"`
//	}
//
//	func init() {
//		orm.RegisterDataBase("default", "mysql", "root:root@/my_db?charset=utf8", 30)
//	}
//
//	func main() {
//		o := orm.NewOrm()
//		user := User{Name: "slene"}
//		// insert
//		id, err := o.Insert(&user)
//		// update
//		user.Name = "astaxie"
//		num, err := o.Update(&user)
//		// read one
//		u := User{Id: user.Id}
//		err = o.Read(&u)
//		// delete
//		num, err = o.Delete(&u)
//	}
//
// more docs: http://beego.me/docs/mvc/model/overview.md
package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/k81/kate/log"
)

// Define common vars
var (
	Debug            = false
	DefaultRowsLimit = 1000
	DefaultTimeLoc   = time.Local
	ErrTxHasBegan    = errors.New("<Ormer.Begin> transaction already begin")
	ErrTxDone        = errors.New("<Ormer.Commit/Rollback> transaction not begin")
	ErrMultiRows     = errors.New("<QuerySeter> return multi rows")
	ErrNoRows        = errors.New("<QuerySeter> no row found")
	ErrStmtClosed    = errors.New("<QuerySeter> stmt already closed")
	ErrArgs          = errors.New("<Ormer> args error may be empty")
	ErrNotImplement  = errors.New("have not implement")
	mctx             = log.SetContext(context.Background(), "module", "orm")
)

// Params stores the Params
type Params map[string]interface{}

// ParamsList stores paramslist
type ParamsList []interface{}

type orm struct {
	ctx    context.Context
	alias  *alias
	db     dbQuerier
	dbName string
	isTx   bool
}

var _ Ormer = new(orm)

// get model info and model reflect value
func (o *orm) getMiInd(md interface{}, needPtr bool) (mi *modelInfo, ind reflect.Value) {
	val := reflect.ValueOf(md)
	ind = reflect.Indirect(val)
	typ := ind.Type()
	if needPtr && val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("<Ormer> cannot use non-ptr model struct `%s`", getFullName(typ)))
	}
	name := getFullName(typ)
	if mi, ok := modelCache.getByFullName(name); ok {
		return mi, ind
	}
	panic(fmt.Errorf("<Ormer> table: `%s` not found, maybe not RegisterModel", name))
}

// get field info from model info by given field name
func (o *orm) getFieldInfo(mi *modelInfo, name string) *fieldInfo {
	fi, ok := mi.fields.GetByAny(name)
	if !ok {
		panic(fmt.Errorf("<Ormer> cannot find field `%s` for model `%s`", name, mi.fullName))
	}
	return fi
}

// read data to model
func (o *orm) Read(md interface{}, cols ...string) error {
	mi, ind := o.getMiInd(md, true)
	if !o.isTx {
		o.setDbByMiInd(mi, ind)
	}
	return o.alias.DbBaser.Read(o.db, mi, ind, o.alias.TZ, cols, false)
}

// read data to model, like Read(), but use "SELECT FOR UPDATE" form
func (o *orm) ReadForUpdate(md interface{}, cols ...string) error {
	mi, ind := o.getMiInd(md, true)
	if !o.isTx {
		o.setDbByMiInd(mi, ind)
	}
	return o.alias.DbBaser.Read(o.db, mi, ind, o.alias.TZ, cols, true)
}

// Try to read a row from the database, or insert one if it doesn't exist
func (o *orm) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	cols = append([]string{col1}, cols...)
	mi, ind := o.getMiInd(md, true)
	if !o.isTx {
		o.setDbByMiInd(mi, ind)
	}
	err := o.alias.DbBaser.Read(o.db, mi, ind, o.alias.TZ, cols, false)
	if err == ErrNoRows {
		// Create
		id, err := o.Insert(md)
		return (err == nil), id, err
	}

	return false, ind.FieldByIndex(mi.fields.pk.fieldIndex).Int(), err
}

// insert model data to database
func (o *orm) Insert(md interface{}) (int64, error) {
	mi, ind := o.getMiInd(md, true)
	if !o.isTx {
		o.setDbByMiInd(mi, ind)
	}
	id, err := o.alias.DbBaser.Insert(o.db, mi, ind, o.alias.TZ)
	if err != nil {
		return id, err
	}

	o.setPk(mi, ind, id)

	return id, nil
}

// set auto pk field
func (o *orm) setPk(mi *modelInfo, ind reflect.Value, id int64) {
	if mi.fields.pk.auto {
		if mi.fields.pk.fieldType&IsPositiveIntegerField > 0 {
			ind.FieldByIndex(mi.fields.pk.fieldIndex).SetUint(uint64(id))
		} else {
			ind.FieldByIndex(mi.fields.pk.fieldIndex).SetInt(id)
		}
	}
}

// insert some models to database
func (o *orm) InsertMulti(bulk int, mds interface{}) (int64, error) {
	var cnt int64

	sind := reflect.Indirect(reflect.ValueOf(mds))

	switch sind.Kind() {
	case reflect.Array, reflect.Slice:
		if sind.Len() == 0 {
			return cnt, ErrArgs
		}
	default:
		return cnt, ErrArgs
	}

	if bulk <= 1 {
		for i := 0; i < sind.Len(); i++ {
			ind := reflect.Indirect(sind.Index(i))
			mi, _ := o.getMiInd(ind.Interface(), false)
			if !o.isTx {
				o.setDbByMiInd(mi, ind)
			}
			id, err := o.alias.DbBaser.Insert(o.db, mi, ind, o.alias.TZ)
			if err != nil {
				return cnt, err
			}

			o.setPk(mi, ind, id)

			cnt++
		}
	} else {
		mi, ind := o.getMiInd(sind.Index(0).Interface(), false)
		if !o.isTx {
			o.setDbByMiInd(mi, ind)
		}
		return o.alias.DbBaser.InsertMulti(o.db, mi, sind, bulk, o.alias.TZ)
	}
	return cnt, nil
}

// InsertOrUpdate data to database
func (o *orm) InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error) {
	mi, ind := o.getMiInd(md, true)
	if !o.isTx {
		o.setDbByMiInd(mi, ind)
	}
	return o.alias.DbBaser.InsertOrUpdate(o.db, mi, ind, o.alias, colConflitAndArgs...)
}

// update model to database.
// cols set the columns those want to update.
func (o *orm) Update(md interface{}, cols ...string) (int64, error) {
	mi, ind := o.getMiInd(md, true)
	if !o.isTx {
		o.setDbByMiInd(mi, ind)
	}
	return o.alias.DbBaser.Update(o.db, mi, ind, o.alias.TZ, cols)
}

// delete model in database
// cols shows the delete conditions values read from. default is pk
func (o *orm) Delete(md interface{}, cols ...string) (int64, error) {
	mi, ind := o.getMiInd(md, true)
	if !o.isTx {
		o.setDbByMiInd(mi, ind)
	}
	num, err := o.alias.DbBaser.Delete(o.db, mi, ind, o.alias.TZ, cols)
	if err != nil {
		return num, err
	}
	if num > 0 {
		o.setPk(mi, ind, 0)
	}
	return num, nil
}

// return a QuerySeter for table operations.
// table name can be string or struct.
// e.g. QueryTable("user"), QueryTable(&user{}) or QueryTable((*User)(nil)),
func (o *orm) QueryTable(ptrStructOrTableName interface{}) (qs QuerySeter) {
	var name string
	if table, ok := ptrStructOrTableName.(string); ok {
		name = snakeString(table)
		if mi, ok := modelCache.get(name); ok {
			qs = newQuerySet(o, mi)
		}
	} else {
		name = getFullName(indirectType(reflect.TypeOf(ptrStructOrTableName)))
		if mi, ok := modelCache.getByFullName(name); ok {
			qs = newQuerySet(o, mi)
		}
	}
	if qs == nil {
		panic(fmt.Errorf("<Ormer.QueryTable> table name: `%s` not exists", name))
	}
	return
}

func (o *orm) Using(dbName string) error {
	o.setDb(dbName)
	return nil
}

func (o *orm) setDb(dbName string) {
	if o.isTx {
		panic(fmt.Errorf("<Ormer.Using> transaction has been start, cannot change db"))
	}

	al, ok := dataBaseCache.get(dbName)
	if !ok {
		al = dataBaseCache.getDefault()
	}

	o.dbName = dbName
	o.alias = al

	db := al.DB

	if Debug {
		o.db = newDbQueryLog(o.ctx, dbName, db)
	} else {
		o.db = db
	}
}

// switch to another registered database driver by given name.
func (o *orm) setDbByMiInd(mi *modelInfo, ind reflect.Value) {
	var (
		dbName = mi.db
	)
	if o.db == nil {
		o.setDb(dbName)
	}
}

// begin transaction
func (o *orm) Begin() error {
	if o.isTx {
		return ErrTxHasBegan
	}
	var tx *sql.Tx
	tx, err := o.db.(txer).Begin()
	if err != nil {
		return err
	}
	o.isTx = true
	if Debug {
		o.db.(*dbQueryLog).SetDB(tx)
	} else {
		o.db = tx
	}
	return nil
}

// commit transaction
func (o *orm) Commit() error {
	if !o.isTx {
		return ErrTxDone
	}
	err := o.db.(txEnder).Commit()
	if err == nil {
		o.isTx = false
	} else if err == sql.ErrTxDone {
		return ErrTxDone
	}
	return err
}

// rollback transaction
func (o *orm) Rollback() error {
	if !o.isTx {
		return ErrTxDone
	}
	err := o.db.(txEnder).Rollback()
	if err == nil {
		o.isTx = false
	} else if err == sql.ErrTxDone {
		return ErrTxDone
	}
	return err
}

// return a raw query seter for raw sql string.
func (o *orm) Raw(query string, args ...interface{}) RawSeter {
	return newRawSet(o, query, args)
}

// return current using database Driver
func (o *orm) Driver() Driver {
	return driver(o.alias.Name)
}

func NewOrm(ctx context.Context) Ormer {
	BootStrap() // execute only once

	o := new(orm)
	o.ctx = ctx
	return o
}
