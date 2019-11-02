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
	*baseArgument
}

func NewSet(value interface{}) *set {
	result := &set{
		baseArgument: &baseArgument{
			value: value,
		},
	}
	result.Resolve()
	return result
}

func (p *set) Resolve() {
	reflectType := getReflectType(p.value)
	inputType := reflectType.String()
	p.argType = inputType[2:]
}

func (p *set) Generate(name string, lang Language) string {
	var builder strings.Builder
	val := reflect.ValueOf(p.value)
	for i := 0; i < val.Len(); i++ {
		addToSet := lang.GenerateAddToSetTemplate(p, name, lang.Value(val.Index(i), p.Type()))
		builder.WriteString(addToSet)
	}
	return lang.GenerateSetTemplate(p, name, builder.String())
}
