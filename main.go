package main

import (
	"github.com/Strovala/crackview/cmd"
)

func main() {
	cmd.Execute()
}

// func main() {
// 	n := 5
// 	arr := []int{1, 3, 5}
// 	aMap := map[int]int{1: 2, 3: 4}
// 	bMap := map[string]bool{"foo": false, "bar": true}
// 	set := []float64{3.4, 5.6}
// 	inputN := generator.NewSimple(n)
// 	inputArr := generator.NewArray(arr)
// 	inputMapA := generator.NewMap(aMap)
// 	inputMapB := generator.NewMap(bMap)
// 	inputSet := generator.NewSet(set)
// 	args := []generator.Argument{
// 		inputN, inputArr, inputMapA, inputMapB, inputSet,
// 	}
// 	result := generator.Generate(args, execution.Cpp)
// 	fmt.Println(result)
// }
