package orm

import (
	"context"
	"reflect"
	"strings"

	"github.com/k81/kate/log"
)

// 1 is attr
// 2 is tag
var supportTag = map[string]int{
	"-":      1,
	"pk":     1,
	"auto":   1,
	"json":   1,
	"column": 2,
	"type":   2,
}

// get reflect.Type name with package path.
func getFullName(typ reflect.Type) string {
	return typ.PkgPath() + "." + typ.Name()
}

// getTableName get struct table name.
// If the struct implement the TableName, then get the result as tablename
// else use the struct name which will apply snakeString.
func getTableName(val reflect.Value) string {
	if fun := val.MethodByName("TableName"); fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		// has return and the first val is string
		if len(vals) > 0 && vals[0].Kind() == reflect.String {
			return vals[0].String()
		}
	}
	return snakeString(reflect.Indirect(val).Type().Name())
}

func isSharded(val reflect.Value) bool {
	if fun := val.MethodByName("TableSuffix"); fun.IsValid() {
		return true
	}
	return false
}

func getTableSuffix(ind reflect.Value) string {
	if !ind.CanAddr() {
		return ""
	}

	val := ind.Addr()
	if fun := val.MethodByName("TableSuffix"); fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		if len(vals) > 0 && vals[0].Kind() == reflect.String {
			return vals[0].String()
		}
	}
	return ""
}

// get snaked column name
func getColumnName(sf reflect.StructField, col string) string {
	column := col
	if col == "" {
		column = snakeString(sf.Name)
	}
	return column
}

// parse struct tag string
func parseStructTag(data string) (attrs map[string]bool, tags map[string]string) {
	attrs = make(map[string]bool)
	tags = make(map[string]string)
	for _, v := range strings.Split(data, defaultStructTagDelim) {
		if v == "" {
			continue
		}
		v = strings.TrimSpace(v)
		if t := strings.ToLower(v); supportTag[t] == 1 {
			attrs[t] = true
		} else if i := strings.Index(v, "("); i > 0 && strings.Index(v, ")") == len(v)-1 {
			name := t[:i]
			if supportTag[name] == 2 {
				v = v[i+1 : len(v)-1]
				tags[name] = v
			}
		} else {
			log.Error(context.TODO(), "unsupport orm tag", "tag", v)
		}
	}
	return
}
