package queryparams

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/k81/kate/utils"
)

// QueryParams for pagination
type QueryParams struct {
	filters map[string]interface{}
	orderBy []string
	page    int
	perPage int
}

var queryRequiredFields = map[string]reflect.Type{
	"Page":    reflect.TypeOf(int(0)),
	"PerPage": reflect.TypeOf(int(0)),
	"Sort":    reflect.TypeOf([]string{}),
}

// NewQueryParams return a QueryParams
func NewQueryParams() *QueryParams {
	p := &QueryParams{
		filters: make(map[string]interface{}),
	}
	return p
}

// NewQueryParamsFromTag create QueryParams from struct
func NewQueryParamsFromTag(ptr interface{}) *QueryParams {
	params := NewQueryParams()

	val := reflect.ValueOf(ptr)
	ind := reflect.Indirect(val)
	typ := ind.Type()
	fullName := typ.PkgPath() + "." + typ.Name()

	if val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("GetQueryParams: cannot use non-ptr struct `%s`", fullName))
	}

	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("GetQueryParams: only allow ptr of struct"))
	}

	for name, expectedType := range queryRequiredFields {
		if !utils.IsType(ind.FieldByName(name), expectedType) {
			panic(fmt.Errorf("`%s` field should be defined as `%s`", name, expectedType))
		}
	}

	var (
		sort    = reflect.Indirect(ind.FieldByName("Sort"))
		page    = reflect.Indirect(ind.FieldByName("Page"))
		perPage = reflect.Indirect(ind.FieldByName("PerPage"))
	)

	if !utils.IsEmptyValue(sort) {
		if sorts, ok := sort.Interface().([]string); ok {
			params.SetOrderBy(sorts)
		}
	}

	if !utils.IsEmptyValue(page) && !utils.IsEmptyValue(perPage) {
		params.SetPagination(int(page.Int()), int(perPage.Int()))
	}

	for i := 0; i < ind.NumField(); i++ {
		structField := ind.Type().Field(i)
		field := ind.Field(i)
		kind := field.Kind()

		if kind == reflect.Ptr && field.IsNil() {
			continue
		}

		filter := structField.Tag.Get("filter")
		if filter == "" {
			continue
		}

		value := reflect.Indirect(field).Interface()

		params.SetFilter(filter, value)
	}

	return params
}

// GetFilters return the filters map
func (p *QueryParams) GetFilters() map[string]interface{} {
	return p.filters
}

// SetFilter add a filter condition
func (p *QueryParams) SetFilter(name string, value interface{}) {
	p.filters[name] = value
}

// SetOrderBy set the order by condition
func (p *QueryParams) SetOrderBy(orderBy []string) {
	p.orderBy = orderBy
}

// SetPagination set the pagination info
func (p *QueryParams) SetPagination(page, perPage int) {
	p.page = page
	p.perPage = perPage
}

// Offset return the sql offset
func (p *QueryParams) Offset() int {
	return (p.page - 1) * p.perPage
}

// Limit return the sql limit
func (p *QueryParams) Limit() int {
	return p.perPage
}

// OrderBy return the order by exprs
func (p *QueryParams) OrderBy() []string {
	return p.orderBy
}

// String return the string representation
// nolint:errcheck
func (p *QueryParams) String() string {
	buf := &bytes.Buffer{}

	buf.WriteString("Filters:[")
	for filter, value := range p.filters {
		buf.WriteString(filter)
		buf.WriteString("=")
		buf.WriteString(fmt.Sprint(value))
		buf.WriteString(",")
	}

	if len(p.filters) > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	buf.WriteString("]")

	buf.WriteString(";OrderBy:[")
	for _, orderBy := range p.orderBy {
		buf.WriteString(orderBy)
		buf.WriteString(",")
	}

	if len(p.orderBy) > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	buf.WriteString("]")

	buf.WriteString(";Page:")
	buf.WriteString(fmt.Sprint(p.page))
	buf.WriteString(";PerPage:")
	buf.WriteString(fmt.Sprint(p.perPage))

	return buf.String()
}
