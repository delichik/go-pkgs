package pipeline

import (
	"reflect"
)

type _ParamProvider struct {
	function *Function
	param    *Param
}

type Pipeline struct {
	paramProviders map[string]*_ParamProvider
	functions      []*Function
	invokers       []*Function
	ordered        []*Function
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		paramProviders: make(map[string]*_ParamProvider),
	}
}

func (p *Pipeline) Provide(fun any) *Pipeline {
	var fo *Function
	switch fun.(type) {
	case *Function:
		fo = fun.(*Function)
	default:
		fo = readFunction(fun)
	}
	for _, param := range fo.out {
		p.paramProviders[param.name] = &_ParamProvider{
			function: fo,
			param:    param,
		}
	}
	p.functions = append(p.functions, fo)

	return p
}

func (p *Pipeline) Invoke(fun any) *Pipeline {
	fo := readFunction(fun)
	for _, param := range fo.out {
		p.paramProviders[param.name] = &_ParamProvider{
			function: fo,
			param:    param,
		}
	}
	p.invokers = append(p.invokers, fo)
	p.functions = append(p.functions, fo)
	return p
}

func (p *Pipeline) Prepare() *Pipeline {
	for _, fun := range p.invokers {
		p.call(fun)
	}
	return p
}

func (p *Pipeline) Run() {
	for _, fo := range p.ordered {
		params := make([]reflect.Value, 0, len(fo.in))
		for _, pr := range fo.in {
			pv := pr.param.value
			if pr.needPointer {
				if pv.CanAddr() {
					pv = pv.Addr()
				} else {
					t := reflect.New(pv.Type())
					t.Elem().Set(pv)
					pv = t
				}
			}
			params = append(params, pv)
		}
		outputs := fo.Call(params...)
		for i, p := range fo.out {
			if p.pointerRemoved {
				p.value = outputs[i].Elem()
			} else {
				t := reflect.New(outputs[i].Type())
				t.Elem().Set(outputs[i])
				p.value = t.Elem()
			}
			if !p.value.IsValid() {
				panic("nil or invalid value for " + p.name)
			}
		}
	}
}

func (p *Pipeline) call(fo *Function) {
	if fo.loaded {
		return
	}

	if fo.inLoad {
		panic(fo.name + " in loop call")
	}

	fo.inLoad = true
	for _, param := range fo.in {
		pp, ok := p.paramProviders[param.name]
		if !ok && param.type_.Kind() == reflect.Interface {
			for _, pp2 := range p.paramProviders {
				if pp2.param.type_.Kind() == reflect.Interface {
					continue
				}
				if reflect.PointerTo(pp2.param.type_).Implements(param.type_) {
					if pp != nil {
						panic("confused implements of " + param.name)
					}
					pp = pp2
					param.needPointer = true
				} else if pp2.param.type_.Implements(param.type_) {
					if pp != nil {
						panic("confused implements of " + param.name)
					}
					pp = pp2
				}
			}
		}
		if pp == nil {
			panic("no provider for " + param.name)
		}
		param.param = pp.param
		p.call(pp.function)
	}
	p.ordered = append(p.ordered, fo)
	fo.loaded = true
}
