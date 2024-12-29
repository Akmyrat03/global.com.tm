package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akmyrat/global/pkg/database"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := InitConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := database.ConnectToDB(database.Config{
		Host:     viper.GetString("storage.host"),
		Port:     viper.GetString("storage.port"),
		Username: viper.GetString("storage.username"),
		Password: viper.GetString("storage.password"),
		DBName:   viper.GetString("storage.dbname"),
		SSLMode:  viper.GetString("storage.sslmode"),
	})

	if err != nil {
		log.Fatalf("Failed to initialize db: %v", err.Error())
	}

	app := gin.Default()

	server := &http.Server{
		Addr:    viper.GetString("APP.host"),
		Handler: app,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	log.Println("Server is running...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if err := db.Close(); err != nil {
		log.Fatalf("Failed to close database: %v", err)
	}
	log.Println("Server exited gracefully ")

}

func InitConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.SetDefault("APP.host", "localhost:3232")
	return nil
}
