package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	crackviewHttp "github.com/Strovala/crackview/http"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(serveCmdNew)
}

func makeServer() *http.Server {
	router := crackviewHttp.NewCrackviewHandler()

	return &http.Server{
		Addr:         fmt.Sprintf(":%v", viper.GetString("port")),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
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
