package querier

func (a *APIGen) List() (string, []interface{}, error) {
	sql := a.sqlBuilder.
		Select("*").
		From("public.countries")
	return sql.ToSql()
}
