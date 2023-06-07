package querier

import (
	sq "github.com/Masterminds/squirrel"
)

func (a *APIGen) SelectSubquery() (string, []interface{}, error) {
	sub := a.sqlBuilder.
		Select("salary").
		From("employees")

	sql := a.sqlBuilder.Select("*").
		FromSelect(sub, "aliasSub").
		PlaceholderFormat(sq.Dollar)
	return sql.ToSql()
}

/**
SELECT * FROM ( SELECT salary FROM employees ) AS aliasSub
*/
