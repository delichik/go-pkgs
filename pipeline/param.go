package pipeline

import "reflect"

type Param struct {
	RemovePointer bool
	Name          string
	Value         reflect.Value
}

type ParamRequire struct {
	AddPointer bool
	Name       string
	Param      *Param
}
