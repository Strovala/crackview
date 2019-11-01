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
}

type typeResolver interface {
	Resolve()
}

type base struct {
	Type  string
	Value interface{}
}

func Generate(args []Argument, lang string) string {
	template := getTemplate(lang, solutionTemplate)
	var argsInit strings.Builder
	var argsPass strings.Builder
	for i, arg := range args {
		argName := fmt.Sprintf("arg_%v", i)
		var argInit string
		switch lang {
		case execution.Python:
			argInit = arg.GeneratePython(argName)
		case execution.Java:
			argInit = arg.GenerateJava(argName)
		case execution.Cpp:
			argInit = arg.GenerateCpp(argName)
		}
		fmt.Fprintf(&argsInit, "%v\n", argInit)
		if i != len(args)-1 {
			fmt.Fprintf(&argsPass, "%v,", argName)
		} else {
			fmt.Fprintf(&argsPass, "%v", argName)
		}
	}
	return fmt.Sprintf(template, argsInit.String(), argsPass.String())
}
