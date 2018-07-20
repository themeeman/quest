package inventory

import (
	"encoding/json"
	"fmt"
	"database/sql/driver"
)

type Chest struct {
	Name string
}

type Chests []*Chest

func (c Chests) Scan(val interface{}) error {
	switch v := val.(type){
	case []byte:
		json.Unmarshal(v, c)
		return nil
	case string:
		json.Unmarshal([]byte(v), c)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (c Chests) Value() (driver.Value, error) {
	return json.Marshal(c)
}
