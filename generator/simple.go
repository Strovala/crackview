package generator

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

func (p *simple) Generate(name string, lang Language) string {
	return lang.GenerateSimpleTemplate(p, name, lang.Value(p.value, p.Type()))
}
