package tracks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"tracks/repository"

	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
)

func updateTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var t repository.Track
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if id == t.Id {
			if n := repository.Update(t); n > 0 {
				w.WriteHeader(204) /* No Content */
			} else if n := repository.Insert(t); n > 0 {
				w.WriteHeader(201) /* Created */
			} else {
				w.WriteHeader(500) /* Internal Server Error */
			}
		} else {
			w.WriteHeader(400) /* Bad Request */
		}
	} else {
		w.WriteHeader(400) /* Bad Request */
	}
}

func listTrack(w http.ResponseWriter, r *http.Request) {
	if t, n := repository.List(); n > 0 {
		w.WriteHeader(200) /* OK */
		json.NewEncoder(w).Encode(t)
	} else {
		w.WriteHeader(500) /* Internal Server Error */
	}
}

func readTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if t, n := repository.Read(id); n > 0 {
		d := repository.Track{Id: t.Id, Audio: t.Audio}
		w.WriteHeader(200) /* OK */
		json.NewEncoder(w).Encode(d)
	} else if n == 0 {
		w.WriteHeader(404) /* Not Found */
	} else {
		w.WriteHeader(500) /* Internal Server Error */
	}
}

func deleteTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if num := repository.Delete(id); num > 0 {
		w.WriteHeader(204) /* No Content */
	} else if num == 0 {
		w.WriteHeader(404) /* No Track */
	} else {
		w.WriteHeader(500) /* Internal Server Error */
	}
}

func Router() http.Handler {
	r := mux.NewRouter()
	fmt.Println("Microservice running...\n")
	/* Create */
	r.HandleFunc("/tracks/{id}", updateTrack).Methods("PUT")
	/* List */
	r.HandleFunc("/tracks", listTrack).Methods("GET")
	/* Read */
	r.HandleFunc("/tracks/{id}", readTrack).Methods("GET")
	/* Delete */
	r.HandleFunc("/tracks/{id}", deleteTrack).Methods("DELETE")
	return r
}
