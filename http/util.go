package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/debug"
)

// Errors when parsing JSON body
var (
	ErrParseJSONBody     = errors.New("Unable to parse JSON body")
	ErrParseJSONBodyType = errors.New("Cannot unmarshal JSON body into type")
)

func errorHandler(next func(http.ResponseWriter, *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			JSONErrorResponse(w, err.Error(), 500)
		}
	}
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
		if closeErr := r.Body.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	err = json.Unmarshal(body, dest)
	if err, ok := err.(*json.SyntaxError); ok && err != nil {
		return ErrParseJSONBody
	}
	if err, ok := err.(*json.UnmarshalTypeError); ok && err != nil {
		return ErrParseJSONBodyType
	}
	return err
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
