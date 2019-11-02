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
	*baseArgument
}

func NewMap(value interface{}) *inputMap {
	result := &inputMap{
		baseArgument: &baseArgument{
			value: value,
		},
	}
	result.Resolve()
	return result
}

func (p *inputMap) Resolve() {
	reflectType := getReflectType(p.value)
	inputType := reflectType.String()
	types := strings.Split(inputType, "]")
	p.keyType = types[0][4:]
	p.valueType = types[1]
}

func (p *inputMap) Generate(name string, lang Language) string {
	var builder strings.Builder
	val := reflect.ValueOf(p.value)
	for _, k := range val.MapKeys() {
		v := val.MapIndex(k)
		addToMap := lang.GenerateAddToMapTemplate(p, name, lang.Value(k, p.KeyType()), lang.Value(v, p.ValueType()))
		builder.WriteString(addToMap)
	}
	return lang.GenerateMapTemplate(p, name, builder.String())
}
