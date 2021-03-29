package model

import (
	"fmt"
	"strings"
	"time"
)

// JSONTime : wrapper time to hold a specific format for a date
//swagger:strfmt date
type JSONTime time.Time

// AsTime : unwrap primitive time type
func (t JSONTime) AsTime() time.Time {
	return time.Time(t)
}

// NewJSONTime : wrap primitive time type to JSONTime
func NewJSONTime(time time.Time) JSONTime {
	return JSONTime(time)
}

// MarshalJSON : convert a time type to string with a specific format YYYY-MM-DD
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02"))
	return []byte(stamp), nil
}

// UnmarshalJSON : convert date string with format YYYY-MM-DD to time.Time
func (t *JSONTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	time, err := time.Parse("2006-01-02", s)
	*t = JSONTime(time)
	return
}
