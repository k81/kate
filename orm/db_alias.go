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
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/k81/kate/log"
)

// DriverType database driver constant int.
type DriverType int

// Enum the Database driver
const (
	_       DriverType = iota // int enum type
	DRMySQL                   // mysql
)

// database driver string.
type driver string

// get type constant int of current driver..
func (d driver) Type() DriverType {
	a, _ := dataBaseCache.get(string(d))
	return a.Driver
}

// get name of current driver
func (d driver) Name() string {
	return string(d)
}

// check driver iis implemented Driver interface or not.
var _ Driver = new(driver)

var (
	dataBaseCache = &_dbCache{cache: make(map[string]*alias)}
	drivers       = map[string]DriverType{
		"mysql": DRMySQL,
	}
	dbBasers = map[DriverType]dbBaser{
		DRMySQL: newdbBaseMysql(),
	}
)

// database alias cacher.
type _dbCache struct {
	mux   sync.RWMutex
	cache map[string]*alias
}

// add database alias with original name.
func (ac *_dbCache) add(name string, al *alias) (added bool) {
	ac.mux.Lock()
	defer ac.mux.Unlock()
	if _, ok := ac.cache[name]; !ok {
		ac.cache[name] = al
		added = true
	}
	return
}

// get database alias if cached.
func (ac *_dbCache) get(name string) (al *alias, ok bool) {
	ac.mux.RLock()
	defer ac.mux.RUnlock()
	al, ok = ac.cache[name]
	return
}

// get default alias.
func (ac *_dbCache) getDefault() (al *alias) {
	al, _ = ac.get("default")
	return
}

type alias struct {
	Name            string
	Driver          DriverType
	DriverName      string
	DataSource      string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	DB              *sql.DB
	DbBaser         dbBaser
	TZ              *time.Location
	Engine          string
}

func detectTZ(al *alias) {
	// orm timezone system match database
	// default use Local
	al.TZ = DefaultTimeLoc

	switch al.Driver {
	case DRMySQL:
		row := al.DB.QueryRow("SELECT TIMEDIFF(NOW(), UTC_TIMESTAMP)")
		var tz string
		row.Scan(&tz)
		if len(tz) >= 8 {
			if tz[0] != '-' {
				tz = "+" + tz
			}
			t, err := time.Parse("-07:00:00", tz)
			if err == nil {
				if t.Location().String() != "" {
					al.TZ = t.Location()
				}
			} else {
				log.Error(mctx, "Detect DB timezone", "time_zone", tz, "error", err)
			}
		}

		// get default engine from current database
		row = al.DB.QueryRow("SELECT ENGINE, TRANSACTIONS FROM information_schema.engines WHERE SUPPORT = 'DEFAULT'")
		var engine string
		var tx bool
		row.Scan(&engine, &tx)

		if engine != "" {
			al.Engine = engine
		} else {
			al.Engine = "INNODB"
		}
	}
}

// RegisterDataBase Setting the database connect params. Use the database driver self dataSource args.
func RegisterDataBase(aliasName, driverName, dataSource string, params ...interface{}) {
	var (
		dr  DriverType
		ok  bool
		al  *alias
		err error
	)

	if dr, ok = drivers[driverName]; !ok {
		panic(fmt.Errorf("driver name `%s` have not registered", driverName))
	}

	al = new(alias)
	al.Name = aliasName
	al.DriverName = driverName
	al.DbBaser = dbBasers[dr]
	al.DataSource = dataSource
	al.Driver = dr

	if al.DB, err = sql.Open(driverName, dataSource); err != nil {
		panic(fmt.Errorf("register db `%s`, %s", aliasName, err.Error()))
	}

	if dataBaseCache.add(aliasName, al) == false {
		panic(fmt.Errorf("DataBase alias name `%s` already registered, cannot reuse", aliasName))
	}

	detectTZ(al)

	for i, v := range params {
		switch i {
		case 0:
			SetMaxIdleConns(al.Name, v.(int))
		case 1:
			SetMaxOpenConns(al.Name, v.(int))
		case 2:
			SetConnMaxLifetime(al.Name, v.(time.Duration))
		}
	}
}

// RegisterDriver Register a database driver use specify driver name, this can be definition the driver is which database type.
func RegisterDriver(driverName string, typ DriverType) {
	if _, ok := drivers[driverName]; ok {
		panic(fmt.Errorf("driverName `%s` db driver already registered", driverName))
	}

	drivers[driverName] = typ
}

// SetDataBaseTZ Change the database default used timezone
func SetDataBaseTZ(aliasName string, tz *time.Location) {
	al, ok := dataBaseCache.get(aliasName)
	if !ok {
		panic(fmt.Errorf("DataBase alias name `%s` not registered", aliasName))
	}

	al.TZ = tz
}

// SetMaxIdleConns Change the max idle conns for *sql.DB, use specify database alias name
func SetMaxIdleConns(aliasName string, maxIdleConns int) {
	al := getDbAlias(aliasName)

	al.MaxIdleConns = maxIdleConns
	al.DB.SetMaxIdleConns(maxIdleConns)
}

// SetMaxOpenConns Change the max open conns for *sql.DB, use specify database alias name
func SetMaxOpenConns(aliasName string, maxOpenConns int) {
	al := getDbAlias(aliasName)

	al.MaxOpenConns = maxOpenConns
	al.DB.SetMaxOpenConns(maxOpenConns)
}

func SetConnMaxLifetime(aliasName string, d time.Duration) {
	al := getDbAlias(aliasName)

	al.ConnMaxLifetime = d
	al.DB.SetConnMaxLifetime(d)
}
