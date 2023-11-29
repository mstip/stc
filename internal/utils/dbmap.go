package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type DBMap struct {
	M map[string]string
}

func NewDBMap(m map[string]string) DBMap {
	return DBMap{M: m}
}

func (m DBMap) Value() (driver.Value, error) {
	str, err := json.Marshal(m.M)
	if err != nil {
		return nil, err
	}
	return str, nil
}

func (m DBMap) String() string {
	str, err := json.Marshal(m.M)
	if err != nil {
		return ""
	}
	return string(str)
}

func (m *DBMap) Scan(src any) error {
	var source []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	case nil:
		return nil
	default:
		return fmt.Errorf("incompatible type for DBMap")
	}
	var mm map[string]string
	err := json.Unmarshal(source, &mm)
	m.M = mm
	return err
}
