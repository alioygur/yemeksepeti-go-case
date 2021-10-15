package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alioygur/yemekepeti-go-case/herodb"
)

// Composition root
func main() {
	// listen for SIGINT or SIGTERM for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	appCtx, appCtxCancel := context.WithCancel(context.Background())

	// init db
	db := herodb.New()
	{
		// register file persister (creates a new file if it not exist)
		f, err := herodb.NewFilepersist("data.db")
		if err != nil {
			log.Fatal(err)
		}
		if err := db.MakePersistent(appCtx, f, time.Second*60); err != nil {
			log.Fatal(err)
		}
	}

	// init http handler
	handler := &Handler{Storage: db}

	// init router
	mux := http.NewServeMux()

	mux.HandleFunc("/get", handler.get)
	mux.HandleFunc("/set", handler.set)

	// init http server
	srv := &http.Server{Handler: mux, Addr: serverAddr()}

	// start http server in goroutine (graceful shutdown)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	log.Printf("server started %s", srv.Addr)

	// exited!
	oscall := <-done
	log.Printf("system call: %+v", oscall)

	// graceful shutdown the app
	appCtxCancel()

	// gracefull shutdown the http server
	serverCtx, serverCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer serverCtxCancel()
	if err := srv.Shutdown(serverCtx); err != nil {
		log.Fatalf("server shutdown failed: %+v", err)
	}
	log.Print("server exited properly")
}

func serverAddr() string {
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	}
	return ":8080"
}
