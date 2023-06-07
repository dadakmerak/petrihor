package constant

const (
	JoinTypeLeft  string = "left"
	JoinTypeRight string = "right"
	JoinTypeInner string = "inner"

	JoinType        string = "type"
	JoinSource      string = "source"
	JoinTargetTable string = "target"
	JoinSourceField string = "from"
	JoinTargetField string = "to"
	JoinOperator    string = "operator"
)

const (
	FilterTypeOr  string = "or"
	FilterTypeIf  string = "if"
	FilterTypeAnd string = "and"

	FilterCondition string = "condition"
	FilterOperator  string = "operator"
	FilterField     string = "field"
	FilterValue     string = "value"
)

// query uri
const (
	OperatorEQ           string = "eq" // =
	OperatorNotEQ        string = "ne" // != <>
	OperatorLessThan     string = "lt" // <
	OperatorGreaterThan  string = "gt" // >
	OperatorLessEqual    string = "le" // <=
	OperatorGreaterEqual string = "ge" // >=

	//special operator
	OperatorIN      string = "in"
	OperatorNotIN   string = "nin"
	OperatorLike    string = "like"
	OperatorNotLike string = "nlike"
	OperatorBetween string = "between"
)

var ToSQLOperator = map[string]string{
	OperatorEQ:           "=",
	OperatorNotEQ:        "<>",
	OperatorLessThan:     "<",
	OperatorGreaterThan:  ">",
	OperatorLessEqual:    "<=",
	OperatorGreaterEqual: ">=",
	OperatorIN:           "IN",
	OperatorNotIN:        "NOT IN",
	OperatorLike:         "LIKE",
	OperatorNotLike:      "NOT LIKE",
	OperatorBetween:      "BETWEEN",
}

const (
	AggregateQuery   string = "agg"
	AggregateGroupBy string = "groupBy"
	AggregateHaving  string = "having"
	AggregateOrder   string = "order"

	AggregateType  string = "type"
	AggregateField string = "field"
	AggregateAlias string = "as"
	AggregateCast  string = "cast"

	AggregateOrderDirection string = "dir"

	AggregateOperator string = "operator"
	AggregateValue    string = "value"
)

const (
	AggregateFieldIf       string = "if"
	AggregateFieldElse     string = "else"
	AggregateFieldThen     string = "then"
	AggregateFieldThenType string = "thenType"
)

const (
	AggregateSUM   string = "sum"
	AggregateCount string = "count"
	AggregateMIN   string = "min"
	AggregateMAX   string = "max"
	AggregateAVG   string = "avg"
)
