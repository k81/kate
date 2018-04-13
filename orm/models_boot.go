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
	"fmt"
	"os"
	"reflect"
	"strings"
)

// register models.
// PrefixOrSuffix means table name prefix or suffix.
// isPrefix whether the prefix is prefix or suffix
func registerModel(PrefixOrSuffix, db string, model interface{}, isPrefix bool) {
	val := reflect.ValueOf(model)
	typ := reflect.Indirect(val).Type()

	if val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("<orm.RegisterModel> cannot use non-ptr model struct `%s`", getFullName(typ)))
	}
	// For this case:
	// u := &User{}
	// registerModel(&u)
	if typ.Kind() == reflect.Ptr {
		panic(fmt.Errorf("<orm.RegisterModel> only allow ptr model struct, it looks you use two reference to the struct `%s`", typ))
	}

	table := getTableName(val)

	if PrefixOrSuffix != "" {
		if isPrefix {
			table = PrefixOrSuffix + table
		} else {
			table = table + PrefixOrSuffix
		}
	}
	// models's fullname is pkgpath + struct name
	name := getFullName(typ)
	if _, ok := modelCache.getByFullName(name); ok {
		panic(fmt.Errorf("<orm.RegisterModel> model `%s` repeat register, must be unique\n", name))
	}

	if _, ok := modelCache.get(table); ok {
		fmt.Printf("<orm.RegisterModel> table name `%s` repeat register, must be unique\n", table)
		os.Exit(2)
	}

	mi := newModelInfo(val)
	if mi.fields.pk == nil {
	outFor:
		for _, fi := range mi.fields.fieldsDB {
			if strings.ToLower(fi.name) == "id" {
				switch fi.addrValue.Elem().Kind() {
				case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint32, reflect.Uint64:
					fi.auto = true
					fi.pk = true
					mi.fields.pk = fi
					break outFor
				}
			}
		}

		if mi.fields.pk == nil {
			panic(fmt.Errorf("<orm.RegisterModel> `%s` need a primary key field, default use 'id' if not set", name))
		}

	}

	mi.db = db
	mi.table = table
	mi.pkg = typ.PkgPath()
	mi.model = model
	mi.manual = true

	modelCache.set(table, mi)
}

// boostrap models
func bootStrap() {
	if modelCache.done {
		return
	}

	if dataBaseCache.getDefault() == nil {
		panic(fmt.Errorf("must have one register DataBase alias named `default`"))
	}
}

// RegisterModel register models
func RegisterModel(db string, models ...interface{}) {
	if modelCache.done {
		panic(fmt.Errorf("RegisterModel must be run before BootStrap"))
	}
	RegisterModelWithPrefix("", db, models...)
}

// RegisterModelWithPrefix register models with a prefix
func RegisterModelWithPrefix(prefix, db string, models ...interface{}) {
	if modelCache.done {
		panic(fmt.Errorf("RegisterModelWithPrefix must be run before BootStrap"))
	}

	for _, model := range models {
		registerModel(prefix, db, model, true)
	}
}

// RegisterModelWithSuffix register models with a suffix
func RegisterModelWithSuffix(suffix, db string, models ...interface{}) {
	if modelCache.done {
		panic(fmt.Errorf("RegisterModelWithSuffix must be run before BootStrap"))
	}

	for _, model := range models {
		registerModel(suffix, db, model, false)
	}
}

// BootStrap bootrap models.
// make all model parsed and can not add more models
func BootStrap() {
	if modelCache.done {
		return
	}
	modelCache.Lock()
	defer modelCache.Unlock()
	bootStrap()
	modelCache.done = true
}
