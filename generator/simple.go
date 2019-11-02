package generator

import (
	"fmt"

	"github.com/Strovala/crackview/execution"
)

type simple struct {
	*baseArgument
}

func NewSimple(value interface{}) *simple {
	result := &simple{
		baseArgument: &baseArgument{
			value: value,
		},
	}
	result.Resolve()
	return result
}

func (p *simple) Resolve() {
	reflectType := getReflectType(p.value)
	inputType := reflectType.String()
	p.argType = inputType
}

func (p *simple) GeneratePython(name string) string {
	template := getTemplate(execution.Python, simpleTemplate)
	return fmt.Sprintf(template, name, valuePython(p.Value, p.Type()))
}

func (p *simple) GenerateJava(name string) string {
	template := getTemplate(execution.Java, simpleTemplate)
	javaType := getJavaType(p.Type())
	return fmt.Sprintf(template, javaType, name, value(p.Value, p.Type()))
}

func (p *simple) GenerateCpp(name string) string {
	template := getTemplate(execution.Cpp, simpleTemplate)
	cppType := getCppType(p.Type())
	return fmt.Sprintf(template, cppType, name, value(p.value, p.Type()))
}

func (p *simple) Generate(name string, lang Language) string {
	return lang.GenerateSimpleTemplate(p, name, lang.Value(p.value, p.Type()))
}
