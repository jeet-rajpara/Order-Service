package main

import (
	"log"
	"net/http"
	"order_service/database"

	// "os"
	"order_service/controller"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	db, err := database.Connect()
	if err != nil {
		log.Printf("Error in establishing database connection: %v", err)
	}

	r := chi.NewRouter()
	r.Use(database.Middleware(db))
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Post("/createOrder", controller.CreateOrder)
	var port = ":3500"
	srv := &http.Server{
		Addr:        port,
		Handler:     r,
		IdleTimeout: 2 * time.Minute,
	}
	log.Fatal(srv.ListenAndServe())

}
