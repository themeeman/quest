package utility

import (
	"reflect"
	"strings"
)

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

func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.ToLower(s)[0:len(prefix)] == strings.ToLower(prefix)
}

func TrimPrefix(s, prefix string) string {
	if HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}
