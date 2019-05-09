package utility

import (
	"reflect"
	"strings"
)

// Contains returns whether a value is in a slice
// If it is found, what the index of it is
// slice: []T
// value: T
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
	return false, -1
}

// ContainsFunc returns whether a value, when applied to a func, is in a slice
// If it is found, what the index of it is
// slice: []T
// value: U
// f:     func(U) T
func ContainsFunc(slice, value, fun interface{}) (bool, int) {
	validFunction := func(fun interface{}) bool {
		f := reflect.TypeOf(fun)
		return f.Kind() == reflect.Func &&
			f.NumIn() == 1 && f.In(0) == reflect.TypeOf(value) &&
			f.NumOut() == 1 && reflect.SliceOf(f.Out(0)) == reflect.TypeOf(slice)
	}
	s := reflect.ValueOf(slice)
	if !(s.Kind() == reflect.Slice || s.Kind() == reflect.Array) {
		panic("Slice must be a slice!")
	}
	f := reflect.ValueOf(slice)
	if !validFunction(fun) {
		panic("Function not in form func(U) T")
	}
	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(value, f.Call([]reflect.Value{s.Index(i)})) {
			return true, i
		}
	}
	return false, -1
}

// HasPrefix returns if a string has the given prefix, case insensitive
func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.ToLower(s)[0:len(prefix)] == strings.ToLower(prefix)
}

// TrimPrefix returns the given string without the given prefix, case insensitive
func TrimPrefix(s, prefix string) string {
	if HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}
