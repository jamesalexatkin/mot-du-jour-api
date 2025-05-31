package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "embed"

	"github.com/PuerkitoBio/goquery"
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

type Word struct {
	Name     string
	Meanings []Meaning
}

type Gender string

const (
	MasculineGender           Gender = "masc."
	FeminineGender            Gender = "fem."
	NoGender                  Gender = ""
	MasculineOrFeminineGender Gender = "masc. ou fem."
)

type Meaning struct {
	Type        string // TODO: enum
	Gender      Gender
	Definitions []string
}

func serveMotDuJour(w http.ResponseWriter, r *http.Request) {
	doc, err := getWiktionaryRandomFrenchWordPage()
	if err != nil {
		// TODO: handle error response

		return
	}

	word := Word{}

	heading := doc.Find("h1").First()
	word.Name = heading.Text()

	var meanings []Meaning
	var current Meaning

	// Find the h2 for "French"
	selection := doc.Find("h2#French").Parent()
	if selection.Length() == 0 {
		log.Fatal("French section not found")
	}

	// Iterate over all HTML node elements after the French heading
	// Wiktionary's HTML structure is flat not hierarchical, so we have to examine
	// all the siblings
	for n := selection.Next(); n.Length() > 0; n = n.Next() {
		if n.HasClass("mw-heading3") {
			h3Section := n.Find("h3").First()
			if h3Section == nil {
				// TODO: defend here, continue?
				continue
			}

			// If we had an in-progress meaning, store it
			if current.Type != "" && len(current.Definitions) > 0 {
				meanings = append(meanings, current)
			}
			current = Meaning{
				Type: strings.TrimSpace(h3Section.Text()),
			}
		}

		// Stop if we hit another h2 (e.g. next language)
		if n.HasClass("mw-heading2") {
			h2Section := n.Find("h2").First()
			if h2Section == nil {
				// TODO: defend here, continue?
				continue
			}
			if h2Section.Text() != "French" {
				fmt.Printf("found another language '%s', stopping\n", h2Section.Text())

				break
			}
		}

		switch goquery.NodeName(n) {
		// Parse grammatical gender out of the paragraph
		case "p":
			genderSpan := n.Find("span.gender").First()
			if genderSpan == nil {
				// TODO: defend here, continue?
				continue
			}

			switch genderSpan.Text() {
			case "m":
				current.Gender = MasculineGender
			case "f":
				current.Gender = FeminineGender
			case "":
				current.Gender = NoGender
			default:
				current.Gender = MasculineOrFeminineGender
			}
		case "ol":
			// The ordered list contains the definitions as list items
			n.Find("li").Each(func(_ int, listItem *goquery.Selection) {
				// Some definitions come with usage examples
				// These are always on a second line, so we split on new lines to cut them out
				justDefinition := strings.Split(listItem.Text(), "\n")[0]

				// TODO: fix bug in here with quotations being pulled through on some words

				current.Definitions = append(current.Definitions, justDefinition)
			})
		}
	}

	// Append last meaning
	if current.Type != "" && len(current.Definitions) > 0 {
		meanings = append(meanings, current)
	}

	word.Meanings = meanings

	word.CleanMeanings()

	fmt.Printf("%+v\n", word)
}

// CleanMeanings removes any other Wiktionary elements parsed into the word's meanings.
// Sections like Etymology and Derived Terms are present on the same level with h2 headings
// so can easily end up here.
func (w *Word) CleanMeanings() {
	var newMeanings []Meaning

	for _, meanings := range w.Meanings {
		if meanings.Type == "Adjective" ||
			meanings.Type == "Adverb" ||
			meanings.Type == "Article" ||
			meanings.Type == "Conjunction" ||
			meanings.Type == "Contraction" ||
			meanings.Type == "Determiner" ||
			meanings.Type == "Interfix" ||
			meanings.Type == "Interjection" ||
			meanings.Type == "Morpheme" ||
			meanings.Type == "Multiword term" ||
			meanings.Type == "Letter" ||
			meanings.Type == "Noun" ||
			meanings.Type == "Numeral" ||
			meanings.Type == "Phrase" ||
			meanings.Type == "Prefix" ||
			meanings.Type == "Preposition" ||
			meanings.Type == "Postposition" ||
			meanings.Type == "Proverb" ||
			meanings.Type == "Proper noun" ||
			meanings.Type == "Suffix" ||
			meanings.Type == "Verb" {
			newMeanings = append(newMeanings, meanings)
		}
	}

	w.Meanings = newMeanings
}

func getWiktionaryRandomFrenchWordPage() (*goquery.Document, error) {
	// URL to make the HTTP request to
	url := "https://en.wiktionary.org/wiki/Special:RandomInCategory/French_lemmas#French"

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		log.Println("couldn't complete request to wiktionary", "err", err)

		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		// TODO: handle error

		return nil, nil
	}
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc, nil
}
