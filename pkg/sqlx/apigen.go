package sqlx

import (
	"fmt"

	"github.com/dadakmerak/petrihor/pkg/database"
	"github.com/dadakmerak/petrihor/pkg/maps"
	"github.com/dadakmerak/petrihor/pkg/querier"
)

type Repository struct {
	db    *database.DB
	query *querier.APIGen
}

func NewRepository(db *database.DB, query *querier.APIGen) *Repository {
	return &Repository{
		db:    db,
		query: query,
	}
}

func (r *Repository) List() (maps.ListResponses, error) {
	q, a, err := r.query.JoinSubquery()
	fmt.Println(q)
	fmt.Println(a...)
	//
	//
	//
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Conn.Queryx(q, a...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	responses := make(maps.ListResponses, 0)
	cols, _ := rows.Columns()
	columns := make([]interface{}, len(cols))
	for rows.Next() {
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}
		response := make(maps.Response, 1)
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			response[colName] = *val
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func (r *Repository) Detail() (maps.Response, error) {
	response := make(maps.Response, 0)
	return response, nil
}

func (r *Repository) Create() error {
	return nil
}

func (r *Repository) Update() error {
	return nil
}

func (r *Repository) Delete() error {
	return nil
}
