package generator

import "reflect"

type simple struct {
	*argument
}

func NewSimple(value interface{}) *simple {
	result := &simple{
		argument: &argument{
			value: value,
		},
	}
	result.ResolveType()
	return result
}

func (p *simple) ResolveType() {
	inputType := reflect.TypeOf(p.value).String()
	p.argType = inputType
}

func (p *simple) Generate(name string, lang Language) string {
	return lang.GenerateSimpleTemplate(name, p.value, p.argType)
}
