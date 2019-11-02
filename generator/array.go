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
	*baseArgument
}

func NewArray(value interface{}) *array {
	result := &array{
		baseArgument: &baseArgument{
			value: value,
		},
	}
	result.Resolve()
	return result
}

func (p *array) Resolve() {
	reflectType := getReflectType(p.value)
	inputType := reflectType.String()
	p.argType = inputType[2:]
}

func (p *array) Generate(name string, lang Language) string {
	var builder strings.Builder
	val := reflect.ValueOf(p.value)
	for i := 0; i < val.Len(); i++ {
		addToSet := lang.GenerateAddToArrayTemplate(p, name, lang.Value(val.Index(i), p.Type()))
		builder.WriteString(addToSet)
	}
	return lang.GenerateArrayTemplate(p, name, builder.String())
}
