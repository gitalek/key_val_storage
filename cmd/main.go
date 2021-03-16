package main

import (
	"flag"
	"github.com/gorilla/mux"
	"key_val_storage/pkg/service"
	"key_val_storage/pkg/storage"
	"log"
	"net/http"
	"time"
)

func main() {
	addr := flag.String("addr", ":5555", "HTTP network address")
	buPath := flag.String("file", "./db.json", "path to backup file")
	buInterval := flag.Int("buInterval", 1000, "backup interval in milliseconds")
	emptyInitStateAllowed := flag.Bool("allowEmptyDBOnStart", true, "allow empty db if backup file didn't find")
	flag.Parse()

	strg, err := storage.NewStorageFromFile(*buPath, *emptyInitStateAllowed)
	if err != nil {
		log.Fatal(err)
	}
	app, err := service.NewApp(strg)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/get/{key}", app.Get).Methods("GET")
	r.HandleFunc("/list", app.List).Methods("GET")
	r.HandleFunc("/delete/{key}", app.Delete).Methods("DELETE", "POST")
	r.HandleFunc("/upsert", app.Upsert).Methods("PUT", "POST")

	quit := app.Backup(*buPath, time.Duration(*buInterval))
	defer func() {
		quit <- struct{}{}
	}()

	log.Printf("starting server on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, r))
}
