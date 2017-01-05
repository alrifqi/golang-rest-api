package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/gorilla/mux"
)

type Note struct {
	Title		string		`json:"title"`
	Description	string		`json:"description"`
	CreatedOn	time.Time 	`json:"createdon"`
}

var noteStore = make(map[string]Note)

var id int = 0

func PostNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note Note

	err:=json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		panic(err)
	}
	note.CreatedOn = time.Now()
	id++
	k:=strconv.Itoa(id)
	noteStore[k]=note

	j, err:=json.Marshal(note)
	if err != nil{
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	var notes []Note

	for _, v:= range noteStore {
		notes = append(notes, v)
	}

	w.Header().Set("Content-Type","application/json")
	j, err := json.Marshal(notes)
	if err != nil{
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func PutNoteHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	k := vars["id"]
	var noteToUpd Note

	err = json.NewDecoder(r.Body).Decode(&noteToUpd)
	if err != nil {
		panic(err)
	}
	if note, ok := noteStore[k]; ok {
		noteToUpd.CreatedOn = note.CreatedOn
		delete(noteStore, k)
		noteStore[k] = noteToUpd
	}else{
		log.Printf("Could not find key of Note %s to delete", k)
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]

	if _, ok := noteStore[k]; ok {
		delete(noteStore, k)
	}else{
		log.Printf("Could not find key of Note %s to delete", k)
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetMain(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Main Site"))
}

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", GetMain).Methods("GET")
	r.HandleFunc("/api/notes", GetNoteHandler).Methods("GET")
	r.HandleFunc("/api/notes", PostNoteHandler).Methods("POST")
	r.HandleFunc("/api/notes/{id}", PutNoteHandler).Methods("PUT")
	r.HandleFunc("/api/notes/{id}", DeleteNoteHandler).Methods("DELETE")

	server := &http.Server{
		Addr: ":8080",
		Handler: r,	
	}
	log.Println("Listening....")
	server.ListenAndServe()
}