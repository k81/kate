package queryparams_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/k81/kate/queryparams"
	"github.com/stretchr/testify/require"
)

type queryReq struct {
	IDs             *[]int64 `filter:"ID__in"`
	Name            *string  `filter:"Name__icontains"`
	CreateTimeBegin *int64   `filter:"CreateTime__gte"`
	CreateTimeEnd   *int64   `filter:"CreateTime__lt"`
	UpdateTimeBegin *int64   `filter:"UpdateTime__gte"`
	UpdateTimeEnd   *int64   `filter:"UpdateTime__lt"`
	Sort            []string `valid:"sortfield(id)"`
	Page            int
	PerPage         int
}

func TestNewQueryParams(t *testing.T) {
	ids := []int64{1, 2}
	name := "zhangsan"
	createTimeBegin := int64(111)
	createTimeEnd := int64(222)
	sortBy := []string{"+id"}
	page := 1
	perPage := 10
	req := &queryReq{
		IDs:             &ids,
		Name:            &name,
		CreateTimeBegin: &createTimeBegin,
		CreateTimeEnd:   &createTimeEnd,
		Sort:            sortBy,
		Page:            page,
		PerPage:         perPage,
	}

	expectedFilters := map[string]interface{}{
		"ID__in":          ids,
		"Name__icontains": name,
		"CreateTime__gte": createTimeBegin,
		"CreateTime__lt":  createTimeEnd,
	}

	p := queryparams.NewQueryParamsFromTag(req)
	require.NotNil(t, p)
	require.Equal(t, expectedFilters, p.GetFilters())
	spew.Dump(p)
}
