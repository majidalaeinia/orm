package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {

	t.Run("simple insert for psql", func(t *testing.T) {
		sql, args := newInsert().
			Table("users").
			Into("name", "password").
			Values("$1", "$2").
			WithArgs("amirreza", "password").
			Build()

		assert.Equal(t, []interface{}{"amirreza", "password"}, args)
		assert.Equal(t, "INSERT INTO users (name,password) VALUES ($1,$2)", sql)
	})

	t.Run("simple insert for mysql", func(t *testing.T) {
		sql, args := newInsert().
			Table("users").
			Into("name", "password").
			Values("?", "?").
			WithArgs("amirreza", "password").
			Build()

		assert.Equal(t, []interface{}{"amirreza", "password"}, args)
		assert.Equal(t, "INSERT INTO users (name,password) VALUES (?,?)", sql)
	})

}