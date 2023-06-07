package querier

import (
	sq "github.com/Masterminds/squirrel"
)

type APIGen struct {
	sqlBuilder sq.StatementBuilderType
}

func NewAPIGen() *APIGen {
	sqlBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &APIGen{sqlBuilder: sqlBuilder}
}
