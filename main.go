package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AlirGG/mongoapi/router"
	"github.com/gorilla/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// create a new MongoDB client and connect to the database
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://user_test:password_test@alirggtestscluster.hlxprlo.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// create a new instance of the router and define the server settings
	router := router.Router()

	port := ":8000"
	server := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, router),
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// start the server in a goroutine
	go func() {
		log.Printf("Server started on port %s\n", port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	//shut down the server when the application is terminated
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server shut down gracefully")
}
