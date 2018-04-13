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
	"reflect"
)

// single model info
type modelInfo struct {
	pkg       string
	name      string
	fullName  string
	db        string
	table     string
	model     interface{}
	fields    *fields
	manual    bool
	addrField reflect.Value //store the original struct value
	uniques   []string
	isThrough bool
}

// new model info
func newModelInfo(val reflect.Value) (mi *modelInfo) {
	mi = &modelInfo{}
	mi.fields = newFields()
	ind := reflect.Indirect(val)
	mi.addrField = val
	mi.name = ind.Type().Name()
	mi.fullName = getFullName(ind.Type())
	addModelFields(mi, ind, "", []int{})
	return
}

// index: FieldByIndex returns the nested field corresponding to index
func addModelFields(mi *modelInfo, ind reflect.Value, mName string, index []int) {
	var (
		err error
		fi  *fieldInfo
		sf  reflect.StructField
	)

	for i := 0; i < ind.NumField(); i++ {
		field := ind.Field(i)
		sf = ind.Type().Field(i)
		// if the field is unexported skip
		if sf.PkgPath != "" {
			continue
		}
		// add anonymous struct fields
		if sf.Anonymous {
			addModelFields(mi, field, mName+"."+sf.Name, append(index, i))
			continue
		}

		fi, err = newFieldInfo(mi, field, sf, mName)
		if err == errSkipField {
			err = nil
			continue
		} else if err != nil {
			break
		}
		//record current field index
		fi.fieldIndex = append(index, i)
		fi.mi = mi
		fi.inModel = true
		if !mi.fields.Add(fi) {
			err = fmt.Errorf("duplicate column name: %s", fi.column)
			break
		}
		if fi.pk {
			if mi.fields.pk != nil {
				err = fmt.Errorf("one model must have one pk field only")
				break
			} else {
				mi.fields.pk = fi
			}
		}
	}

	if err != nil {
		panic(fmt.Errorf("field: %s.%s, %s", ind.Type(), sf.Name, err))
	}
}
