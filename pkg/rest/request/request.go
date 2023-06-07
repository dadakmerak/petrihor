package request

import (
	"fmt"
	"strings"

	"github.com/dadakmerak/petrihor/pkg/rest/constant"
	"github.com/dadakmerak/petrihor/pkg/rest/sanitize"
	"github.com/dadakmerak/petrihor/pkg/shared"
)

// Generate SQL : field to select as JSON ouput, and join clause
func MakeJoin(schema, sourceTable string, joinMap shared.Map, fields map[string][]string) (column string, tables map[string][]string) {
	column = ""
	tables = make(map[string][]string, 0)

	for _, dataJoin := range joinMap {
		datas := dataJoin.([]interface{})
		for _, dataJoin := range datas {
			dataJoin := dataJoin.(map[string]interface{})

			typeJoin := getStringFromMap(dataJoin, constant.JoinType, "")
			targetTable := getStringFromMap(dataJoin, constant.JoinTargetTable, "")
			sourceField := getStringFromMap(dataJoin, constant.JoinSourceField, fmt.Sprintf("%sId", targetTable))
			targetField := getStringFromMap(dataJoin, constant.JoinTargetField, "id")
			operator := getStringFromMap(dataJoin, constant.JoinOperator, "=")
			sourceTable := getStringFromMap(dataJoin, constant.JoinSource, sourceTable)

			targetTable = sanitize.ToSnakeCase(targetTable)
			sourceField = sanitize.ToSnakeCase(sourceField)
			targetField = sanitize.ToSnakeCase(targetField)
			sourceTable = sanitize.ToSnakeCase(sourceTable)

			_, aliasTargetTable := sanitize.MakeTableNameWithAlias(schema, targetTable)

			if len(fields) > 0 {
				if len(fields[sanitize.ToSnakeCase(targetTable)]) > 0 {
					c := makeFieldString(fields[sanitize.ToSnakeCase(targetTable)])
					if c == "*" {
						column = fmt.Sprintf("%s, json_build_array(%s.*) %s", column, aliasTargetTable, targetTable)
					} else {
						column = fmt.Sprintf("%s, json_build_array(jsonb_build_object(%s)) %s", column, c, targetTable)
					}
				}
			} else {
				column = fmt.Sprintf("%s, json_build_array(%s.*) %s", column, aliasTargetTable, targetTable)
			}

			tables[typeJoin] = append(tables[typeJoin],
				sanitize.MakeJoinString(schema, sourceTable, targetTable, sourceField, targetField, operator))
		}
	}

	return column, tables
}

// Get value from a map, by a key then cast to string
func getStringFromMap(m map[string]interface{}, key string, defaultVal string) string {
	if val, ok := m[key]; ok {
		return val.(string)
	}
	return defaultVal
}

// make string for nested json, if query URI field not empty
//
// output : 'title', post.title
func makeFieldString(fields []string) string {
	c := ""
	for _, field := range fields {
		f := strings.Split(field, ".")
		if f[1] == "*" {
			return "*"
		}
		c = fmt.Sprintf(`%s '%s', %s,`, c, f[1], field) // previus column string, 'fieldName', value ,
	}
	return strings.Trim(c, ",")
}

// Extract Joined table from query uri filter
func ExtractTableFromQueryURIJoin(queryURI string) ([]string, error) {
	tables := make([]string, 0)
	joins, err := sanitize.QueryURIToMap(queryURI)

	if err != nil {
		return nil, err
	}
	for _, values := range joins {
		datas := values.([]interface{})
		for _, join := range datas {
			dataJoin := join.(map[string]interface{})
			if val, ok := dataJoin[constant.JoinTargetTable]; ok {
				tables = append(tables, sanitize.ToCamelCase(val.(string)))
			}
		}
	}
	return tables, nil
}

// Build SQL WHERE
func MakeCondition(sourceTable string, dataCondition map[string]interface{}) (condition string, args []interface{}) {
	isFirst := true

	for _, datas := range dataCondition {
		datas := datas.([]interface{})
		for _, dataFilter := range datas {
			dataFilter := dataFilter.(map[string]interface{})
			var valueField interface{}
			conditionType := getStringFromMap(dataFilter, constant.FilterCondition, "")
			operator := getStringFromMap(dataFilter, constant.FilterOperator, "")
			field := getStringFromMap(dataFilter, constant.FilterField, "")

			if val, ok := dataFilter[constant.FilterValue]; ok {
				valueField = val
			}

			if isFirst {
				isFirst = false
				conditionType = ""
			}

			field = columnWithTable(sourceTable, field)
			whereStr, values := conditionString(conditionType, operator, field, valueField)
			condition = condition + whereStr
			args = append(args, values...)
		}
	}
	return condition, args
}

// output : AND field  > value
func conditionString(conditionType, operator, field string, values interface{}) (condition string, args []interface{}) {
	condition = ""
	conditionType = sanitize.MakeOperatorCondition(conditionType)

	switch operator {
	case constant.OperatorIN, constant.OperatorNotIN:
		newValue := values.([]interface{})
		condition = fmt.Sprintf("%s (%s %s (%s)) ",
			conditionType,
			field,
			constant.ToSQLOperator[operator],
			placeholders(len(newValue)),
		)
		args = append(args, newValue...)
	case constant.OperatorLike, constant.OperatorNotLike:
		value := fmt.Sprint("%", strings.ToLower(values.(string)), "%")
		field := fmt.Sprintf("LOWER(%s)", field)
		condition = fmt.Sprintf("%s (%s %s ?) ",
			conditionType,
			field,
			constant.ToSQLOperator[operator])
		args = append(args, value)
	case constant.OperatorBetween:
		newValue := values.([]interface{})
		condition = fmt.Sprintf("%s (%s %s ? AND ?) ",
			conditionType,
			field,
			constant.ToSQLOperator[operator],
		)
		args = append(args, newValue...)
	default:
		condition = fmt.Sprintf("%s (%s %s ?) ",
			conditionType,
			field,
			constant.ToSQLOperator[operator])
		args = append(args, values)
	}
	return condition, args
}

func placeholders(count int) string {
	if count < 1 {
		return ""
	}

	return strings.Repeat(",?", count)[1:]
}

// Generate SQL : for aggreate function
//
// Column : SUM(table.id) as sum_table_id
//
// GroupBy: table2.name (for column and GROUP By)
//
// Order: table2.name ASC
//
// Having: SUM(table.id) > 1000
func MakeAggreate(schema, sourceTable string, queryAgg shared.Map) (column string, groupBy []string, order []string, having map[string]interface{}, expression map[string][]interface{}) {
	column = ""
	expression = make(map[string][]interface{}, 0)
	if groupBys, ok := queryAgg[constant.AggregateGroupBy]; ok {
		for _, field := range groupBys.([]interface{}) {
			field := field.(string)
			field = columnWithTable(sourceTable, field)
			groupBy = append(groupBy, field+",")
		}
	}

	if agg, ok := queryAgg[constant.AggregateQuery]; ok {
		for _, dataAgg := range agg.([]interface{}) {
			dataAgg := dataAgg.(map[string]interface{})
			if _, ok := dataAgg[constant.AggregateField].(map[string]interface{}); ok {
				exprs_query, exprs_arg := aggreateExpression(dataAgg)
				expression[exprs_query] = exprs_arg
				continue
			}
			aggType := dataAgg[constant.AggregateType].(string)
			aggField := dataAgg[constant.AggregateField].(string)
			field := ""

			aggField = columnWithTable(sourceTable, aggField)
			alias := strings.ReplaceAll(aggField, ".", "_")
			alias = fmt.Sprintf("%s_%s", aggType, alias)
			if as, ok := dataAgg[constant.AggregateAlias].(string); ok {
				alias = sanitize.ToSnakeCase(as)
			}

			cast := ""
			if c, ok := dataAgg[constant.AggregateCast].(string); ok {
				cast = fmt.Sprintf("::%s", c)
			}

			switch aggType {
			case constant.AggregateSUM:
				field = fmt.Sprintf(" SUM(%s)%s AS %s,", aggField, cast, alias)
			case constant.AggregateCount:
				field = fmt.Sprintf(" COUNT(%s)%s AS %s,", aggField, cast, alias)
			case constant.AggregateMIN:
				field = fmt.Sprintf(" MIN(%s)%s AS %s,", aggField, cast, alias)
			case constant.AggregateMAX:
				field = fmt.Sprintf(" MAX(%s)%s AS %s,", aggField, cast, alias)
			case constant.AggregateAVG:
				field = fmt.Sprintf(" AVG(%s)%s AS %s,", aggField, cast, alias)
			}
			column = column + field
		}
	}

	if orders, ok := queryAgg[constant.AggregateOrder]; ok { // Order by aggregate value
		for _, dataOrder := range orders.([]interface{}) {
			dataOrder := dataOrder.(map[string]interface{})
			orderField := dataOrder[constant.AggregateField].(string)
			orderDir := dataOrder[constant.AggregateOrderDirection].(string)
			orderField = columnWithTable(sourceTable, orderField)

			if orderType, ok := dataOrder[constant.AggregateType].(string); ok {
				switch orderType {
				case constant.AggregateSUM:
					orderField = fmt.Sprintf("SUM(%s)", orderField)
				case constant.AggregateCount:
					orderField = fmt.Sprintf("COUNT(%s)", orderField)
				case constant.AggregateMIN:
					orderField = fmt.Sprintf("MIN(%s)", orderField)
				case constant.AggregateMAX:
					orderField = fmt.Sprintf("MAX(%s)", orderField)
				case constant.AggregateAVG:
					orderField = fmt.Sprintf("AVG(%s)", orderField)
				}
			}
			orderDir = fmt.Sprintf("%s %s", orderField, orderDir)
			order = append(order, orderDir)
		}
	}

	if havings, ok := queryAgg[constant.AggregateHaving]; ok {
		mapHaving := make(map[string]interface{})
		for _, dataHaving := range havings.([]interface{}) {
			dataHaving := dataHaving.(map[string]interface{})

			havingType := dataHaving[constant.AggregateType].(string)
			havingField := dataHaving[constant.AggregateField].(string)
			havingOperator := dataHaving[constant.AggregateOperator].(string)
			havingValue := dataHaving[constant.AggregateValue]
			field := ""

			havingField = columnWithTable(sourceTable, havingField)

			switch havingType {
			case constant.AggregateSUM:
				field = fmt.Sprintf("SUM(%s) %s ?", havingField, constant.ToSQLOperator[havingOperator])
			case constant.AggregateCount:
				field = fmt.Sprintf("COUNT(%s) %s ?", havingField, constant.ToSQLOperator[havingOperator])
			case constant.AggregateMIN:
				field = fmt.Sprintf("MIN(%s) %s ?", havingField, constant.ToSQLOperator[havingOperator])
			case constant.AggregateMAX:
				field = fmt.Sprintf("MAX(%s) %s ?", havingField, constant.ToSQLOperator[havingOperator])
			case constant.AggregateAVG:
				field = fmt.Sprintf("AVG(%s) %s ?", havingField, constant.ToSQLOperator[havingOperator])
			}
			mapHaving[field] = havingValue
		}
		having = mapHaving
	}
	column = strings.Trim(column, ",")
	return column, groupBy, order, having, expression
}

// Create field with table name, to avoid ambiguous column name
//
// return table_name.field_name
func columnWithTable(table, column string) string {
	if strings.Contains(column, ".") {
		splittedField := strings.Split(column, ".")
		return fmt.Sprintf("%s.%s", sanitize.ToSnakeCase(splittedField[0]), sanitize.ToSnakeCase(splittedField[1]))
	}
	return fmt.Sprintf("%s.%s", sanitize.ToSnakeCase(table), sanitize.ToSnakeCase(column))
}

// Organize field into a map[string][string]
//
// Example Output:
//
// post : ["id", "title", "content"]
//
// category : ["id", "name"]
func MappingColumnName(table string, fields []string) map[string][]string {
	field := make(map[string][]string, 0)
	for _, f := range fields {
		if f == "" {
			continue
		}
		if strings.Contains(f, ".") {
			splittedField := strings.Split(f, ".")

			fieldTable := fmt.Sprintf("%s.%s", sanitize.ToSnakeCase(splittedField[0]), sanitize.ToSnakeCase(splittedField[1]))
			field[splittedField[0]] = append(field[splittedField[0]], fieldTable)
		} else {
			fieldTable := fmt.Sprintf("%s.%s", sanitize.ToSnakeCase(table), sanitize.ToSnakeCase(f))
			field[table] = append(field[table], fieldTable)
		}
	}
	return field
}

func aggreateExpression(dataAggregate map[string]interface{}) (string, []interface{}) {
	args := make([]interface{}, 0)
	query, alias, caseQuery := "", "", "CASE "
	field := dataAggregate[constant.AggregateField].(map[string]interface{})
	aggType := dataAggregate[constant.AggregateType].(string)

	cases := field[constant.AggregateFieldIf].([]interface{})

	for _, caseEx := range cases {
		expression := caseEx.(map[string]interface{})
		expField := getStringFromMap(expression, constant.AggregateField, "")
		expOpt := getStringFromMap(expression, constant.AggregateOperator, "eq")
		expOpt = constant.ToSQLOperator[expOpt]
		expVal := expression[constant.AggregateValue]

		expThen := expression[constant.AggregateFieldThen]
		expThenType := getStringFromMap(expression, constant.AggregateFieldThenType, "")

		caseQuery = fmt.Sprintf("%s WHEN %s %s ? THEN", caseQuery, expField, expOpt)
		args = append(args, expVal)
		if expThenType == "text" {
			caseQuery = fmt.Sprintf("%s ?", caseQuery)
			args = append(args, expThen)
		}
		if expThenType == "col" {
			caseQuery = fmt.Sprintf("%s %s", caseQuery, expThen)
		}

		if alias == "" {
			alias = strings.ReplaceAll(expField, ".", "_")
			alias = fmt.Sprintf("%s_%s", aggType, alias)
		}
	}

	if elses, ok := field[constant.AggregateFieldElse].(map[string]interface{}); ok {
		elseThen := elses[constant.AggregateFieldThen]
		elseThenType := elses[constant.AggregateFieldThenType].(string)
		if elseThenType == "text" {
			caseQuery = fmt.Sprintf("%s ELSE ?", caseQuery)
			args = append(args, elseThen)
		}
		if elseThenType == "col" {
			caseQuery = fmt.Sprintf("%s ELSE %s", caseQuery, elseThen)
		}
	}

	caseQuery = fmt.Sprintf("%s END", caseQuery)
	if as, ok := dataAggregate[constant.AggregateAlias].(string); ok {
		alias = as
	}

	alias = sanitize.ToSnakeCase(alias)

	switch aggType {
	case constant.AggregateSUM:
		query = fmt.Sprintf("SUM(%s) %s", caseQuery, alias)
	case constant.AggregateCount:
		query = fmt.Sprintf("COUNT(%s) %s", caseQuery, alias)
	case constant.AggregateMIN:
		query = fmt.Sprintf("MIN(%s) %s", caseQuery, alias)
	case constant.AggregateMAX:
		query = fmt.Sprintf("MAX(%s) %s", caseQuery, alias)
	case constant.AggregateAVG:
		query = fmt.Sprintf("AVG(%s) %s", caseQuery, alias)
	}

	return query, args
}

func ExtractTableAndCols(table string, fields []string, join, aggregate shared.Map) ([]string, []string) {
	field := make([]string, 0)
	for _, f := range fields {
		if f == "" {
			continue
		}
		if strings.Contains(f, ".") {
			splittedField := strings.Split(f, ".")
			field = append(field, sanitize.ToSnakeCase(splittedField[1]))
		} else {
			field = append(field, sanitize.ToSnakeCase(f))
		}
	}
	tables := make([]string, 0)
	for _, dataJoin := range join {
		datas := dataJoin.([]interface{})
		for _, dataJoin := range datas {
			dataJoin := dataJoin.(map[string]interface{})
			targetTable := getStringFromMap(dataJoin, constant.JoinTargetTable, "")
			targetTable = sanitize.ToSnakeCase(targetTable)
			tables = append(tables, targetTable)
		}
	}
	return field, tables
}
