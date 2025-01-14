package orm

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID   int64
	Name string
	Timestamps
}

func (u User) ConfigureEntity(e *EntityConfigurator) {
	e.Table("users")
}

type Address struct {
	ID   int
	Path string
}

func TestBind(t *testing.T) {
	t.Run("single result", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		mock.
			ExpectQuery("SELECT .* FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
				AddRow(1, "amirreza", sql.NullTime{Time: time.Now(), Valid: true}, sql.NullTime{Time: time.Now(), Valid: true}, sql.NullTime{}))
		rows, err := db.Query(`SELECT * FROM users`)
		assert.NoError(t, err)

		u := &User{}
		md := schemaOfHeavyReflectionStuff(u)
		err = newBinder[User](md).bind(rows, u)
		assert.NoError(t, err)

		assert.Equal(t, "amirreza", u.Name)
	})

	t.Run("multi result", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		mock.
			ExpectQuery("SELECT .* FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "amirreza").AddRow(2, "milad"))

		rows, err := db.Query(`SELECT * FROM users`)
		assert.NoError(t, err)

		md := schemaOfHeavyReflectionStuff(&User{})
		var users []*User
		err = newBinder[User](md).bind(rows, &users)
		assert.NoError(t, err)

		assert.Equal(t, "amirreza", users[0].Name)
		assert.Equal(t, "milad", users[1].Name)
	})
}

func TestBindMap(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.
		ExpectQuery("SELECT .* FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "amirreza", sql.NullTime{Time: time.Now(), Valid: true}, sql.NullTime{Time: time.Now(), Valid: true}, sql.NullTime{}))
	rows, err := db.Query(`SELECT * FROM users`)
	assert.NoError(t, err)

	ms := bindToMap(rows)

	assert.NotEmpty(t, ms)

	assert.Len(t, ms, 1)
}
