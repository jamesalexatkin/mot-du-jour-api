package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jamesalexatkin/mot-du-jour-api/wiktionary"
)

//go:embed page.html
var page string

func main() {
	http.HandleFunc("/", serveHello)
	http.HandleFunc("/mot_du_jour", serveMotDuJour)

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

func serveHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, page)
}

func serveMotDuJour(w http.ResponseWriter, r *http.Request) {
	doc, err := wiktionary.GetRandomFrenchWordPage()
	if err != nil {
		// TODO: introduce better error handling
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	word := wiktionary.ParseWordFromPage(doc)

	// fmt.Printf("%+v\n", word)

	marshalledWord, err := json.Marshal(&word)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Write(marshalledWord)
}
