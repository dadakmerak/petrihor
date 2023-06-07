package querier

func (a *APIGen) WhereSubquery2() (string, []interface{}, error) {
	sub, arg, _ := a.sqlBuilder.
		Select("department_id").
		From("departments d").
		Join("locations l ON l.location_id = d.location_id").
		Where("l.country_id = ?", "US").
		Prefix("(").Suffix(")").
		ToSql()

	sql := a.sqlBuilder.Select("*").
		From("employees").
		Where("department_id in "+sub, arg...)

	return sql.ToSql()
}

/**
select e.*
from employees e
where e.department_id in
	(
		select d.department_id
			from departments d
			join locations l on l.location_id = d.location_id and l.country_id = 'US'
	)
*/
