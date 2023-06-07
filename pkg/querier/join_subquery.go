package querier

import (
	sq "github.com/Masterminds/squirrel"
)

func (a *APIGen) JoinSubquery() (string, []interface{}, error) {
	sub, arg, _ := a.sqlBuilder.
		Select("d.department_id").
		From("departments d").
		Join("locations l on l.location_id = d.location_id").
		Where(sq.Eq{"country_id": "US"}).
		Prefix("(").Suffix(")").
		ToSql()

	sql := a.sqlBuilder.Select("*").
		From("employees e").
		LeftJoin(sub+" d on d.department_id = e.department_id ", arg...)

	return sql.ToSql()
}

/**
select e.*
from employees e
left join
	(
		select d.department_id
			from departments d
			join locations l on l.location_id = d.location_id and l.country_id = 'US'
	) d
on d.department_id = e.employee_id
*/
