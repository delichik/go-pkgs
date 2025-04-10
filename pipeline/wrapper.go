package pipeline

import "reflect"

type fieldValue struct {
	index       int
	needPointer bool
	value       reflect.Value
}

// ProvideValue 直接提供一个变量
func (p *Pipeline) ProvideValue(t any) *Pipeline {
	rv := reflect.ValueOf(t)
	rt := rv.Type()
	param := &Param{
		pointerRemoved: false,
	}
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
		param.pointerRemoved = true
	}
	param.value = rv
	param.type_ = rt
	param.name = getName(rt)

	return p.Provide(&Function{
		in:      []*ParamRequire{},
		out:     []*Param{param},
		name:    "value_provider." + getName(rt),
		dynamic: true,
	})
}

// ProvideStruct 直接提供一个结构体，这个结构体中的所有非匿名参数将自动被可用填充
func (p *Pipeline) ProvideStruct(model any) *Pipeline {
	rt := reflect.TypeOf(model)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		panic("AutoProvider[T]: T must be a struct")
	}

	prv := reflect.New(rt)
	rv := prv.Elem()

	in := []*ParamRequire{}
	fieldValues := []fieldValue{}
	dynamicMapping := []reflect.Value{}
	for i := range rt.NumField() {
		rf := rt.Field(i)
		if rf.Anonymous {
			continue
		}
		rft := rf.Type
		pr := &ParamRequire{}

		if rft.Kind() == reflect.Pointer {
			rft = rft.Elem()
			pr.needPointer = true
		}

		pr.name = getName(rft)
		pr.type_ = rft
		dynamicMapping = append(dynamicMapping, rv.Field(i))
		fieldValues = append(fieldValues, fieldValue{
			index:       i,
			needPointer: pr.needPointer,
			value:       rv.Field(i),
		})
		in = append(in, pr)
	}

	return p.Provide(&Function{
		in: in,
		out: []*Param{{
			pointerRemoved: true,
			name:           getName(rt),
			type_:          rt,
			value:          prv,
		}},
		name:           "auto_provider." + getName(rt),
		dynamic:        true,
		dynamicMapping: dynamicMapping,
	})
}
