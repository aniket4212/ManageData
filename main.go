package main

import (
	"context"
	"log"
	"managedata/config"
	"managedata/db/mysql"
	"managedata/db/redis"
	"managedata/routes"
	"managedata/utils"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	config.GetConfigurations()

	// connect mysql database
	mysql.ConnectMysqlDB()

	redis.ConnectRedis()
	utils.SetLogLevel()

}

func main() {
	// read config file
	config.GetConfigurations()

	err := run()
	if err != nil {
		log.Fatalln(err)
	}

	// close database connections
	defer func() {
		log.Println("Closing Mysql database connection")
		mysql.CloseMysql()
	}()
}

func run() (err error) {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// server setup
	gin.SetMode(gin.ReleaseMode)
	// Create a new Gin router instance.
	router := gin.New()
	// Register handlers.
	routes.RegisterRouter(router)

	// Start HTTP server.
	srv := &http.Server{
		Addr:         ":" + config.AppConfig.Server.Port,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	log.Printf("Server is up and running on http://localhost:%v%v/ping\n", config.AppConfig.Server.Port, config.AppConfig.Prefix)

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return err
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()

		// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
		err = srv.Shutdown(context.Background())
	}

	return nil
}
