package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Dummy struct{}

func (d Dummy) ConfigureEntity(e *EntityConfigurator) {
	// TODO implement me
	panic("implement me")
}

func TestSelect(t *testing.T) {
	t.Run("only select * from Table", func(t *testing.T) {
		s := NewQueryBuilder[Dummy]()
		s.Table("users").SetSelect()
		str, args, err := s.ToSql()
		assert.NoError(t, err)
		assert.Empty(t, args)
		assert.Equal(t, "SELECT * FROM users", str)
	})
	t.Run("select with whereClause", func(t *testing.T) {
		s := NewQueryBuilder[Dummy]()

		s.Table("users").SetDialect(Dialects.MySQL).
			Where("age", 10).
			AndWhere("age", "<", 10).
			Where("name", "Amirreza").
			OrWhere("age", GT, 11).
			SetSelect()

		str, args, err := s.ToSql()
		assert.NoError(t, err)
		assert.EqualValues(t, []interface{}{10, 10, "Amirreza", 11}, args)
		assert.Equal(t, "SELECT * FROM users WHERE age = ? AND age < ? AND name = ? OR age > ?", str)
	})
	t.Run("select with order by", func(t *testing.T) {
		s := NewQueryBuilder[Dummy]().Table("users").OrderBy("created_at", ASC).OrderBy("updated_at", DESC)
		str, args, err := s.ToSql()
		assert.NoError(t, err)
		assert.Empty(t, args)
		assert.Equal(t, "SELECT * FROM users ORDER BY created_at ASC,updated_at DESC", str)
	})

	t.Run("select with group by", func(t *testing.T) {
		s := NewQueryBuilder[Dummy]().Table("users").GroupBy("created_at", "updated_at")
		str, args, err := s.ToSql()
		assert.NoError(t, err)
		assert.Empty(t, args)
		assert.Equal(t, "SELECT * FROM users GROUP BY created_at,updated_at", str)
	})

	t.Run("Select with limit", func(t *testing.T) {
		s := NewQueryBuilder[Dummy]().Table("users").Limit(10)
		str, args, err := s.ToSql()
		assert.NoError(t, err)
		assert.Empty(t, args)
		assert.Equal(t, "SELECT * FROM users LIMIT 10", str)
	})

	t.Run("Select with offset", func(t *testing.T) {
		s := NewQueryBuilder[Dummy]().Table("users").Offset(10)
		str, args, err := s.ToSql()
		assert.NoError(t, err)
		assert.Empty(t, args)
		assert.Equal(t, "SELECT * FROM users OFFSET 10", str)
	})

	t.Run("select with join", func(t *testing.T) {
		s := NewQueryBuilder[Dummy]().Table("users").Select("id", "name").RightJoin("addresses", "users.id", "addresses.user_id")
		str, args, err := s.ToSql()
		assert.NoError(t, err)
		assert.Empty(t, args)
		assert.Equal(t, `SELECT id,name FROM users RIGHT JOIN addresses ON users.id = addresses.user_id`, str)
	})

	t.Run("select with multiple joins", func(t *testing.T) {
		s := NewQueryBuilder[Dummy]().Table("users").
			Select("id", "name").
			RightJoin("addresses", "users.id", "addresses.user_id").
			LeftJoin("user_credits", "users.id", "user_credits.user_id")
		sql, args, err := s.ToSql()
		assert.NoError(t, err)
		assert.Empty(t, args)
		assert.Equal(t, `SELECT id,name FROM users RIGHT JOIN addresses ON users.id = addresses.user_id LEFT JOIN user_credits ON users.id = user_credits.user_id`, sql)
	})

	t.Run("select with subquery", func(t *testing.T) {
		s := NewQueryBuilder[Dummy]().SetDialect(Dialects.MySQL)
		s.FromQuery(NewQueryBuilder[Dummy]().Table("users").Where("age", "<", 10))
		sql, args, err := s.ToSql()
		assert.NoError(t, err)
		assert.EqualValues(t, []interface{}{10}, args)
		assert.Equal(t, `SELECT * FROM (SELECT * FROM users WHERE age < ? )`, sql)

	})

	t.Run("select with inner join", func(t *testing.T) {
		s := NewQueryBuilder[Dummy]().Table("users").Select("id", "name").InnerJoin("addresses", "users.id", "addresses.user_id")
		str, args, err := s.ToSql()
		assert.NoError(t, err)
		assert.Empty(t, args)
		assert.Equal(t, `SELECT id,name FROM users INNER JOIN addresses ON users.id = addresses.user_id`, str)
	})

	t.Run("select with join", func(t *testing.T) {
		s := NewQueryBuilder[Dummy]().Table("users").Select("id", "name").Join("addresses", "users.id", "addresses.user_id")
		str, args, err := s.ToSql()
		assert.NoError(t, err)
		assert.Empty(t, args)
		assert.Equal(t, `SELECT id,name FROM users INNER JOIN addresses ON users.id = addresses.user_id`, str)
	})

	t.Run("select with full outer join", func(t *testing.T) {
		s := NewQueryBuilder[Dummy]().Table("users").Select("id", "name").FullOuterJoin("addresses", "users.id", "addresses.user_id")
		str, args, err := s.ToSql()
		assert.NoError(t, err)
		assert.Empty(t, args)
		assert.Equal(t, `SELECT id,name FROM users FULL OUTER JOIN addresses ON users.id = addresses.user_id`, str)
	})
	t.Run("raw where", func(t *testing.T) {
		sql, args, err :=
			NewQueryBuilder[Dummy]().
				SetDialect(Dialects.MySQL).
				Table("users").
				Where(Raw("id = ?", 1)).
				AndWhere(Raw("age < ?", 10)).
				SetSelect().
				ToSql()

		assert.NoError(t, err)
		assert.EqualValues(t, []interface{}{1, 10}, args)
		assert.Equal(t, `SELECT * FROM users WHERE id = ? AND age < ?`, sql)
	})
	t.Run("no sql type matched", func(t *testing.T) {
		sql, args, err := NewQueryBuilder[Dummy]().ToSql()
		assert.Error(t, err)
		assert.Empty(t, args)
		assert.Empty(t, sql)
	})

	t.Run("raw where in", func(t *testing.T) {
		sql, args, err :=
			NewQueryBuilder[Dummy]().
				SetDialect(Dialects.MySQL).
				Table("users").
				WhereIn("id", Raw("SELECT user_id FROM user_books WHERE book_id = ?", 10)).
				SetSelect().
				ToSql()

		assert.NoError(t, err)
		assert.EqualValues(t, []interface{}{10}, args)
		assert.Equal(t, `SELECT * FROM users WHERE id IN (SELECT user_id FROM user_books WHERE book_id = ?)`, sql)
	})
	t.Run("where in", func(t *testing.T) {
		sql, args, err :=
			NewQueryBuilder[Dummy]().
				SetDialect(Dialects.MySQL).
				Table("users").
				WhereIn("id", 1, 2, 3, 4, 5, 6).
				SetSelect().
				ToSql()

		assert.NoError(t, err)
		assert.EqualValues(t, []interface{}{1, 2, 3, 4, 5, 6}, args)
		assert.Equal(t, `SELECT * FROM users WHERE id IN (?,?,?,?,?,?)`, sql)

	})
}
func TestUpdate(t *testing.T) {
	t.Run("update no whereClause", func(t *testing.T) {
		u := NewQueryBuilder[Dummy]().Table("users").Set("name", "amirreza").SetDialect(Dialects.MySQL)
		sql, args, err := u.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, `UPDATE users SET name=?`, sql)
		assert.Equal(t, []interface{}{"amirreza"}, args)
	})
	t.Run("update with whereClause", func(t *testing.T) {
		u := NewQueryBuilder[Dummy]().Table("users").Set("name", "amirreza").Where("age", "<", 18).SetDialect(Dialects.MySQL)
		sql, args, err := u.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, `UPDATE users SET name=? WHERE age < ?`, sql)
		assert.Equal(t, []interface{}{"amirreza", 18}, args)

	})
}
func TestDelete(t *testing.T) {
	t.Run("delete without whereClause", func(t *testing.T) {
		d := NewQueryBuilder[Dummy]().Table("users").SetDelete()
		sql, args, err := d.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, `DELETE FROM users`, sql)
		assert.Empty(t, args)
	})
	t.Run("delete with whereClause", func(t *testing.T) {
		d := NewQueryBuilder[Dummy]().Table("users").SetDialect(Dialects.MySQL).Where("created_at", ">", "2012-01-10").SetDelete()
		sql, args, err := d.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, `DELETE FROM users WHERE created_at > ?`, sql)
		assert.EqualValues(t, []interface{}{"2012-01-10"}, args)
	})
}
