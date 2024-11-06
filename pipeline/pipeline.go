package pipeline

import (
	"reflect"
)

type ParamProvider struct {
	Function *Function
	Param    *Param
}

type Pipeline struct {
	paramProviders map[string]*ParamProvider
	functions      []*Function
	invokers       []*Function
	ordered        []*Function
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		paramProviders: make(map[string]*ParamProvider),
	}
}

func (p *Pipeline) Provide(fun any) *Pipeline {
	fo := readFunction(fun)
	for _, param := range fo.Out {
		p.paramProviders[param.Name] = &ParamProvider{
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
		p.paramProviders[param.Name] = &ParamProvider{
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
			if pr.AddPointer {
				params = append(params, pr.Param.Value.Addr())
			} else {
				params = append(params, pr.Param.Value)
			}
		}
		outputs := fo.Call(params...)
		for i, p := range fo.Out {
			if p.RemovePointer {
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

func (p *Pipeline) call(fo *Function) {
	if fo.Loaded {
		return
	}

	if fo.InLoad {
		panic(fo.Name + " in loop call")
	}

	fo.InLoad = true
	for _, param := range fo.In {
		pp, ok := p.paramProviders[param.Name]
		if !ok {
			panic("no provider for " + param.Name)
		}
		param.Param = pp.Param
		p.call(pp.Function)
	}
	p.ordered = append(p.ordered, fo)
	fo.Loaded = true
}
