package pipeline

import "reflect"

type _Function struct {
	In     []*ParamRequire
	Out    []*_Param
	Name   string
	InLoad bool
	Loaded bool
	frv    reflect.Value
}

func (f *_Function) Call(params ...reflect.Value) (outputs []reflect.Value) {
	return f.frv.Call(params)
}

func readFunction(fun any) *_Function {
	frv := reflect.ValueOf(fun)
	frt := frv.Type()

	fo := &_Function{
		Name: frt.PkgPath() + "." + frt.Name(),
		frv:  frv,
	}

	for i := 0; i < frt.NumIn(); i++ {
		pt := frt.In(i)
		pr := &ParamRequire{}
		if pt.Kind() == reflect.Pointer {
			pt = pt.Elem()
			pr.NeedPointer = true
		}

		pr.Name = pt.PkgPath() + "." + pt.Name()
		pr.Type = pt
		fo.In = append(fo.In, pr)
	}

	for i := 0; i < frt.NumOut(); i++ {
		pt := frt.Out(i)
		param := &_Param{}
		if pt.Kind() == reflect.Pointer {
			pt = pt.Elem()
			param.PointerRemoved = true
		}
		param.Name = pt.PkgPath() + "." + pt.Name()
		param.Type = pt
		fo.Out = append(fo.Out, param)
	}

	return fo
}
