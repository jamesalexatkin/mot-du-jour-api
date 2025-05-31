package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jamesalexatkin/mot-du-jour-api/model"
	"github.com/jamesalexatkin/mot-du-jour-api/wiktionary"
)

var motDuJour model.Word
var lastFetched time.Time

func main() {
	http.HandleFunc("/mot_du_jour", serveMotDuJour)
	http.HandleFunc("/mot_spontane", serveMotSpontane)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server listening on http://localhost:" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Error starting server:", err)
	}
}

func serveMotDuJour(w http.ResponseWriter, r *http.Request) {
	// Fetch new word if it's been more than a day since the last one
	if time.Since(lastFetched) > 24*time.Hour {
		doc, err := wiktionary.GetRandomFrenchWordPage()
		if err != nil {
			// TODO: introduce better error handling
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		motDuJour = wiktionary.ParseWordFromPage(doc)
		lastFetched = time.Now()
	}

	marshalledWord, err := json.Marshal(&motDuJour)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Write(marshalledWord)
}

func serveMotSpontane(w http.ResponseWriter, r *http.Request) {
	doc, err := wiktionary.GetRandomFrenchWordPage()
	if err != nil {
		// TODO: introduce better error handling
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	word := wiktionary.ParseWordFromPage(doc)

	marshalledWord, err := json.Marshal(&word)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Write(marshalledWord)
}
