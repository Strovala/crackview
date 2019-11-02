package http

import (
	"fmt"
	"net/http"

	"github.com/Strovala/crackview/execution"
	"github.com/Strovala/crackview/generator"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

func newCodeHandler() http.Handler {
	c := code{}
	mux := chi.NewMux()
	mux.Get("/info", errorHandler(c.Info))
	mux.Post("/execute", errorHandler(c.Execute))
	return mux
}

type code struct{}

func (c *code) Info(w http.ResponseWriter, r *http.Request) error {
	JSONResponse(w, fmt.Sprintf("Of languages we support: %v", viper.GetStringSlice("languages")), http.StatusOK)
	return nil
}

func (c *code) Execute(w http.ResponseWriter, r *http.Request) error {
	var data CodeRequest
	if err := Unmarshal(&data, r); err != nil {
		return err
	}
	args := initArgs()
	var executor execution.Executor
	switch data.Lang {
	case execution.Python:
		generator.Generate(args, execution.Python, data.Text)
		executor = execution.NewPythonExecutor()
	case execution.Java:
		generator.Generate(args, execution.Java, data.Text)
		executor = execution.NewJavaExecutor()
	case execution.Cpp:
		generator.Generate(args, execution.Cpp, data.Text)
		executor = execution.NewCppExecutor()
	}
	resp, err := executor.Execute()
	if err != nil {
		return err
	}
	JSONResponse(w, *resp, http.StatusOK)
	return nil
}

// CodeRequest is DTO for request for code run
type CodeRequest struct {
	Text string `json:"text"`
	Lang string `json:"lang"`
}

func initArgs() []generator.Argument {
	n := 5
	arr := []int{1, 3, 5}
	aMap := map[int]int{1: 2, 3: 4}
	bMap := map[string]bool{"foo": false, "bar": true}
	set := []float64{3.4, 5.6}
	inputN := generator.NewSimple(n)
	inputArr := generator.NewArray(arr)
	inputMapA := generator.NewMap(aMap)
	inputMapB := generator.NewMap(bMap)
	inputSet := generator.NewSet(set)
	args := []generator.Argument{
		inputN, inputArr, inputMapA, inputMapB, inputSet,
	}
	return args
}
