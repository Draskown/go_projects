package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/Draskown/WBL0/server"
	"github.com/Draskown/WBL0/server/handler"
	"github.com/Draskown/WBL0/server/repository"
	"github.com/Draskown/WBL0/server/service"
)

func main() {
	// Set formatter for logrus
	logrus.SetFormatter(new(logrus.JSONFormatter))
	// Find Viper config
	if err := initConfig(); err != nil {
		logrus.Fatalf("Could not read Viper config (%s)\n", err)
	}
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Coul not load environment variables (%s)", err.Error())
	}

	// Create a struct for database from config info
	db, err := repository.NewPostgresDB(repository.DBConfig{
		Host:     viper.GetString("postgres.host"),
		Port:     viper.GetString("postgres.port"),
		Username: viper.GetString("postgres.username"),
		DBName:   viper.GetString("postgres.dbname"),
		SSLMode:  viper.GetString("postgres.sslmode"),
		// Get password from environment variable
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("Could not connect to Postgres (%s)", err.Error())
	}

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	hnd := handler.NewHandler(service)

	srv := new(server.Server)

	go func() {
		// Launch the server with a port from Viper config
		if err := srv.Run(viper.GetString("port"), hnd.InitRoutes()); err != nil {
			logrus.Fatalf("Could not run the server (%s)\n", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("Could not shutdown the server (%s)", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Fatalf("Could not close database connection (%s)", err.Error())
	}
}

func initConfig() error {
	// Config relative path
	viper.AddConfigPath("config")
	// Config filename
	viper.SetConfigName("viperConfig")

	return viper.ReadInConfig()
}
