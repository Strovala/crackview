package generator

import (
	"reflect"
	"strings"
)

const (
	addToSetTemplateCpp    = "%v.insert(%v);"
	addToSetTemplateJava   = "add(%v);"
	addToSetTemplatePython = "%v,"
)

type set struct {
	*argument
}

func NewSet(value interface{}) *set {
	result := &set{
		argument: &argument{
			value: value,
		},
	}
	result.ResolveType()
	return result
}

func (p *set) ResolveType() {
	inputType := reflect.TypeOf(p.value).String()
	p.argType = inputType[2:]
}

func (p *set) Generate(name string, lang Language) string {
	var builder strings.Builder
	val := reflect.ValueOf(p.value)
	for i := 0; i < val.Len(); i++ {
		addToSet := lang.GenerateAddToSetTemplate(name, val.Index(i), p.argType)
		builder.WriteString(addToSet)
	}
	return lang.GenerateSetTemplate(p.argType, name, builder.String())
}
