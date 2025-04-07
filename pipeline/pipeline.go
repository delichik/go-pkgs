package pipeline

import (
	"reflect"
)

type _ParamProvider struct {
	Function *_Function
	Param    *_Param
}

type Pipeline struct {
	paramProviders map[string]*_ParamProvider
	functions      []*_Function
	invokers       []*_Function
	ordered        []*_Function
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		paramProviders: make(map[string]*_ParamProvider),
	}
}

func (p *Pipeline) Provide(fun any) *Pipeline {
	fo := readFunction(fun)
	for _, param := range fo.Out {
		p.paramProviders[param.Name] = &_ParamProvider{
			Function: fo,
			Param:    param,
		}
	}
	p.functions = append(p.functions, fo)

	return p
}

func (p *Pipeline) Invoke(fun any) *Pipeline {
	fo := readFunction(fun)
	for _, param := range fo.Out {
		p.paramProviders[param.Name] = &_ParamProvider{
			Function: fo,
			Param:    param,
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
		params := make([]reflect.Value, 0, len(fo.In))
		for _, pr := range fo.In {
			pv := pr.Param.Value
			if pr.NeedPointer {
				pv = pv.Addr()
			}
			params = append(params, pv)
		}
		outputs := fo.Call(params...)
		for i, p := range fo.Out {
			if p.PointerRemoved {
				p.Value = outputs[i].Elem()
			} else {
				t := reflect.New(outputs[i].Type())
				t.Elem().Set(outputs[i])
				p.Value = t.Elem()
			}
			if !p.Value.IsValid() {
				panic("nil or invalid value for " + p.Name)
			}
		}
	}
}

func (p *Pipeline) call(fo *_Function) {
	if fo.Loaded {
		return
	}

	if fo.InLoad {
		panic(fo.Name + " in loop call")
	}

	fo.InLoad = true
	for _, param := range fo.In {
		pp, ok := p.paramProviders[param.Name]
		if !ok && param.Type.Kind() == reflect.Interface {
			for _, pp2 := range p.paramProviders {
				if pp2.Param.Type.Kind() == reflect.Interface {
					continue
				}
				if reflect.PointerTo(pp2.Param.Type).Implements(param.Type) {
					if pp != nil {
						panic("confused implements of " + param.Name)
					}
					pp = pp2
					param.NeedPointer = true
				} else if pp2.Param.Type.Implements(param.Type) {
					if pp != nil {
						panic("confused implements of " + param.Name)
					}
					pp = pp2
				}
			}
		}
		if pp == nil {
			panic("no provider for " + param.Name)
		}
		param.Param = pp.Param
		p.call(pp.Function)
	}
	p.ordered = append(p.ordered, fo)
	fo.Loaded = true
}
