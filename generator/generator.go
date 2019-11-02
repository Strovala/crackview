package generator

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/Strovala/crackview/execution"
)

const (
	templatesPath    = "generator/templates"
	arrayTemplate    = "array.txt"
	setTemplate      = "set.txt"
	simpleTemplate   = "simple.txt"
	mapTemplate      = "map.txt"
	solutionTemplate = "solution.txt"
)

const (
	intType        = "int"
	stringType     = "string"
	boolType       = "bool"
	floatTypeSmall = "float32"
	floatTypeBig   = "float64"
	floatType      = floatTypeBig
)

const (
	javaIntWrapperType = "Integer"
)

const (
	cppFloatType  = "float"
	cppDoubleType = "double"
)

var cppLookUpType = map[string]string{
	floatType:      cppDoubleType,
	floatTypeSmall: cppFloatType,
}

const (
	javaStringType = "String"
	javaBoolType   = "boolean"
	javaFloatType  = cppFloatType
	javaDoubleType = cppDoubleType
)

var javaLookUpType = map[string]string{
	stringType:     javaStringType,
	boolType:       javaBoolType,
	floatType:      javaDoubleType,
	floatTypeSmall: javaFloatType,
}

var javaLookUpWrapperClass = map[string]string{
	intType: javaIntWrapperType,
}

func capitalizeFirstLetter(value string) string {
	return strings.Title(strings.ToLower(value))
}

func getCppType(valueType string) string {
	javaType, ok := cppLookUpType[valueType]
	if !ok {
		javaType = valueType
	}
	return javaType
}

func javaWrapperClass(valueType string) string {
	javaType, ok := javaLookUpWrapperClass[valueType]
	if !ok {
		javaType = capitalizeFirstLetter(valueType)
	}
	return javaType
}

func getJavaType(valueType string) string {
	javaType, ok := javaLookUpType[valueType]
	if !ok {
		javaType = valueType
	}
	return javaType
}

func valuePython(value interface{}, valueType string) string {
	if valueType == stringType {
		return fmt.Sprintf("\"%v\"", value)
	} else if valueType == boolType {
		return capitalizeFirstLetter(fmt.Sprintf("%v", value))
	}
	return fmt.Sprintf("%v", value)
}

func value(value interface{}, valueType string) string {
	if valueType == stringType {
		return fmt.Sprintf("\"%v\"", value)
	}
	return fmt.Sprintf("%v", value)
}

func getTemplate(lang, file string) string {
	dat, _ := ioutil.ReadFile(fmt.Sprintf("%v/%v/%v", templatesPath, lang, file))
	return string(dat)
}

func getReflectType(value interface{}) reflect.Type {
	return reflect.TypeOf(value)
}

type Argument interface {
	GeneratePython(name string) string
	GenerateJava(name string) string
	GenerateCpp(name string) string
	Generate(name string, lang Language) string

	Type() string
	KeyType() string
	ValueType() string
	Value() interface{}
}

type typeResolver interface {
	Resolve()
}

type baseArgument struct {
	argType   string
	keyType   string
	valueType string
	value     interface{}
}

func (b *baseArgument) Generate(name string, lang Language) string {
	return "Not Implemented"
}

func (b *baseArgument) Type() string {
	return b.argType
}

func (b *baseArgument) KeyType() string {
	return b.keyType
}

func (b *baseArgument) ValueType() string {
	return b.valueType
}

func (b *baseArgument) Value() interface{} {
	return b.value
}

// func Generate(args []Argument, lang string) string {
// 	template := getTemplate(lang, solutionTemplate)
// 	var argsInit strings.Builder
// 	var argsPass strings.Builder
// 	for i, arg := range args {
// 		argName := fmt.Sprintf("arg_%v", i)
// 		var argInit string
// 		switch lang {
// 		case execution.Python:
// 			argInit = arg.GeneratePython(argName)
// 		case execution.Java:
// 			argInit = arg.GenerateJava(argName)
// 		case execution.Cpp:
// 			argInit = arg.GenerateCpp(argName)
// 		}
// 		fmt.Fprintf(&argsInit, "%v\n", argInit)
// 		if i != len(args)-1 {
// 			fmt.Fprintf(&argsPass, "%v,", argName)
// 		} else {
// 			fmt.Fprintf(&argsPass, "%v", argName)
// 		}
// 	}
// 	return fmt.Sprintf(template, argsInit.String(), argsPass.String())
// }

// TODO: Maybe something like this
type Language interface {
	Name() string
	Value(val interface{}, valType string) string
	GenerateSimpleTemplate(arg Argument, name, value string) string
	GenerateSetTemplate(arg Argument, name, value string) string
	GenerateAddToSetTemplate(arg Argument, name string, value string) string
	GenerateArrayTemplate(arg Argument, name, value string) string
	GenerateAddToArrayTemplate(arg Argument, name string, value string) string
	GenerateMapTemplate(arg Argument, name, value string) string
	GenerateAddToMapTemplate(arg Argument, name string, keyValue string, Value string) string
}

type baseLang struct {
	name               string
	addToSetTemplate   string
	addToArrayTemplate string
	addToMapTemplate   string
}

func (l *baseLang) Name() string {
	return l.name
}

func (l *baseLang) GenerateSimpleTemplate(arg Argument, name, value string) string {
	return "Not Implemented"
}

func (l *baseLang) GenerateArrayTemplate(arg Argument, name, value string) string {
	return "Not Implemented"
}

func (l *baseLang) GenerateMapTemplate(arg Argument, name, value string) string {
	return "Not Implemented"
}

func (l *baseLang) GenerateAddToArrayTemplate(arg Argument, name string, value string) string {
	return "Not Implemented"
}

func (l *baseLang) GenerateAddToMapTemplate(arg Argument, name string, keyValue string, Value string) string {
	return "Not Implemented"
}

func (l *baseLang) GenerateAddToSetTemplate(arg Argument, name string, value string) string {
	return "Not Implemented"
}

type python struct {
	*baseLang
}

func (l *python) Value(val interface{}, valType string) string {
	return valuePython(val, valType)
}

func (l *python) GenerateSimpleTemplate(arg Argument, name, value string) string {
	template := getTemplate(l.name, simpleTemplate)
	return fmt.Sprintf(template, name, value)
}

func (l *python) GenerateArrayTemplate(arg Argument, name, value string) string {
	template := getTemplate(l.name, arrayTemplate)
	return fmt.Sprintf(template, name, value)
}

func (l *python) GenerateAddToArrayTemplate(arg Argument, name string, value string) string {
	return fmt.Sprintf(l.addToArrayTemplate, value)
}

func (l *python) GenerateSetTemplate(arg Argument, name, value string) string {
	template := getTemplate(l.name, setTemplate)
	return fmt.Sprintf(template, name, value)
}

func (l *python) GenerateAddToSetTemplate(arg Argument, name string, value string) string {
	return fmt.Sprintf(l.addToSetTemplate, value)
}

func (l *python) GenerateMapTemplate(arg Argument, name, value string) string {
	template := getTemplate(l.name, mapTemplate)
	return fmt.Sprintf(template, name, value)
}

func (l *python) GenerateAddToMapTemplate(arg Argument, name string, keyValue string, value string) string {
	return fmt.Sprintf(l.addToMapTemplate, keyValue, value)
}

type java struct {
	*baseLang
}

func (l *java) Value(val interface{}, valType string) string {
	return value(val, valType)
}

func (l *java) GenerateSimpleTemplate(arg Argument, name, value string) string {
	template := getTemplate(l.name, simpleTemplate)
	javaType := getJavaType(arg.Type())
	return fmt.Sprintf(template, javaType, name, value)
}

func (l *java) GenerateArrayTemplate(arg Argument, name, value string) string {
	template := getTemplate(l.name, arrayTemplate)
	javaType := getJavaType(arg.Type())
	return fmt.Sprintf(template, javaType, name, javaType, value)
}

func (l *java) GenerateAddToArrayTemplate(arg Argument, name string, value string) string {
	return fmt.Sprintf(l.addToArrayTemplate, value)
}

func (l *java) GenerateSetTemplate(arg Argument, name, value string) string {
	template := getTemplate(l.name, setTemplate)
	javaType := javaWrapperClass(getJavaType(arg.Type()))
	return fmt.Sprintf(template, javaType, name, javaType, value)
}

func (l *java) GenerateAddToSetTemplate(arg Argument, name string, value string) string {
	return fmt.Sprintf(l.addToSetTemplate, value)
}

func (l *java) GenerateMapTemplate(arg Argument, name, value string) string {
	template := getTemplate(l.name, mapTemplate)
	javaKeyType := javaWrapperClass(getJavaType(arg.KeyType()))
	javaValueType := javaWrapperClass(getJavaType(arg.ValueType()))
	return fmt.Sprintf(template, javaKeyType, javaValueType, name, javaKeyType, javaValueType, value)
}

func (l *java) GenerateAddToMapTemplate(arg Argument, name string, keyValue string, value string) string {
	return fmt.Sprintf(l.addToMapTemplate, keyValue, value)
}

type cpp struct {
	*baseLang
}

func (l *cpp) Value(val interface{}, valType string) string {
	return value(val, valType)
}

func (l *cpp) GenerateSimpleTemplate(arg Argument, name, value string) string {
	template := getTemplate(l.name, simpleTemplate)
	cppType := getCppType(arg.Type())
	return fmt.Sprintf(template, cppType, name, value)
}

func (l *cpp) GenerateArrayTemplate(arg Argument, name, value string) string {
	template := getTemplate(l.name, arrayTemplate)
	cppType := getCppType(arg.Type())
	return fmt.Sprintf(template, cppType, name, value)
}

func (l *cpp) GenerateAddToArrayTemplate(arg Argument, name string, value string) string {
	return fmt.Sprintf(l.addToArrayTemplate, name, value)
}

func (l *cpp) GenerateSetTemplate(arg Argument, name, value string) string {
	template := getTemplate(l.name, setTemplate)
	cppType := getCppType(arg.Type())
	return fmt.Sprintf(template, cppType, name, value)
}

func (l *cpp) GenerateAddToSetTemplate(arg Argument, name string, value string) string {
	return fmt.Sprintf(l.addToSetTemplate, name, value)
}

func (l *cpp) GenerateMapTemplate(arg Argument, name, value string) string {
	template := getTemplate(l.name, mapTemplate)
	cppKeyType := getCppType(arg.KeyType())
	cppValueType := getCppType(arg.ValueType())
	return fmt.Sprintf(template, cppKeyType, cppValueType, name, value)
}

func (l *cpp) GenerateAddToMapTemplate(arg Argument, name string, keyValue string, value string) string {
	return fmt.Sprintf(l.addToMapTemplate, name, keyValue, value)
}

var (
	pythonObj = &python{baseLang: &baseLang{
		name:               execution.Python,
		addToSetTemplate:   addToSetTemplatePython,
		addToArrayTemplate: addToArrayTemplatePython,
		addToMapTemplate:   addToMapTemplatePython,
	}}

	javaObj = &java{baseLang: &baseLang{
		name:               execution.Java,
		addToSetTemplate:   addToSetTemplateJava,
		addToArrayTemplate: addToArrayTemplateJava,
		addToMapTemplate:   addToMapTemplateJava,
	}}

	cppObj = &cpp{baseLang: &baseLang{
		name:               execution.Cpp,
		addToSetTemplate:   addToSetTemplateCpp,
		addToArrayTemplate: addToArrayTemplateCpp,
		addToMapTemplate:   addToMapTemplateCpp,
	}}
)

func Generate(args []Argument, lang string) string {
	template := getTemplate(lang, solutionTemplate)
	var argsInit strings.Builder
	var argsPass strings.Builder
	for i, arg := range args {
		argName := fmt.Sprintf("arg_%v", i)
		var langObj Language
		switch lang {
		case execution.Python:
			langObj = pythonObj
		case execution.Java:
			langObj = javaObj
		case execution.Cpp:
			langObj = cppObj
		}
		argInit := arg.Generate(argName, langObj)

		fmt.Fprintf(&argsInit, "%v\n", argInit)
		if i != len(args)-1 {
			fmt.Fprintf(&argsPass, "%v,", argName)
		} else {
			fmt.Fprintf(&argsPass, "%v", argName)
		}
	}
	return fmt.Sprintf(template, argsInit.String(), argsPass.String())
}
