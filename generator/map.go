package generator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Strovala/crackview/execution"
)

const (
	addToMapTemplateCpp    = "%v[%v]=%v;"
	addToMapTemplateJava   = "put(%v,%v);"
	addToMapTemplatePython = "%v:%v,"
)

type inputMap struct {
	*base
	KeyType   string
	ValueType string
}

func NewMap(value interface{}) *inputMap {
	result := &inputMap{
		base: &base{
			Value: value,
		},
	}
	result.Resolve()
	return result
}

func (p *inputMap) Resolve() {
	reflectType := getReflectType(p.Value)
	inputType := reflectType.String()
	types := strings.Split(inputType, "]")
	p.KeyType = types[0][4:]
	p.ValueType = types[1]
}

func (p *inputMap) GeneratePython(name string) string {
	template := getTemplate(execution.Python, mapTemplate)
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for _, k := range val.MapKeys() {
		v := val.MapIndex(k)
		fmt.Fprintf(&builder, addToMapTemplatePython, valuePython(k, p.KeyType), valuePython(v, p.ValueType))
	}
	return fmt.Sprintf(template, name, builder.String())
}

func (p *inputMap) GenerateJava(name string) string {
	template := getTemplate(execution.Java, mapTemplate)
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for _, k := range val.MapKeys() {
		v := val.MapIndex(k)
		fmt.Fprintf(&builder, addToMapTemplateJava, valuePython(k, p.KeyType), value(v, p.ValueType))
	}
	javaKeyType := javaWrapperClass(getJavaType(p.KeyType))
	javaValueType := javaWrapperClass(getJavaType(p.ValueType))
	return fmt.Sprintf(template, javaKeyType, javaValueType, name, javaKeyType, javaValueType, builder.String())
}

func (p *inputMap) GenerateCpp(name string) string {
	template := getTemplate(execution.Cpp, mapTemplate)
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for _, k := range val.MapKeys() {
		v := val.MapIndex(k)
		fmt.Fprintf(&builder, addToMapTemplateCpp, name, valuePython(k, p.KeyType), value(v, p.ValueType))
	}
	cppKeyType := getCppType(p.KeyType)
	cppValueType := getCppType(p.ValueType)
	return fmt.Sprintf(template, cppKeyType, cppValueType, name, builder.String())
}
