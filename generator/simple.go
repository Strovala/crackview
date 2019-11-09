package generator

import "reflect"

type Simple struct {
	*argument
}

func NewSimple(value interface{}) *Simple {
	result := &Simple{
		argument: &argument{
			value: value,
		},
	}
	result.ResolveType()
	return result
}

func (p *Simple) ResolveType() {
	inputType := reflect.TypeOf(p.value).String()
	p.argType = inputType
}

func (p *Simple) Generate(name string, lang Language) string {
	return lang.GenerateSimpleTemplate(name, p.value, p.argType)
}
