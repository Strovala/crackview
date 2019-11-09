package generator

import (
	"reflect"
	"strings"
)

const (
	addToArrayTemplateCpp    = "%v.push_back(%v);"
	addToArrayTemplatePython = "%v,"
	addToArrayTemplateJava   = addToArrayTemplatePython
)

type Array struct {
	*argument
}

func NewArray(value interface{}) *Array {
	result := &Array{
		argument: &argument{
			value: value,
		},
	}
	result.ResolveType()
	return result
}

func (p *Array) ResolveType() {
	inputType := reflect.TypeOf(p.value).String()
	p.argType = inputType[2:]
}

func (p *Array) Generate(name string, lang Language) string {
	var builder strings.Builder
	val := reflect.ValueOf(p.value)
	for i := 0; i < val.Len(); i++ {
		addToSet := lang.GenerateAddToArrayTemplate(name, val.Index(i), p.argType)
		builder.WriteString(addToSet)
	}
	return lang.GenerateArrayTemplate(p.argType, name, builder.String())
}
