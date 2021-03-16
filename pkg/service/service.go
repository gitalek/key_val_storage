package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	storage2 "key_val_storage/pkg/storage"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	ErrNilStorageProvided serviceError = "service: provided nil storage"
)

type serviceError string

func (s serviceError) Error() string {
	return string(s)
}

type Result struct {
	Data  interface{} `json:"data"`
	Error string      `json:"err"`
}

type App struct {
	storage *storage2.Storage
}

// NewApp is an App constructor.
func NewApp(s *storage2.Storage) (*App, error) {
	if s == nil {
		return nil, ErrNilStorageProvided

	}
	return &App{storage: s}, nil
}

func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	result := &Result{}
	enc := json.NewEncoder(w)
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		result.Error = http.StatusText(http.StatusBadRequest)
		enc.Encode(result)
		return
	}
	val, err := a.storage.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		result.Error = err.Error()
		enc.Encode(result)
		return
	}
	result.Data = val
	enc.Encode(result)
}

func (a *App) List(w http.ResponseWriter, r *http.Request) {
	result := &Result{}
	enc := json.NewEncoder(w)
	val := a.storage.List()
	result.Data = val
	enc.Encode(result)
}

func (a *App) Delete(w http.ResponseWriter, r *http.Request) {
	result := &Result{}
	enc := json.NewEncoder(w)
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		result.Error = http.StatusText(http.StatusBadRequest)
		enc.Encode(result)
		return
	}
	val, err := a.storage.Delete(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		result.Error = err.Error()
		enc.Encode(result)
		return
	}
	result.Data = val
	enc.Encode(result)

}

func (a *App) Upsert(w http.ResponseWriter, r *http.Request) {
	result := &Result{}
	enc := json.NewEncoder(w)
	items, err := parseQuery(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result.Error = http.StatusText(http.StatusBadRequest)
		enc.Encode(result)
		return
	}
	a.storage.Upsert(items)
	result.Data = items
	enc.Encode(result)
}

func (a *App) Backup(pathToFile string, interval time.Duration) chan struct{} {
	ticker := time.NewTicker(interval * time.Millisecond)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				saveToFile(pathToFile, a.storage.Backup())
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()
	return quit
}

func parseQuery(r *http.Request) (map[string]string, error) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, err
	}
	items := make(map[string]string, len(params))
	var curVal string
	for key, val := range params {
		curVal = val[0]
		if len(val) > 1 {
			// choose last val for provided key
			curVal = val[len(val)-1]
		}
		items[key] = curVal
	}
	return items, nil
}

func saveToFile(pathToFile string, state map[string]string) {
	err := os.Remove(pathToFile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := os.Create(pathToFile)
	defer f.Close()
	if err != nil {
		log.Println(err)
		return
	}
	err = json.NewEncoder(f).Encode(state)
	if err != nil {
		log.Println(err)
		return
	}
}
