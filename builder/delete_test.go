package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {

	t.Run("simple delete equality of id", func(t *testing.T) {
		s, err := NewDelete("users").
			Where("id", "=", "$1").
			Stmt().
			SQL()
		assert.NoError(t, err)
		assert.Equal(t, `DELETE FROM users WHERE id = $1`, s)
	})
}