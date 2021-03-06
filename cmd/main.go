package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gitalek/key_val_storage/pkg/service"
	"github.com/gitalek/key_val_storage/pkg/storage"
	"github.com/gorilla/mux"
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

	get := service.HeaderMW(app.Get)
	list := service.HeaderMW(app.List)
	del := service.HeaderMW(app.Delete)
	upsert := service.HeaderMW(app.Upsert)

	r := mux.NewRouter()
	r.HandleFunc("/get/{key}", get).Methods("GET")
	r.HandleFunc("/list", list).Methods("GET")
	r.HandleFunc("/delete/{key}", del).Methods("DELETE", "POST")
	r.HandleFunc("/upsert", upsert).Methods("PUT", "POST")

	quit := app.Backup(*buPath, time.Duration(*buInterval))

	defer func() {
		quit <- struct{}{}
	}()

	log.Printf("starting server on %s\n", *addr)
	log.Println(http.ListenAndServe(*addr, r))
}
