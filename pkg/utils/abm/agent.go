package abm

import (
	"reflect"
)

type Agent interface {
	Run()
	IsAlive() bool
}

func CopyAgent(src Agent) Agent {
	if src == nil {
		return nil
	}
	typ := reflect.TypeOf(src)
	val := reflect.ValueOf(src)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	elem := reflect.New(typ).Elem()
	elem.Set(val)
	return elem.Addr().Interface().(Agent)
}
