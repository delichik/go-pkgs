package pipeline

import "reflect"

type Param struct {
	pointerRemoved bool
	name           string
	type_          reflect.Type
	value          reflect.Value
}

type ParamRequire struct {
	needPointer bool
	name        string
	type_       reflect.Type
	param       *Param
}
