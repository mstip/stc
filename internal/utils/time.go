package utils

import (
	"encoding/json"
	"time"
)

func TimeToStr(t time.Time) string {
	tb, _ := json.Marshal(t)
	return string(tb)
}

func StrToTime(tstr string) (time.Time, error) {
	var t time.Time
	err := json.Unmarshal([]byte(tstr), &t)
	return t, err
}
