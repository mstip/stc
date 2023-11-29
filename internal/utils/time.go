package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type DBTime struct {
	T time.Time
}

func NewDBTimeNow() DBTime {
	return DBTime{T: time.Now()}
}

func (t DBTime) Value() (driver.Value, error) {
	return t.T.String(), nil
}

func (t DBTime) String() string {
	return TimeToStr(t.T)
}

func (t *DBTime) Scan(src any) error {
	var source []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	case nil:
		return nil
	default:
		return fmt.Errorf("incompatible type for DBTime")
	}

	var err error
	t.T, err = StrToTime(string(source))
	return err
}

func TimeToStr(t time.Time) string {
	tb, _ := json.Marshal(t)
	return string(tb)
}

func StrToTime(tstr string) (time.Time, error) {
	var t time.Time
	err := json.Unmarshal([]byte(tstr), &t)
	return t, err
}
