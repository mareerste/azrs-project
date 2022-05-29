package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// GRACEFUL SHUTDOWN :

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	cf, err := New()
	if err != nil {
		log.Fatal(err)
	}

	server := Service{
		cf: cf,
	}

	router.HandleFunc("/config/", count(server.createConfigHandler)).Methods("POST")
	router.HandleFunc("/configs/", count(server.getAllConfig)).Methods("GET")
	router.HandleFunc("/config/{id}/", count(server.getConfigHandler)).Methods("GET")
	router.HandleFunc("/config/{id}/{version}", count(server.getConfigHandlerVersion)).Methods("GET")
	router.HandleFunc("/config/{id}/{version}/{labels}", count(server.getFilteredConfigHandler)).Methods("GET")
	router.HandleFunc("/config/{id}/", count(server.delConfigHandler)).Methods("DELETE")
	router.HandleFunc("/config/{id}/{version}", count(server.createNewVersionHandler)).Methods("POST")
	router.HandleFunc("/config/{id}/{version}", count(server.addConfigToExistingGroupHandler)).Methods("PUT")
	router.HandleFunc("/config/{id}/{version}/", count(server.delConfigHandlerVersion)).Methods("DELETE")
	// *Server's scraped metrics UI Path (localhost:9090
	router.Path("/metrics").Handler(metricsHandler())

	// start server
	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
