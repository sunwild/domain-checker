package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/sunwild/api/internal/datebase"
	"github.com/sunwild/api/internal/restapi"
	"github.com/sunwild/api/internal/service/domains"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	db, _ := datebase.NewRepository(ctx)
	dService := domains.NewService(db)
	dHandler := restapi.NewDomainHandler(dService)

	// routes
	r := mux.NewRouter() // use gorilla/mux for routes
	r.HandleFunc("/domains", dHandler.GetAllDomains).Methods("GET")
	r.HandleFunc("/domains/{id}", dHandler.GetDomainById).Methods("GET")
	r.HandleFunc("/domains", dHandler.AddDomain).Methods("POST")
	r.HandleFunc("/domains/{id}", dHandler.UpdateDomainById).Methods("PUT")
	r.HandleFunc("/domains/{id}", dHandler.DeleteDomainById).Methods("DELETE")
	r.HandleFunc("/check-domains", dHandler.HandleManualDomainCheck).Methods("POST")

	// settings http service
	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:8000",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	fmt.Println("Starting server on http://localhost:8000")
	log.Fatal(srv.ListenAndServe())
}
