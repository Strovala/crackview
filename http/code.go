package http

import (
	"fmt"
	"net/http"

	"github.com/Strovala/crackview/execution"
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
	var executor execution.Executor
	switch data.Lang {
	case execution.Python:
		executor = execution.NewPythonExecutor()
	case execution.Java:
		executor = execution.NewJavaExecutor()
	case execution.Cpp:
		executor = execution.NewCppExecutor()
	}
	resp, err := executor.Execute(data.Text)
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
