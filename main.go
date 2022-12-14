package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manureddy7143/GolangStarter/utils/database"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var configPath *string

func main() {
	//Setting config flag
	configPath = flag.String("config-path", "conf", "conf")
	flag.Parse()
	loadconfig()
	setUpdb()
	r := gin.New()
	setupRoutes(r)
	startServer(r)

}

//Load configuration file
func loadconfig() {
	viper.AddConfigPath(*configPath)
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		if readErr, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Panic().Msgf("No config file found at %s\n", *configPath)
		} else {
			log.Panic().Msgf("Error reading config file: %s\n", readErr)
		}
	}
}

//connect to mysql database
func setUpdb() {
	database.GetInstance()
}

// startServer - Start server
func startServer(r *gin.Engine) {
	srv := &http.Server{
		Addr:    viper.GetString("server.port"),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("Listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 5)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...\n")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("Server forced to shutdown: %s\n", err)
	}

	log.Info().Msg("Server exiting\n")
}
