package utils

import (
	"testing"
	"time"
)

func TestTimeToStr(t *testing.T) {
	str := TimeToStr(time.Now())
	if str == "" {
		t.Fatal("time str is empty")
	}
}

func TestStrToTime(t *testing.T) {
	now := time.Now()
	tt, err := StrToTime(TimeToStr(now))
	if err != nil {
		t.Fatal(err)
	}

	if !now.Equal(tt) {
		t.Fatalf("StrToTime: expect %#v, got %#v", now, tt)
	}
}
