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
		logrus.Fatalf("Coul not load environment variables (%s)\n", err.Error())
	}

	// Connect to db using a struct made from config info
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
		logrus.Fatalf("Could not connect to Postgres (%s)\n", err.Error())
	}

	// Connect to stan using a struct made from config info
	sc, sub, err := repository.NewStanConn(repository.StanConn{
		ClientId:  viper.GetString("stan.clientid"),
		ClusterId: viper.GetString("stan.clusterid"),
		Subject:   viper.GetString("stan.subject"),
		DB:        db,
	})
	if err != nil {
		logrus.Fatalf("Could not connect to Stan (%s)\n", err.Error())
	}

	// Dependancy injection
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

	// Wait for all of the goroutines to finish
	// By waiting for SIGTERM and SIGINT signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// Shutdown server, db conection, stan subscribtion and stan connection
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("Could not shutdown the server (%s)\n", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Fatalf("Could not close database connection (%s)\n", err.Error())
	}
	if err := sub.Unsubscribe(); err != nil {
		logrus.Fatalf("Could not unsubscribe from stan subject (%s)\n", err.Error())
	}
	if err := sc.Close(); err != nil {
		logrus.Fatalf("Could not close stan connection (%s)\n", err.Error())
	}
}

// Initialises application configuration
func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
