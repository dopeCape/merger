package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	server "github.com/dopeCape/schduler/internal"
	"github.com/dopeCape/schduler/pkg/apikey"
	"github.com/dopeCape/schduler/pkg/broker"
	rdb "github.com/dopeCape/schduler/pkg/db"
	"github.com/dopeCape/schduler/pkg/inspector"
	"github.com/dopeCape/schduler/pkg/scheduler"
	"github.com/dopeCape/schduler/pkg/shared"
	"github.com/dopeCape/schduler/pkg/suscriber"
	"github.com/joho/godotenv"
)

const redisAddr string = "127.0.0.1:6379"

func main() {
	// load env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// connect to turso
	err = rdb.NewDB()
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	// core
	// quque config
	config := shared.Config{Concurrency: 1000, RedisAddress: redisAddr}
	// app api
	// quque worker
	apiKeyService := apikey.NewApiKeySerice()
	broker, client := broker.RunBroker(config)
	inspectorClient, inspector := inspector.GetInspector(config)
	schedulerClient := scheduler.GetSchduler()
	// Run the subscriber in a goroutine

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		schedulerClient.Start()
	}()

	// Run the subscriber in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		suscriber.Run(config)
	}()
	httpServer := server.GetHandler(broker, inspector, apiKeyService)
	server := &http.Server{
		Addr:    ":8000",
		Handler: httpServer,
	}

	// Run the HTTP server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("HTTP server failed: %v", err)
		}
	}()

	// Shutdown the HTTP server with a timeout

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	wg.Wait()
	defer client.Close()
	defer inspectorClient.Close()
	defer schedulerClient.Shutdown()
}
