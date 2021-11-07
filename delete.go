package orm

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type DeleteStmt struct {
	table string
	where string
}

func (q *DeleteStmt) Where(parts ...string) *DeleteStmt {
	if q.where == "" {
		q.where = fmt.Sprintf("%s", strings.Join(parts, " "))
		return q
	}
	q.where = fmt.Sprintf("%s AND %s", q.where, strings.Join(parts, " "))
	return q
}

func (q *DeleteStmt) OrWhere(parts ...string) *DeleteStmt {
	q.where = fmt.Sprintf("%s OR %s", q.where, strings.Join(parts, " "))
	return q
}

func (q *DeleteStmt) AndWhere(parts ...string) *DeleteStmt {
	return q.Where(parts...)
}

func (d *DeleteStmt) SQL() (string, error) {
	return fmt.Sprintf("DELETE FROM %s WHERE %s", d.table, d.where), nil
}

func (d *DeleteStmt) ExecContext(ctx context.Context, db *sql.DB, args ...interface{}) (sql.Result, error) {
	s, err := d.SQL()
	if err != nil {
		return nil, err
	}
	return exec(context.Background(), db, s, args)
}
func (d *DeleteStmt) Exec(db *sql.DB, args ...interface{}) (sql.Result, error) {
	query, err := d.SQL()
	if err != nil {
		return nil, err
	}
	return exec(context.Background(), db, query, args)

}

func NewDelete(table string) *DeleteStmt {
	return &DeleteStmt{table: table}
}