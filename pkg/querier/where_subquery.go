package querier

import (
	sq "github.com/Masterminds/squirrel"
)

func (a *APIGen) Subquery1() (string, []interface{}, error) {
	sub, arg, _ := a.sqlBuilder.
		Select("salary").
		From("employees").
		Where(sq.Eq{"employee_id": 163}).
		Prefix("(").Suffix(")").
		ToSql()

	sql := a.sqlBuilder.Select("*").
		From("employees").
		Where("salary > "+sub, arg...)

	return sql.ToSql()
}

/**
SELECT * FROM employees WHERE salary > (SELECT salary FROM employees WHERE employee_id = 163)
*/
