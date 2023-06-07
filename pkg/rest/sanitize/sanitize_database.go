package sanitize

import (
	"fmt"
	"strings"
)

const (
	SORT_ASC  string = "ASC"
	SORT_DESC string = "DESC"
)

// filter sort direction : must ASC or DESC
func MakeOrderDirection(orderDir string) string {
	orderDir = strings.ToUpper(strings.TrimSpace(orderDir))
	if orderDir == "" {
		orderDir = SORT_ASC
	}
	if orderDir != SORT_ASC && orderDir != SORT_DESC {
		orderDir = SORT_ASC
	}
	return orderDir
}

// output with alias "schemaName.tableName t"
func MakeTableNameWithAlias(schema, table string) (tableName string, alias string) {
	tableName = MakeTableName(schema, table)
	alias = strings.ToLower(allowedInput(table))
	return fmt.Sprintf(" %s %s", tableName, alias), alias
}

func MakeTableName(schema, table string) string {
	tableName := fmt.Sprintf("%s.%s", schema, table)
	tableName = strings.ToLower(allowedInput(tableName))
	return tableName
}

// create string join ex : schema.table2.field t2 ON t1.id_table2 = t2.id
func MakeJoinString(schema, sourceTable, targetTable, sourceField, targetField, operator string) string {
	_, aliasPrimary := MakeTableNameWithAlias(schema, sourceTable)
	tableForeign, aliasForeign := MakeTableNameWithAlias(schema, targetTable)

	return fmt.Sprintf("%s ON %s.%s %s %s.%s",
		tableForeign, aliasPrimary, sourceField, operator, aliasForeign, targetField)
}

func MakeOperatorCondition(operator string) string {
	return allowedInput(operator)
}
