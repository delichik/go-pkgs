package pipeline

import "reflect"

type Function struct {
	In     []*ParamRequire
	Out    []*Param
	Name   string
	InLoad bool
	Loaded bool
	frv    reflect.Value
}

func (f *Function) Call(params ...reflect.Value) (outputs []reflect.Value) {
	return f.frv.Call(params)
}

func readFunction(fun any) *Function {
	frv := reflect.ValueOf(fun)
	frt := frv.Type()

	fo := &Function{
		Name: frt.PkgPath() + "." + frt.Name(),
		frv:  frv,
	}

	for i := 0; i < frt.NumIn(); i++ {
		pt := frt.In(i)
		pr := &ParamRequire{}
		if pt.Kind() == reflect.Pointer {
			pt = pt.Elem()
			pr.AddPointer = true
		}
		pr.Name = pt.PkgPath() + "." + pt.Name()
		fo.In = append(fo.In, pr)
	}

	for i := 0; i < frt.NumOut(); i++ {
		pt := frt.Out(i)
		param := &Param{}
		if pt.Kind() == reflect.Pointer {
			pt = pt.Elem()
			param.RemovePointer = true
		}
		param.Name = pt.PkgPath() + "." + pt.Name()
		fo.Out = append(fo.Out, param)
	}

	return fo
}
