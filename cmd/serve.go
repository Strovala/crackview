package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	chiCors "github.com/go-chi/cors"
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
				log.Println("Panic occurred!", rvr)
				JSONResponse(w, "Unexpected error occurred", 500)
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

// JSONResponse writes json response
func JSONResponse(w http.ResponseWriter, data interface{}, httpCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	transformedData := map[string]interface{}{
		"data": data,
	}
	w.WriteHeader(httpCode)
	fmt.Fprint(w, serializeToString(transformedData)+"\n")
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

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		JSONResponse(w, "blant", 200)
	})

	s := &http.Server{
		Addr:        fmt.Sprintf(":%v", viper.GetString("port")),
		Handler:     router,
		ReadTimeout: 2 * time.Minute,
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
