package pipeline

import "reflect"

type _Param struct {
	PointerRemoved bool
	Name           string
	Type           reflect.Type
	Value          reflect.Value
}

type ParamRequire struct {
	NeedPointer bool
	Name        string
	Type        reflect.Type
	Param       *_Param
}
