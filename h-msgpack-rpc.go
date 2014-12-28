package main

import (
	"errors"
	"reflect"
)

type MPResolver struct {
	*Config
}

func (r *MPResolver) Resolve(name string, arguments []reflect.Value) (val reflect.Value, eerr error) {
	switch name {

	case "geo":
		val = reflect.ValueOf(
			func(ip string) map[string]interface{} {
				return r.VBase.Load().(*GeoBase).Find(ip).ToMap()
			},
		)

	default:
		eerr = errors.New("Not implemented")

	}
	return
}
