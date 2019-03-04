package structures

import (
	"github.com/fatih/structs"
	"reflect"
)

type Option struct {
	Name string
	Type string
}

func GetOptions() map[string]*Option {
	g := new(Guild)
	t := reflect.TypeOf(*g)
	options := make(map[string]*Option)
	for _, s := range structs.Names(*g) {
		field, _ := t.FieldByName(s)
		t, ok := field.Tag.Lookup("type")
		if ok {
			name, _ := field.Tag.Lookup("db")
			options[s] = &Option{
				Name: name,
				Type: t,
			}
		}
	}
	return options
}
