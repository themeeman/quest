package inventory

import (
	"encoding/json"
	"fmt"
	"database/sql/driver"
)

const (
	ChestDaily uint = iota
)

type Chest struct {
	Name        string
	Description string
}

type ChestsInventory map[uint]uint
type Chests map[uint]*Chest

func (c ChestsInventory) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &c)
		if err != nil {
			return err
		}
		return nil
	case string:
		err := json.Unmarshal([]byte(v), &c)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (c ChestsInventory) Value() (driver.Value, error) {
	return json.Marshal(c)
}
