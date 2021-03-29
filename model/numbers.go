package model

import (
	"fmt"
	"strconv"
	"strings"
)

// Number : Object to hold a number with 2 decimal positions to represent Money or Percentage fields
// swagger:type number
type Number float64

// AsFloat : Get Number as primitive float
func (n Number) AsFloat() float64 {
	return float64(n)
}

// NewNumber : Create a Number object from float64 type
func NewNumber(n float64) Number {
	return Number(n)
}

// MarshalJSON : Marshall a Number object to a string
func (t Number) MarshalJSON() ([]byte, error) {
	text := fmt.Sprintf("%.2f", t)
	return []byte(text), nil
}

// UnmarshalJSON : Unmarshal a String with number to a Number
func (n *Number) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	f, err := strconv.ParseFloat(s, 64)
	*n = NewNumber(f)
	return err
}
