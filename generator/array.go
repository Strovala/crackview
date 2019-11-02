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

type array struct {
	*argument
}

func NewArray(value interface{}) *array {
	result := &array{
		argument: &argument{
			value: value,
		},
	}
	result.ResolveType()
	return result
}

func (p *array) ResolveType() {
	inputType := reflect.TypeOf(p.value).String()
	p.argType = inputType[2:]
}

func (p *array) Generate(name string, lang Language) string {
	var builder strings.Builder
	val := reflect.ValueOf(p.value)
	for i := 0; i < val.Len(); i++ {
		addToSet := lang.GenerateAddToArrayTemplate(name, val.Index(i), p.argType)
		builder.WriteString(addToSet)
	}
	return lang.GenerateArrayTemplate(p.argType, name, builder.String())
}
