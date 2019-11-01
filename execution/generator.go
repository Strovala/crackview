package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

const templatesPath = "execution/templates"

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

type generator interface {
	GenerateSnippetPython(name string) string
	GenerateSnippetJava(name string) string
	GenerateSnippetCpp(name string) string
}

type typeResolver interface {
	Resolve()
}

type base struct {
	Type  string
	Value interface{}
}

func (p *base) GenerateSnippetCpp(name string) string {
	return "Not Implemented!"
}

type simple struct {
	*base
}

func newSimple(value interface{}) *simple {
	result := &simple{
		base: &base{
			Value: value,
		},
	}
	result.Resolve()
	return result
}

func (p *simple) Resolve() {
	reflectType := getReflectType(p.Value)
	inputType := reflectType.String()
	p.Type = inputType
}

func (p *simple) GenerateSnippetPython(name string) string {
	template := getTemplate("python", "simple.txt")
	return fmt.Sprintf(template, name, valuePython(p.Value, p.Type))
}

func (p *simple) GenerateSnippetJava(name string) string {
	template := getTemplate("java", "simple.txt")
	javaType := getJavaType(p.Type)
	return fmt.Sprintf(template, javaType, name, value(p.Value, p.Type))
}

func (p *simple) GenerateSnippetCpp(name string) string {
	template := getTemplate("cpp", "simple.txt")
	cppType := getCppType(p.Type)
	return fmt.Sprintf(template, cppType, name, value(p.Value, p.Type))
}

type array struct {
	*base
}

func newArray(value interface{}) *array {
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

func (p *array) GenerateSnippetPython(name string) string {
	template := getTemplate("python", "array.txt")
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for i := 0; i < val.Len(); i++ {
		fmt.Fprintf(&builder, "%v,", valuePython(val.Index(i), p.Type))
	}
	return fmt.Sprintf(template, name, builder.String())
}

func (p *array) GenerateSnippetJava(name string) string {
	template := getTemplate("java", "array.txt")
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for i := 0; i < val.Len(); i++ {
		fmt.Fprintf(&builder, "%v,", value(val.Index(i), p.Type))
	}
	javaType := getJavaType(p.Type)
	return fmt.Sprintf(template, javaType, name, javaType, builder.String())
}

func (p *array) GenerateSnippetCpp(name string) string {
	template := getTemplate("cpp", "array.txt")
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for i := 0; i < val.Len(); i++ {
		fmt.Fprintf(&builder, "%v.push_back(%v);", name, value(val.Index(i), p.Type))
	}
	cppType := getCppType(p.Type)
	return fmt.Sprintf(template, cppType, name, builder.String())
}

type inputMap struct {
	*base
	KeyType   string
	ValueType string
}

func newMap(value interface{}) *inputMap {
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

func (p *inputMap) GenerateSnippetPython(name string) string {
	template := getTemplate("python", "map.txt")
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for _, k := range val.MapKeys() {
		v := val.MapIndex(k)
		fmt.Fprintf(&builder, "%v:%v,", valuePython(k, p.KeyType), valuePython(v, p.ValueType))
	}
	return fmt.Sprintf(template, name, builder.String())
}

func (p *inputMap) GenerateSnippetJava(name string) string {
	template := getTemplate("java", "map.txt")
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for _, k := range val.MapKeys() {
		v := val.MapIndex(k)
		fmt.Fprintf(&builder, "put(%v,%v);", valuePython(k, p.KeyType), value(v, p.ValueType))
	}
	javaKeyType := javaWrapperClass(getJavaType(p.KeyType))
	javaValueType := javaWrapperClass(getJavaType(p.ValueType))
	return fmt.Sprintf(template, javaKeyType, javaValueType, name, javaKeyType, javaValueType, builder.String())
}

func (p *inputMap) GenerateSnippetCpp(name string) string {
	template := getTemplate("cpp", "map.txt")
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for _, k := range val.MapKeys() {
		v := val.MapIndex(k)
		fmt.Fprintf(&builder, "%v[%v]=%v;", name, valuePython(k, p.KeyType), value(v, p.ValueType))
	}
	cppKeyType := getCppType(p.KeyType)
	cppValueType := getCppType(p.ValueType)
	return fmt.Sprintf(template, cppKeyType, cppValueType, name, builder.String())
}

type set struct {
	*base
}

func newSet(value interface{}) *set {
	result := &set{
		base: &base{
			Value: value,
		},
	}
	result.Resolve()
	return result
}

func (p *set) Resolve() {
	reflectType := getReflectType(p.Value)
	inputType := reflectType.String()
	p.Type = inputType[2:]
}

func (p *set) GenerateSnippetPython(name string) string {
	template := getTemplate("python", "set.txt")
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for i := 0; i < val.Len(); i++ {
		fmt.Fprintf(&builder, "%v,", valuePython(val.Index(i), p.Type))
	}
	return fmt.Sprintf(template, name, builder.String())
}

func (p *set) GenerateSnippetJava(name string) string {
	template := getTemplate("java", "set.txt")
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for i := 0; i < val.Len(); i++ {
		fmt.Fprintf(&builder, "add(%v);", value(val.Index(i), p.Type))
	}
	javaType := javaWrapperClass(getJavaType(p.Type))
	return fmt.Sprintf(template, javaType, name, javaType, builder.String())
}

func (p *set) GenerateSnippetCpp(name string) string {
	template := getTemplate("cpp", "set.txt")
	var builder strings.Builder
	val := reflect.ValueOf(p.Value)
	for i := 0; i < val.Len(); i++ {
		fmt.Fprintf(&builder, "%v.insert(%v);", name, value(val.Index(i), p.Type))
	}
	cppType := getCppType(p.Type)
	return fmt.Sprintf(template, cppType, name, builder.String())
}

func generatePython(args []generator) string {
	template := getTemplate("python", "solution.txt")
	var argsInit strings.Builder
	var argsPass strings.Builder
	for i, arg := range args {
		argName := fmt.Sprintf("arg_%v", i)
		fmt.Fprintf(&argsInit, "%v\n", arg.GenerateSnippetPython(argName))
		if i != len(args)-1 {
			fmt.Fprintf(&argsPass, "%v,", argName)
		} else {
			fmt.Fprintf(&argsPass, "%v", argName)
		}
	}
	return fmt.Sprintf(template, argsInit.String(), argsPass.String())
}

func generateJava(args []generator) string {
	template := getTemplate("java", "solution.txt")
	var argsInit strings.Builder
	var argsPass strings.Builder
	for i, arg := range args {
		argName := fmt.Sprintf("arg_%v", i)
		fmt.Fprintf(&argsInit, "%v\n", arg.GenerateSnippetJava(argName))
		if i != len(args)-1 {
			fmt.Fprintf(&argsPass, "%v,", argName)
		} else {
			fmt.Fprintf(&argsPass, "%v", argName)
		}
	}
	return fmt.Sprintf(template, argsInit.String(), argsPass.String())
}

func generateCpp(args []generator) string {
	template := getTemplate("cpp", "solution.txt")
	var argsInit strings.Builder
	var argsPass strings.Builder
	for i, arg := range args {
		argName := fmt.Sprintf("arg_%v", i)
		fmt.Fprintf(&argsInit, "%v\n", arg.GenerateSnippetCpp(argName))
		if i != len(args)-1 {
			fmt.Fprintf(&argsPass, "%v,", argName)
		} else {
			fmt.Fprintf(&argsPass, "%v", argName)
		}
	}
	return fmt.Sprintf(template, argsInit.String(), argsPass.String())
}

func main() {
	n := 5
	arr := []int{1, 3, 5}
	aMap := map[int]int{1: 2, 3: 4}
	bMap := map[string]bool{"foo": false, "bar": true}
	set := []float64{3.4, 5.6}
	inputN := newSimple(n)
	inputArr := newArray(arr)
	inputMapA := newMap(aMap)
	inputMapB := newMap(bMap)
	inputSet := newSet(set)
	result := generateCpp([]generator{
		inputN, inputArr, inputMapA, inputMapB, inputSet,
	})
	fmt.Println(result)
}

// 5 => int
// [1, 3, 5] => array<int>; []int
// {1: 2, 3: 4} => map<int,int>; map[int]int
// {"foo": false, "bar": true} => map<string,bool>; map[string]bool
// (3.4, 5.6) => set<float>; set[float]

// n = 5
// nums = [1,3,5]
// a_map = {1: 2, 3: 4}
// b_map = {"foo": False, "bar": True}
// a_set = (3.4, 5.6)
// result = Solution.code(n, nums, a_map, b_map, a_set)

// int n = 5;
// int[] nums = new int[]{1, 3, 5};
// Map<Integer, Integer> aMap = new HashMap<Integer, Integer>(){{
// 	put(1, 2);
// 	put(3, 4);
// }};
// Map<String, Boolean> bMap = new HashMap<String, Boolean>(){{
// 	put("foo", false);
// 	put("bar", true);
// }};
// Set<Double> aSet = new HashSet<Double>(){{
// 	add(3.4);
// 	add(5.6);
// }};

// int n = 5;
// vector<int> nums;
// nums.push_back(1);
// nums.push_back(3);
// nums.push_back(5);

// map<int, int> a_map;
// a_map[1] = 2;
// a_map[3] = 4;

// map<string, bool> b_map;
// b_map["foo"] = false;
// b_map["bar"] = true;

// set<double> a_set;
// a_set.insert(3.4);
// a_set.insert(5.6);
