package pipeline

import (
	"reflect"
)

type Function struct {
	in             []*ParamRequire
	out            []*Param
	name           string
	inLoad         bool
	loaded         bool
	dynamic        bool
	dynamicMapping []reflect.Value

	frv reflect.Value
}

func (f *Function) Call(params ...reflect.Value) (outputs []reflect.Value) {
	if f.dynamic {
		for i, param := range params {
			dm := f.dynamicMapping[i]
			if dm.CanSet() {
				dm.Set(param)
			} else {
				setPrivate(dm, param)
			}
		}
		for _, v := range f.out {
			outputs = append(outputs, v.value)
		}
		return outputs
	} else {
		return f.frv.Call(params)
	}
}

func readFunction(fun any) *Function {
	frv := reflect.ValueOf(fun)
	frt := frv.Type()

	fo := &Function{
		name: getName(frt),
		frv:  frv,
	}

	for i := 0; i < frt.NumIn(); i++ {
		pt := frt.In(i)
		pr := &ParamRequire{}
		if pt.Kind() == reflect.Pointer {
			pt = pt.Elem()
			pr.needPointer = true
		}

		pr.name = getName(pt)
		pr.type_ = pt
		fo.in = append(fo.in, pr)
	}

	for i := 0; i < frt.NumOut(); i++ {
		pt := frt.Out(i)
		param := &Param{}
		if pt.Kind() == reflect.Pointer {
			pt = pt.Elem()
			param.pointerRemoved = true
		}
		param.name = pt.PkgPath() + "." + pt.Name()
		param.type_ = pt
		fo.out = append(fo.out, param)
	}

	return fo
}
