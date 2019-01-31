package utility

import (
	"time"
	"reflect"
)

func TimeToTimestamp(t time.Time) string {
	return t.Format("2006-01-02T15:04:05+00:00")
}

func Contains(slice, value interface{}) (bool, int) {
	s := reflect.ValueOf(slice)
	if !(s.Kind() == reflect.Slice || s.Kind() == reflect.Array) {
		panic("Slice must be a slice!")
	}
	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(value, s.Index(i).Interface()) {
			return true, i
		}
	}
	return false, 0
}
