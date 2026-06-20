package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hafizul16103123/student-api/internal/config"
	student "github.com/hafizul16103123/student-api/internal/http/handlers"
	"github.com/hafizul16103123/student-api/internal/storage/sqlite"
)
func main(){
	// Load config
	cfg:=config.MustLoad()
	slog.Info("Configuration loaded")
	// Database setup
	db,err:=sqlite.NewSQLite(cfg)
	if err!=nil{
		log.Fatal("Failed to connect to database",err.Error())
	}
	slog.Info("Connected to database successfully",slog.String("env",cfg.Env))

	// setup router
	router:=http.NewServeMux()
	router.HandleFunc("POST /api/students",student.New(db))
	router.HandleFunc("GET /api/students/{id}",student.GetById(db))

	// setup server
	server:=http.Server{
		Addr:cfg.Addr,
		Handler:router,
	}


	done:=make(chan os.Signal,1)// Create a channel to receive OS signals (SIGINT/SIGTERM)

	// Register signals that we want to listen for to done channel
	signal.Notify(done,os.Interrupt,syscall.SIGINT,syscall.SIGTERM)

	// Start server in a separate goroutine
	// Because ListenAndServe() blocks until server stops
	go func(){
		// Start accepting HTTP requests
		err:=server.ListenAndServe()
		if err!= nil{
			log.Fatal("Failed to start server")
		}
	}()

	// Main goroutine waits here until shutdown signal is received
	// It blocks until Ctrl+C, SIGTERM, etc.
	<-done

	// Create context with timeout
	// Server has maximum 5 seconds to finish existing requests
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second*5)
	// Release context resources after shutdown completes
	defer cancel()


	// Gracefully stop the server
	// It:
	// 1. Stops accepting new requests
	// 2. Waits for active requests to finish
	// 3. Closes connections
	// 4. Need to provide context timeout because sometimes Shutdown() can hang forever. 
	//    when context timeout reached then Shutdown cancelled ❌ and error returned.
	err=server.Shutdown(ctx)

	if err!=nil{
		slog.Error("Failed to Shutdown server",slog.String("error",err.Error()))
	}
	slog.Info("Server shutdown successfully")

}

/*
Graceful Shutdown:

Receive SIGTERM
        |
        ↓
Stop accepting new requests
        |
        ↓
Finish existing requests
        |
        ↓
Close DB connections
        |
        ↓
Exit

This pattern is essential for production Go microservices.
*/
