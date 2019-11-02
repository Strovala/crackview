package generator

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Strovala/crackview/execution"
)

const (
	templatesPath     = "generator/templates"
	arrayTemplate     = "array.txt"
	setTemplate       = "set.txt"
	simpleTemplate    = "simple.txt"
	mapTemplate       = "map.txt"
	inputArgsTemplate = "input_args.txt"
	mainTemplate      = "main.txt"
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

func javaWrapperClass(valueType string) string {
	javaType, ok := javaLookUpWrapperClass[valueType]
	if !ok {
		javaType = capitalizeFirstLetter(valueType)
	}
	return javaType
}

func getType(lookUp map[string]string, valueType string) string {
	langType, ok := lookUp[valueType]
	if !ok {
		langType = valueType
	}
	return langType
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

// Argument is interface for generating code snippet for specific argument with value and type
type Argument interface {
	ResolveType()
	Generate(name string, lang Language) string
}

type argument struct {
	argType   string
	keyType   string
	valueType string
	value     interface{}
}

// Language is interface for generating code snippet for specific language
type Language interface {
	// Returns type of given type in code ex. Type("string") -> String in Java
	Type(valType string) string
	// Returns value of given type in code ex. GenerateValue("foo", "string") -> "\"foo\""
	GenerateValue(val interface{}, valType string) string
	GenerateSimpleTemplate(name, val interface{}, valType string) string
	GenerateSetTemplate(argType string, name, value string) string
	GenerateAddToSetTemplate(name string, val interface{}, valType string) string
	GenerateArrayTemplate(argType string, name, value string) string
	GenerateAddToArrayTemplate(name string, val interface{}, valType string) string
	GenerateMapTemplate(argKeyType, argValueType string, name, value string) string
	GenerateAddToMapTemplate(name string, keyValue interface{}, keyType string, val interface{}, valType string) string
	Generate(inputArgs string, solution string)
}

type language struct {
	name               string
	addToSetTemplate   string
	addToArrayTemplate string
	addToMapTemplate   string
	mainName           string
}

func (l *language) getTemplate(template string) string {
	return getTemplate(l.name, template)
}

type python struct {
	*language
}

func (l *python) Type(valType string) string {
	return ""
}

func (l *python) GenerateValue(val interface{}, valType string) string {
	return valuePython(val, valType)
}

func (l *python) GenerateSimpleTemplate(name, val interface{}, valType string) string {
	template := l.getTemplate(simpleTemplate)
	return fmt.Sprintf(template, name, l.GenerateValue(val, valType))
}

func (l *python) GenerateArrayTemplate(argType string, name, value string) string {
	template := l.getTemplate(arrayTemplate)
	return fmt.Sprintf(template, name, value)
}

func (l *python) GenerateAddToArrayTemplate(name string, val interface{}, valType string) string {
	return fmt.Sprintf(l.addToArrayTemplate, l.GenerateValue(val, valType))
}

func (l *python) GenerateSetTemplate(argType string, name, value string) string {
	template := l.getTemplate(setTemplate)
	return fmt.Sprintf(template, name, value)
}

func (l *python) GenerateAddToSetTemplate(name string, val interface{}, valType string) string {
	return fmt.Sprintf(l.addToSetTemplate, l.GenerateValue(val, valType))
}

func (l *python) GenerateMapTemplate(argKeyType, argValueType string, name, value string) string {
	template := l.getTemplate(mapTemplate)
	return fmt.Sprintf(template, name, value)
}

func (l *python) GenerateAddToMapTemplate(name string, keyValue interface{}, keyType string, val interface{}, valType string) string {
	return fmt.Sprintf(l.addToMapTemplate, l.GenerateValue(keyValue, keyType), l.GenerateValue(val, valType))
}

func (l *python) Generate(inputArgs string, solution string) {
	template := l.getTemplate(mainTemplate)
	code := fmt.Sprintf(template, solution, inputArgs)
	_ = ioutil.WriteFile(l.mainName, []byte(code), 0644)
}

type java struct {
	*language
}

func (l *java) wrapperClass(valType string) string {
	return javaWrapperClass(l.Type(valType))
}

func (l *java) Type(valType string) string {
	return getType(javaLookUpType, valType)
}

func (l *java) GenerateValue(val interface{}, valType string) string {
	return value(val, valType)
}

func (l *java) GenerateSimpleTemplate(name, val interface{}, valType string) string {
	template := l.getTemplate(simpleTemplate)
	langType := l.Type(valType)
	return fmt.Sprintf(template, langType, name, l.GenerateValue(val, valType))
}

func (l *java) GenerateArrayTemplate(argType string, name, value string) string {
	template := l.getTemplate(arrayTemplate)
	langType := l.Type(argType)
	return fmt.Sprintf(template, langType, name, langType, value)
}

func (l *java) GenerateAddToArrayTemplate(name string, val interface{}, valType string) string {
	return fmt.Sprintf(l.addToArrayTemplate, l.GenerateValue(val, valType))
}

func (l *java) GenerateSetTemplate(argType string, name, value string) string {
	template := l.getTemplate(setTemplate)
	langType := l.wrapperClass(argType)
	return fmt.Sprintf(template, langType, name, langType, value)
}

func (l *java) GenerateAddToSetTemplate(name string, val interface{}, valType string) string {
	return fmt.Sprintf(l.addToSetTemplate, l.GenerateValue(val, valType))
}

func (l *java) GenerateMapTemplate(argKeyType, argValueType string, name, value string) string {
	template := l.getTemplate(mapTemplate)
	langKeyType := l.wrapperClass(argKeyType)
	langValueType := l.wrapperClass(argValueType)
	return fmt.Sprintf(template, langKeyType, langValueType, name, langKeyType, langValueType, value)
}

func (l *java) GenerateAddToMapTemplate(name string, keyValue interface{}, keyType string, val interface{}, valType string) string {
	return fmt.Sprintf(l.addToMapTemplate, l.GenerateValue(keyValue, keyType), l.GenerateValue(val, valType))
}

func (l *java) Generate(inputArgs string, solution string) {
	template := l.getTemplate(mainTemplate)
	code := fmt.Sprintf(template, inputArgs, solution)
	_ = ioutil.WriteFile(l.mainName, []byte(code), 0644)
}

type cpp struct {
	*language
}

func (l *cpp) Type(valType string) string {
	return getType(cppLookUpType, valType)
}

func (l *cpp) GenerateValue(val interface{}, valType string) string {
	return value(val, valType)
}

func (l *cpp) GenerateSimpleTemplate(name, val interface{}, valType string) string {
	template := l.getTemplate(simpleTemplate)
	langType := l.Type(valType)
	return fmt.Sprintf(template, langType, name, l.GenerateValue(val, valType))
}

func (l *cpp) GenerateArrayTemplate(argType string, name, value string) string {
	template := l.getTemplate(arrayTemplate)
	langType := l.Type(argType)
	return fmt.Sprintf(template, langType, name, value)
}

func (l *cpp) GenerateAddToArrayTemplate(name string, val interface{}, valType string) string {
	return fmt.Sprintf(l.addToArrayTemplate, name, l.GenerateValue(val, valType))
}

func (l *cpp) GenerateSetTemplate(argType string, name, value string) string {
	template := l.getTemplate(setTemplate)
	langType := l.Type(argType)
	return fmt.Sprintf(template, langType, name, value)
}

func (l *cpp) GenerateAddToSetTemplate(name string, val interface{}, valType string) string {
	return fmt.Sprintf(l.addToSetTemplate, name, l.GenerateValue(val, valType))
}

func (l *cpp) GenerateMapTemplate(argKeyType, argValueType string, name, value string) string {
	template := l.getTemplate(mapTemplate)
	langKeyType := l.Type(argKeyType)
	langValueType := l.Type(argValueType)
	return fmt.Sprintf(template, langKeyType, langValueType, name, value)
}

func (l *cpp) GenerateAddToMapTemplate(name string, keyValue interface{}, keyType string, val interface{}, valType string) string {
	return fmt.Sprintf(l.addToMapTemplate, name, l.GenerateValue(keyValue, keyType), l.GenerateValue(val, valType))
}

func (l *cpp) Generate(inputArgs string, solution string) {
	template := l.getTemplate(mainTemplate)
	code := fmt.Sprintf(template, solution, inputArgs)
	_ = ioutil.WriteFile(l.mainName, []byte(code), 0644)
}

var (
	pythonObj = &python{language: &language{
		name:               execution.Python,
		addToSetTemplate:   addToSetTemplatePython,
		addToArrayTemplate: addToArrayTemplatePython,
		addToMapTemplate:   addToMapTemplatePython,
		mainName:           execution.PythonMainName,
	}}

	javaObj = &java{language: &language{
		name:               execution.Java,
		addToSetTemplate:   addToSetTemplateJava,
		addToArrayTemplate: addToArrayTemplateJava,
		addToMapTemplate:   addToMapTemplateJava,
		mainName:           execution.JavaMainName,
	}}

	cppObj = &cpp{language: &language{
		name:               execution.Cpp,
		addToSetTemplate:   addToSetTemplateCpp,
		addToArrayTemplate: addToArrayTemplateCpp,
		addToMapTemplate:   addToMapTemplateCpp,
		mainName:           execution.CppMainName,
	}}
)

// Generate generates code snippet for executing solution with provided arguments
func Generate(args []Argument, lang string, solution string) {
	template := getTemplate(lang, inputArgsTemplate)
	var argsInit strings.Builder
	var argsPass strings.Builder
	var langObj Language
	switch lang {
	case execution.Python:
		langObj = pythonObj
	case execution.Java:
		langObj = javaObj
	case execution.Cpp:
		langObj = cppObj
	}
	for i, arg := range args {
		argName := fmt.Sprintf("input_%v", i)
		argInit := arg.Generate(argName, langObj)
		fmt.Fprintf(&argsInit, "%v\n", argInit)
		if i != len(args)-1 {
			fmt.Fprintf(&argsPass, "%v,", argName)
		} else {
			fmt.Fprintf(&argsPass, "%v", argName)
		}
	}
	inputArgs := fmt.Sprintf(template, argsInit.String(), argsPass.String())
	langObj.Generate(inputArgs, solution)
}
