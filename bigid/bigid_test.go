package bigid

import (
	"testing"

	"github.com/k81/kate/orm/sqlbuilder"
	"github.com/stretchr/testify/assert"
)

func TestBigID(t *testing.T) {
	builder := sqlbuilder.NewSelectBuilder()
	id := Fake(10)
	builder.Select("a", "b").From("c").Where(
		builder.E("id", id),
	)

	query, args := builder.Build()

	assert.Equal(t, "SELECT a, b FROM c WHERE id = ?", query)
	assert.Equal(t, []interface{}{Fake(10)}, args)
}
