package log

import (
	"fmt"
	"reflect"
)

type Formatter interface {
	Format(*Entry) ([]byte, error)
}

func toString(k interface{}) (str string) {
	switch x := k.(type) {
	case string:
		str = x
	case fmt.Stringer:
		str = safeString(x)
	default:
		str = fmt.Sprint(x)
	}
	return
}

func safeString(s fmt.Stringer) (str string) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			if v := reflect.ValueOf(s); v.Kind() == reflect.Ptr && v.IsNil() {
				str = "NULL"
			} else {
				panic(panicVal)
			}
		}
	}()

	str = s.String()
	return
}
