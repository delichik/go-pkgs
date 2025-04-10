package pipeline

import (
	"reflect"
	"unsafe"
)

func getName(rt reflect.Type) string {
	return rt.PkgPath() + "." + rt.Name()
}

func setPrivate(dst reflect.Value, src reflect.Value) {
	const flagRO uintptr = 1<<5 | 1<<6
	flagField := reflect.ValueOf(&dst).Elem().FieldByName("flag")
	flagPtr := (*uintptr)(unsafe.Pointer(flagField.UnsafeAddr()))
	old := *flagPtr
	*flagPtr &= ^(flagRO)
	dst.Set(src)
	*flagPtr = old
}
