package convert

import (
	"strconv"
	"strings"
	"time"
)

//ToFloat converts string to float64. It returns 0 if conversion is impossible
func ToFloat(s string) float32 {
	f, _ := strconv.ParseFloat(strings.Replace(s, ",", ".", -1), 32)
	return float32(f)
}

//ToFloat64 converts string to float64. It returns 0 if conversion is impossible
func ToFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(strings.Replace(s, ",", ".", -1), 64)
	return f
}

//ToInt converts string to int. It returns 0 if conversion is impossible
func ToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

//ToInt64 converts string to int64. It returns 0 if conversion is impossible
func ToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 0, 64)
	return i
}

func BoolP(v bool) *bool {
	return &v
}

func Float64P(v float64) *float64 {
	return &v
}

func FloatP(v float32) *float32 {
	return &v
}

func Int64P(v int64) *int64 {
	return &v
}

func IntP(v int) *int {
	return &v
}

func StrP(s string) *string {
	return &s
}

func TimeP(t time.Time) *time.Time {
	return &t
}
