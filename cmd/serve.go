package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	chiCors "github.com/go-chi/cors"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(serveCmdNew)
}

// JSONRecoverer is a middleware that recovers from panics, and returns a HTTP 500 (Internal Server Error).
func JSONRecoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
				debug.PrintStack()
				JSONResponse(w, "Unexpected error occurred", http.StatusInternalServerError)
				return
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func serializeToString(data interface{}) (s string) {
	var b []byte
	var err error
	b, err = json.MarshalIndent(data, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func writeJSONResponse(w http.ResponseWriter, data interface{}, httpCode int, key string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	transformedData := map[string]interface{}{
		key: data,
	}
	w.WriteHeader(httpCode)
	fmt.Fprint(w, serializeToString(transformedData)+"\n")
}

// JSONResponse writes json response
func JSONResponse(w http.ResponseWriter, data interface{}, httpCode int) {
	writeJSONResponse(w, data, httpCode, "data")
}

// JSONErrorResponse writes json error response
func JSONErrorResponse(w http.ResponseWriter, data interface{}, httpCode int) {
	writeJSONResponse(w, data, httpCode, "error")
}

// Unmarshal deserializes request body
func Unmarshal(dest interface{}, r *http.Request) (err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer func() {
		// do not assign to err directly since this is executing last and
		// we might set it to nil when someone after is setting it to non-nil
		closeErr := r.Body.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()

	err = json.Unmarshal(body, dest)
	if err, ok := err.(*json.SyntaxError); ok && err != nil {
		return errors.New("Unable to parse JSON body")
	}
	if err, ok := err.(*json.UnmarshalTypeError); ok && err != nil {
		return errors.New("Cannot unmarshal JSON body into type")
	}
	return err
}

func homeHandler(w http.ResponseWriter, r *http.Request) error {
	JSONResponse(w, fmt.Sprintf("Of languages we support: %v", viper.GetStringSlice("languages")), http.StatusOK)
	return nil
}

// CodeRequest is DTO for request for code run
type CodeRequest struct {
	Text string `json:"text"`
	Lang string `json:"lang"`
}

// CodeResponse is DTO for response for code run
type CodeResponse struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

func execPython(text string) (*CodeResponse, error) {
	err := ioutil.WriteFile("main.py", []byte(text), 0644)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("python3", "main.py")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	_ = cmd.Run()
	err = os.Remove("main.py")
	if err != nil {
		return nil, err
	}
	resp := &CodeResponse{
		Output: out.String(),
		Error:  errOut.String(),
	}
	return resp, nil
}

func execCpp(text string) (*CodeResponse, error) {
	err := ioutil.WriteFile("main.cpp", []byte(text), 0644)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("c++", "-o", "main", "main.cpp")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	_ = cmd.Run()
	if errOut.Len() != 0 {
		err = os.Remove("main.cpp")
		if err != nil {
			return nil, err
		}
		return &CodeResponse{
			Error: errOut.String(),
		}, nil
	}
	cmd = exec.Command("./main")
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	_ = cmd.Run()
	err = os.Remove("main.cpp")
	err = os.Remove("main")

	return &CodeResponse{
		Output: out.String(),
	}, nil
}

func execJava(text string) (*CodeResponse, error) {
	err := ioutil.WriteFile("main.java", []byte(text), 0644)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("javac", "main.java")
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	_ = cmd.Run()
	if errOut.Len() != 0 {
		err = os.Remove("main.java")
		if err != nil {
			return nil, err
		}
		return &CodeResponse{
			Error: errOut.String(),
		}, nil
	}
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return nil, err
	}

	unparsedExecName := ""
	for _, file := range files {
		if strings.Contains(file.Name(), ".class") {
			unparsedExecName = file.Name()
			break
		}
	}
	execName := strings.Split(unparsedExecName, ".")[0]
	cmd = exec.Command("java", execName)
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	_ = cmd.Run()
	err = os.Remove("main.java")
	err = os.Remove(unparsedExecName)

	return &CodeResponse{
		Output: out.String(),
	}, nil
}

const (
	PYTHON = "python"
	JAVA   = "java"
	CPP    = "c++"
)

func homeHandlerPost(w http.ResponseWriter, r *http.Request) error {
	var data CodeRequest
	err := Unmarshal(&data, r)
	if err != nil {
		return err
	}

	var execFunc func(string) (*CodeResponse, error)
	switch data.Lang {
	case PYTHON:
		execFunc = execPython
	case JAVA:
		execFunc = execJava
	case CPP:
		execFunc = execCpp
	}
	resp, err := execFunc(data.Text)
	if err != nil {
		return err
	}

	JSONResponse(w, *resp, http.StatusOK)
	return nil
}

func errorHandler(next func(http.ResponseWriter, *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		if err != nil {
			JSONErrorResponse(w, err.Error(), 500)
		}
	}
}

func makeServer() *http.Server {
	cors := chiCors.New(chiCors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	})

	router := chi.NewMux()
	router.Use(cors.Handler)
	router.Use(JSONRecoverer)

	router.Get("/", errorHandler(homeHandler))
	router.Post("/", errorHandler(homeHandlerPost))

	s := &http.Server{
		Addr:         fmt.Sprintf(":%v", viper.GetString("port")),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return s
}

var serveCmdNew = &cobra.Command{
	Use:   "server",
	Short: "Start HTTP Server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// init db etc

		srv := makeServer()

		done := make(chan os.Signal, 1)
		signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()
		log.Print("Server Started")

		<-done
		log.Print("Server Stopped")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			// extra handling here
			// close db, etc
			cancel()
		}()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server Shutdown Failed:%+v", err)
		}
		return nil
	},
}
