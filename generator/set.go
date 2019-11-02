package generator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Strovala/crackview/execution"
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

func (p *set) GeneratePython(name string) string {
	template := getTemplate(execution.Python, setTemplate)
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for i := 0; i < val.Len(); i++ {
		fmt.Fprintf(&builder, addToSetTemplatePython, valuePython(val.Index(i), p.Type()))
	}
	return fmt.Sprintf(template, name, builder.String())
}

func (p *set) GenerateJava(name string) string {
	template := getTemplate(execution.Java, setTemplate)
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for i := 0; i < val.Len(); i++ {
		fmt.Fprintf(&builder, addToSetTemplateJava, value(val.Index(i), p.Type()))
	}
	javaType := javaWrapperClass(getJavaType(p.Type()))
	return fmt.Sprintf(template, javaType, name, javaType, builder.String())
}

func (p *set) GenerateCpp(name string) string {
	template := getTemplate(execution.Cpp, setTemplate)
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for i := 0; i < val.Len(); i++ {
		fmt.Fprintf(&builder, addToSetTemplateCpp, name, value(val.Index(i), p.Type()))
	}
	cppType := getCppType(p.Type())
	return fmt.Sprintf(template, cppType, name, builder.String())
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
