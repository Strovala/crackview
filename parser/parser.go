package parser

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Strovala/crackview/generator"
)

const inputFile = "input.txt"

const (
	arrayType = "[]"
	mapType   = "{}"
	setType   = "()"
)

const (
	intType    = "int"
	boolType   = "bool"
	stringType = "string"
	floatType  = "float"
)

func Parse(input string) ([]generator.Argument, error) {
	splitedInput := strings.TrimSpace(input)
	var result []generator.Argument
	lines := strings.Split(splitedInput, "\n")
	for _, line := range lines {
		splited := strings.Split(line, "=>")
		if len(splited) < 2 {
			continue
		}
		val := splited[0]
		valType := strings.TrimSpace(splited[1])
		argObj, err := parse(val, valType)
		if err != nil {
			return nil, err
		}
		result = append(result, argObj)
	}
	return result, nil
}

func parse(val string, valType string) (generator.Argument, error) {
	starting := valType[:2]
	splited := []string{valType}
	var argObj generator.Argument
	var err error
	switch starting {
	case arrayType:
		splited = strings.Split(valType, arrayType)
		vType := splited[1]
		argObj, err = parseArray(val, vType)
		if err != nil {
			return nil, err
		}
	case mapType:
		splited = strings.Split(valType, mapType)
		vType := splited[1]
		vTypeSplited := strings.Split(vType, ",")
		keyType := vTypeSplited[0]
		valueType := vTypeSplited[1]
		argObj, err = parseMap(val, keyType, valueType)
		if err != nil {
			return nil, err
		}
	case setType:
		splited = strings.Split(valType, setType)
		vType := splited[1]
		argObj, err = parseSet(val, vType)
		if err != nil {
			return nil, err
		}
	default:
		argObj, err = parseSimple(val, valType)
		if err != nil {
			return nil, err
		}
	}
	return argObj, nil
}

func parseSimple(val string, baseType string) (*generator.Simple, error) {
	var res *generator.Simple
	switch baseType {
	case intType:
		result, err := parseInt(val)
		if err != nil {
			return nil, err
		}
		res = generator.NewSimple(result)
	case floatType:
		result, err := parseFloat(val)
		if err != nil {
			return nil, err
		}
		res = generator.NewSimple(result)
	case boolType:
		result, err := parseBool(val)
		if err != nil {
			return nil, err
		}
		res = generator.NewSimple(result)
	case stringType:
		result, err := parseString(val)
		if err != nil {
			return nil, err
		}
		res = generator.NewSimple(result)
	default:
		return nil, errors.New("unknown type")
	}
	return res, nil
}

func parseMapIntInt(val string) (map[int]int, error) {
	splited := toStringArray(val)
	result := make(map[int]int)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseInt(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseInt(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMapIntFloat(val string) (map[int]float64, error) {
	splited := toStringArray(val)
	result := make(map[int]float64)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseInt(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseFloat(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMapIntBool(val string) (map[int]bool, error) {
	splited := toStringArray(val)
	result := make(map[int]bool)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseInt(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseBool(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMapIntString(val string) (map[int]string, error) {
	splited := toStringArray(val)
	result := make(map[int]string)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseInt(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseString(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMapFloatInt(val string) (map[float64]int, error) {
	splited := toStringArray(val)
	result := make(map[float64]int)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseFloat(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseInt(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMapFloatFloat(val string) (map[float64]float64, error) {
	splited := toStringArray(val)
	result := make(map[float64]float64)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseFloat(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseFloat(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMapFloatBool(val string) (map[float64]bool, error) {
	splited := toStringArray(val)
	result := make(map[float64]bool)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseFloat(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseBool(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMapFloatString(val string) (map[float64]string, error) {
	splited := toStringArray(val)
	result := make(map[float64]string)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseFloat(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseString(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMapBoolInt(val string) (map[bool]int, error) {
	splited := toStringArray(val)
	result := make(map[bool]int)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseBool(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseInt(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMapBoolFloat(val string) (map[bool]float64, error) {
	splited := toStringArray(val)
	result := make(map[bool]float64)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseBool(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseFloat(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMapBoolBool(val string) (map[bool]bool, error) {
	splited := toStringArray(val)
	result := make(map[bool]bool)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseBool(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseBool(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMapBoolString(val string) (map[bool]string, error) {
	splited := toStringArray(val)
	result := make(map[bool]string)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseBool(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseString(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMapStringInt(val string) (map[string]int, error) {
	splited := toStringArray(val)
	result := make(map[string]int)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseString(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseInt(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMapStringFloat(val string) (map[string]float64, error) {
	splited := toStringArray(val)
	result := make(map[string]float64)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseString(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseFloat(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMapStringBool(val string) (map[string]bool, error) {
	splited := toStringArray(val)
	result := make(map[string]bool)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseString(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseBool(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMapStringString(val string) (map[string]string, error) {
	splited := toStringArray(val)
	result := make(map[string]string)
	for _, value := range splited {
		keyValueList := strings.Split(value, ":")
		parsedKey, err := parseString(keyValueList[0])
		if err != nil {
			return nil, err
		}
		parsedValue, err := parseString(keyValueList[1])
		if err != nil {
			return nil, err
		}
		result[parsedKey] = parsedValue
	}
	return result, nil
}

func parseMap(val, baseKeyType, baseValueType string) (*generator.InputMap, error) {
	var res *generator.InputMap
	switch baseKeyType {
	case intType:
		switch baseValueType {
		case intType:
			result, err := parseMapIntInt(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		case floatType:
			result, err := parseMapIntFloat(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		case boolType:
			result, err := parseMapIntBool(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		case stringType:
			result, err := parseMapIntString(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		default:
			return nil, errors.New("unknown type")
		}
	case floatType:
		switch baseValueType {
		case intType:
			result, err := parseMapFloatInt(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		case floatType:
			result, err := parseMapFloatFloat(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		case boolType:
			result, err := parseMapFloatBool(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		case stringType:
			result, err := parseMapFloatString(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		default:
			return nil, errors.New("unknown type")
		}
	case boolType:
		switch baseValueType {
		case intType:
			result, err := parseMapBoolInt(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		case floatType:
			result, err := parseMapBoolFloat(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		case boolType:
			result, err := parseMapBoolBool(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		case stringType:
			result, err := parseMapBoolString(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		default:
			return nil, errors.New("unknown type")
		}
	case stringType:
		switch baseValueType {
		case intType:
			result, err := parseMapStringInt(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		case floatType:
			result, err := parseMapStringFloat(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		case boolType:
			result, err := parseMapStringBool(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		case stringType:
			result, err := parseMapStringString(val)
			if err != nil {
				return nil, err
			}
			res = generator.NewMap(result)
		default:
			return nil, errors.New("unknown type")
		}
	default:
		return nil, errors.New("unknown type")
	}
	return res, nil
}

func parseInt(value string) (int, error) {
	value = strings.TrimSpace(value)
	parsedValue, err := strconv.ParseInt(value, 10, 64)
	return int(parsedValue), err
}

func parseFloat(value string) (float64, error) {
	value = strings.TrimSpace(value)
	parsedValue, err := strconv.ParseFloat(value, 64)
	return parsedValue, err
}

func parseBool(value string) (bool, error) {
	value = strings.TrimSpace(value)
	parsedValue, err := strconv.ParseBool(value)
	return parsedValue, err
}

func parseString(value string) (string, error) {
	value = strings.TrimSpace(value)
	return value, nil
}

func toStringArray(val string) []string {
	val = val[1 : len(val)-2]
	return strings.Split(val, ",")
}

func parseIntArray(val string) ([]int, error) {
	splited := toStringArray(val)
	var result []int
	for _, value := range splited {
		parsedValue, err := parseInt(value)
		if err != nil {
			return nil, err
		}
		result = append(result, parsedValue)
	}
	return result, nil
}

func parseFloatArray(val string) ([]float64, error) {
	splited := toStringArray(val)
	var result []float64
	for _, value := range splited {
		parsedValue, err := parseFloat(value)
		if err != nil {
			return nil, err
		}
		result = append(result, parsedValue)
	}
	return result, nil
}

func parseBoolArray(val string) ([]bool, error) {
	splited := toStringArray(val)
	var result []bool
	for _, value := range splited {
		parsedValue, err := parseBool(value)
		if err != nil {
			return nil, err
		}
		result = append(result, parsedValue)
	}
	return result, nil
}

func parseStringArray(val string) ([]string, error) {
	splited := toStringArray(val)
	return splited, nil
}

func parseArray(val, baseType string) (*generator.Array, error) {
	var res *generator.Array
	switch baseType {
	case intType:
		result, err := parseIntArray(val)
		if err != nil {
			return nil, err
		}
		res = generator.NewArray(result)
	case floatType:
		result, err := parseFloatArray(val)
		if err != nil {
			return nil, err
		}
		res = generator.NewArray(result)
	case boolType:
		result, err := parseBoolArray(val)
		if err != nil {
			return nil, err
		}
		res = generator.NewArray(result)
	case stringType:
		result, err := parseStringArray(val)
		if err != nil {
			return nil, err
		}
		res = generator.NewArray(result)
	default:
		return nil, errors.New("unknown type")
	}
	return res, nil
}

func parseSet(val, baseType string) (*generator.Set, error) {
	var res *generator.Set
	switch baseType {
	case intType:
		result, err := parseIntArray(val)
		if err != nil {
			return nil, err
		}
		res = generator.NewSet(result)
	case floatType:
		result, err := parseFloatArray(val)
		if err != nil {
			return nil, err
		}
		res = generator.NewSet(result)
	case boolType:
		result, err := parseBoolArray(val)
		if err != nil {
			return nil, err
		}
		res = generator.NewSet(result)
	case stringType:
		result, err := parseStringArray(val)
		if err != nil {
			return nil, err
		}
		res = generator.NewSet(result)
	default:
		return nil, errors.New("unknown type")
	}
	return res, nil
}

// 5 => int
// [1, 3, 5] => []int
// {1: 2, 3: 4} => {}int,int
// {"foo": false, "bar": true} => {}string,bool
// (3.4, 5.6) => ()float
