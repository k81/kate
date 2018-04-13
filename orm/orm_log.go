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

package orm

import (
	"context"
	"database/sql"
	"time"

	"github.com/k81/kate/log"
	"github.com/k81/kate/utils"
)

func debugLogQueies(ctx context.Context, dbName string, operation, query string, t time.Time, err error, args ...interface{}) {
	var (
		elapsed = int64(time.Since(t) / time.Millisecond)
		flag    = "OK"
	)

	if err != nil {
		flag = "FAIL"
	}

	log.Debug(ctx, "debug sql",
		"db", dbName,
		"flag", flag,
		"operation", operation,
		"elapsed_ms", elapsed,
		"sql", query,
		"args", utils.JoinSlice(args, ","),
		"error", err,
	)
}

// statement query logger struct.
// if dev mode, use stmtQueryLog, or use stmtQuerier.
type stmtQueryLog struct {
	dbName string
	query  string
	stmt   stmtQuerier
	ctx    context.Context
}

var _ stmtQuerier = new(stmtQueryLog)

func (d *stmtQueryLog) Close() error {
	a := time.Now()
	err := d.stmt.Close()
	debugLogQueies(d.ctx, d.dbName, "st.Close", d.query, a, err)
	return err
}

func (d *stmtQueryLog) Exec(args ...interface{}) (sql.Result, error) {
	a := time.Now()
	res, err := d.stmt.Exec(args...)
	debugLogQueies(d.ctx, d.dbName, "st.Exec", d.query, a, err, args...)
	return res, err
}

func (d *stmtQueryLog) Query(args ...interface{}) (*sql.Rows, error) {
	a := time.Now()
	res, err := d.stmt.Query(args...)
	debugLogQueies(d.ctx, d.dbName, "st.Query", d.query, a, err, args...)
	return res, err
}

func (d *stmtQueryLog) QueryRow(args ...interface{}) *sql.Row {
	a := time.Now()
	res := d.stmt.QueryRow(args...)
	debugLogQueies(d.ctx, d.dbName, "st.QueryRow", d.query, a, nil, args...)
	return res
}

func newStmtQueryLog(ctx context.Context, dbName string, stmt stmtQuerier, query string) stmtQuerier {
	d := new(stmtQueryLog)
	d.ctx = ctx
	d.stmt = stmt
	d.dbName = dbName
	d.query = query
	return d
}

// database query logger struct.
// if dev mode, use dbQueryLog, or use dbQuerier.
type dbQueryLog struct {
	dbName string
	db     dbQuerier
	tx     txer
	txe    txEnder
	ctx    context.Context
}

var _ dbQuerier = new(dbQueryLog)
var _ txer = new(dbQueryLog)
var _ txEnder = new(dbQueryLog)

func (d *dbQueryLog) Prepare(query string) (*sql.Stmt, error) {
	a := time.Now()
	stmt, err := d.db.Prepare(query)
	debugLogQueies(d.ctx, d.dbName, "db.Prepare", query, a, err)
	return stmt, err
}

func (d *dbQueryLog) Exec(query string, args ...interface{}) (sql.Result, error) {
	a := time.Now()
	res, err := d.db.Exec(query, args...)
	debugLogQueies(d.ctx, d.dbName, "db.Exec", query, a, err, args...)
	return res, err
}

func (d *dbQueryLog) Query(query string, args ...interface{}) (*sql.Rows, error) {
	a := time.Now()
	res, err := d.db.Query(query, args...)
	debugLogQueies(d.ctx, d.dbName, "db.Query", query, a, err, args...)
	return res, err
}

func (d *dbQueryLog) QueryRow(query string, args ...interface{}) *sql.Row {
	a := time.Now()
	res := d.db.QueryRow(query, args...)
	debugLogQueies(d.ctx, d.dbName, "db.QueryRow", query, a, nil, args...)
	return res
}

func (d *dbQueryLog) Begin() (*sql.Tx, error) {
	a := time.Now()
	tx, err := d.db.(txer).Begin()
	debugLogQueies(d.ctx, d.dbName, "db.Begin", "START TRANSACTION", a, err)
	return tx, err
}

func (d *dbQueryLog) Commit() error {
	a := time.Now()
	err := d.db.(txEnder).Commit()
	debugLogQueies(d.ctx, d.dbName, "tx.Commit", "COMMIT", a, err)
	return err
}

func (d *dbQueryLog) Rollback() error {
	a := time.Now()
	err := d.db.(txEnder).Rollback()
	debugLogQueies(d.ctx, d.dbName, "tx.Rollback", "ROLLBACK", a, err)
	return err
}

func (d *dbQueryLog) SetDB(db dbQuerier) {
	d.db = db
}

func newDbQueryLog(ctx context.Context, dbName string, db dbQuerier) dbQuerier {
	d := new(dbQueryLog)
	d.ctx = ctx
	d.dbName = dbName
	d.db = db
	return d
}
