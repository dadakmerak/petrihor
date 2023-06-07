package shared

import (
	"database/sql/driver"
	"encoding/json"
)

type Map map[string]interface{}

func (r Map) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *Map) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, &r)
}

func (r Map) IsEmpty() bool {
	return len(r) == 0
}

func (r Map) Set(key string, value interface{}) {
	r[key] = value
}

func (r Map) GetMap(key string) map[string]interface{} {
	for k, v := range r {
		if k == key {
			if value, ok := v.(map[string]interface{}); ok {
				return value
			}
		}
	}
	return nil
}

func (r Map) GetString(key string) string {
	if value, ok := r[key]; ok {
		if v, ok := value.(string); ok {
			return v
		}
	}
	return ""
}
