package models

import (
	"fmt"
	"strings"
	"time"
)

type JSONTime time.Time


func (t JSONTime)MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

func (t *JSONTime) UnmarshalJSON(buf []byte) (err error) {
	tt, err := time.Parse("2006-01-02 15:04:05", strings.Trim(string(buf), `"`))
	if err != nil {
		return
	}
	*t = JSONTime(tt)
	return
}
