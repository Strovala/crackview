package generator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Strovala/crackview/execution"
)

const (
	addToArrayTemplateCpp    = "%v.push_back(%v);"
	addToArrayTemplatePython = "%v,"
	addToArrayTemplateJava   = addToArrayTemplatePython
)

type array struct {
	*base
}

func NewArray(value interface{}) *array {
	result := &array{
		base: &base{
			Value: value,
		},
	}
	result.Resolve()
	return result
}

func (p *array) Resolve() {
	reflectType := getReflectType(p.Value)
	inputType := reflectType.String()
	p.Type = inputType[2:]
}

func (p *array) GeneratePython(name string) string {
	template := getTemplate(execution.Python, arrayTemplate)
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for i := 0; i < val.Len(); i++ {
		fmt.Fprintf(&builder, addToArrayTemplatePython, valuePython(val.Index(i), p.Type))
	}
	return fmt.Sprintf(template, name, builder.String())
}

func (p *array) GenerateJava(name string) string {
	template := getTemplate(execution.Java, arrayTemplate)
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for i := 0; i < val.Len(); i++ {
		fmt.Fprintf(&builder, addToArrayTemplateJava, value(val.Index(i), p.Type))
	}
	javaType := getJavaType(p.Type)
	return fmt.Sprintf(template, javaType, name, javaType, builder.String())
}

func (p *array) GenerateCpp(name string) string {
	template := getTemplate(execution.Cpp, arrayTemplate)
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for i := 0; i < val.Len(); i++ {
		fmt.Fprintf(&builder, addToArrayTemplateCpp, name, value(val.Index(i), p.Type))
	}
	cppType := getCppType(p.Type)
	return fmt.Sprintf(template, cppType, name, builder.String())
}
