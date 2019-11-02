package generator

import (
	"reflect"
	"strings"
)

const (
	addToMapTemplateCpp    = "%v[%v]=%v;"
	addToMapTemplateJava   = "put(%v,%v);"
	addToMapTemplatePython = "%v:%v,"
)

type inputMap struct {
	*argument
}

func NewMap(value interface{}) *inputMap {
	result := &inputMap{
		argument: &argument{
			value: value,
		},
	}
	result.ResolveType()
	return result
}

func (p *inputMap) ResolveType() {
	inputType := reflect.TypeOf(p.value).String()
	types := strings.Split(inputType, "]")
	p.keyType = types[0][4:]
	p.valueType = types[1]
}

func (p *inputMap) Generate(name string, lang Language) string {
	var builder strings.Builder
	val := reflect.ValueOf(p.value)
	for _, k := range val.MapKeys() {
		v := val.MapIndex(k)
		addToMap := lang.GenerateAddToMapTemplate(name, k, p.keyType, v, p.valueType)
		builder.WriteString(addToMap)
	}
	return lang.GenerateMapTemplate(p.keyType, p.valueType, name, builder.String())
}
