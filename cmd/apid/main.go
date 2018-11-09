package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/patrickvalle/heatmap/cmd/apid/config"
	"github.com/patrickvalle/heatmap/cmd/apid/handlers"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func main() {

	// ============================================================
	// Configuration

	// TODO: Handle error if there is a need for required configuration.
	c := config.New()

	// Start the server.
	server := http.Server{
		Addr:           c.APIHost,
		Handler:        handlers.API(c),
		ReadTimeout:    c.ReadTimeout,
		WriteTimeout:   c.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// Starting the service, listening for requests.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		log.Printf("startup : Listening @ %s", c.APIHost)
		log.Printf("shutdown : Listener closed : %v", server.ListenAndServe())
		wg.Done()
	}()

	// ============================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	<-osSignals

	// Create context for Shutdown call.
	ctx, cancel := context.WithTimeout(context.Background(), c.ShutdownTimeout)
	defer cancel()

	// Asking listener to shutdown and load shed.
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown : Graceful shutdown did not complete in %v : %v", c.ShutdownTimeout, err)

		if err := server.Close(); err != nil {
			log.Printf("shutdown : Error killing server : %v", err)
		}
	}

	// Waiting for service to complete that load shedding.
	wg.Wait()

	log.Println("main : Completed")
}
