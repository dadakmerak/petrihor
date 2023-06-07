package maps

import (
	"database/sql/driver"
	"encoding/json"
)

type Response map[string]interface{}

func (r *Response) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *Response) Scan(value interface{}) error {
	b, err := value.([]byte)
	if !err {
		return nil
	}
	return json.Unmarshal(b, &r)
}

type ListResponses []Response
