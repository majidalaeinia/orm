package builder

import (
	"fmt"
	"strings"
)

type whereHelpers struct {
	Like    func(column string, pattern string) string
	In      func(column string, values ...string) string
	Between func(column string, lower string, higher string) string
	Equal   func(column, value string) string
	Less    func(column, value string) string
	More    func(column, value string) string
	EqualID func(value string) string
	And     func(conds ...string) string
	Or      func(conds ...string) string
	Not     func(cond ...string) string
}

var WhereHelpers = &whereHelpers{
	Like:    like,
	In:      in,
	Between: between,
	Not:     not,
	Equal:   equal,
	EqualID: func(value string) string {
		return equal("id", value)
	},
	Less: less,
	More: more,
	Or:   or,
	And:  and,
}

func less(column, value string) string {
	return fmt.Sprintf("%s < %s", column, value)
}

func more(column, value string) string {
	return fmt.Sprintf("%s > %s", column, value)
}

func equal(column, value string) string {
	return fmt.Sprintf("%s = %s", column, value)
}
func like(column string, pattern string) string {
	return fmt.Sprintf("%s LIKE %s", column, pattern)
}

func in(column string, values ...string) string {
	return fmt.Sprintf("%s IN (%s)", column, strings.Join(values, ", "))
}

func between(column string, lower string, higher string) string {
	return fmt.Sprintf("%s BETWEEN %s AND %s", column, lower, higher)
}

func not(cond ...string) string {
	return fmt.Sprintf("NOT %s", strings.Join(cond, " "))
}

func and(cond ...string) string {
	return fmt.Sprintf("%s", strings.Join(cond, " AND "))
}

func or(cond ...string) string {
	return fmt.Sprintf("%s", strings.Join(cond, " OR "))
}