package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/tidwall/buntdb"

	"github.com/patrickvalle/heatmap/cmd/apid/config"
	"github.com/patrickvalle/heatmap/cmd/apid/handlers"
	"github.com/patrickvalle/heatmap/internal/ipv6"
)

const csvPath = "GeoLite2-City-Blocks-IPv6.csv"

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func main() {

	// Grab an instance of our env config.
	c := config.New()

	// Create an in-memory instance of buntdb for us to use.
	db, err := buntdb.Open(":memory:")
	if err != nil {
		log.Printf("startup : Failed to load BuntDB into memory: %s", err.Error())
	}

	// Process the CSV file and load up our dataset.
	// TODO: Offload this so it doesn't impact cold startup time.
	file, err := os.Open(csvPath)
	if err != nil {
		log.Printf("startup : Failed to load CSV data: %s", err.Error())
	}
	ipv6.LoadData(db, file)

	// Start the server.
	server := http.Server{
		Addr:           c.APIHost,
		Handler:        handlers.API(c, db),
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

	// Boilerplate shutdown logic below.

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
