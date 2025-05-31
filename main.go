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
	MasculineGender Gender = "masc."
	FeminineGender  Gender = "fem."
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

	// Find the <h2> for "French"
	selection := doc.Find("h2#French").Parent()
	if selection.Length() == 0 {
		log.Fatal("French section not found")
	}
	// fmt.Println(selection.First().Text())

	for n := selection.Next(); n.Length() > 0; n = n.Next() {
		// fmt.Println(goquery.NodeName(n))
		// fmt.Println(n.Text())

		if n.HasClass("mw-heading4") {
			h3Section := n.Find("h4").First()
			if h3Section == nil {
				// TODO: defend here, continue?
				continue
			}

			// switch h3Section.Text() {
			// case "Noun", "Adjective", "Verb", "Adverb":
			// If we had an in-progress meaning, store it
			if current.Type != "" && len(current.Definitions) > 0 {
				meanings = append(meanings, current)
			}
			current = Meaning{
				Type: strings.TrimSpace(n.Find("h4").First().Text()),
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
			}
		case "ol":
			// The ordered list contains the definitions as list items
			n.Find("li").Each(func(_ int, listItem *goquery.Selection) {
				// Some definitions come with usage examples
				// These are always on a second line, so we split on new lines to cut them out
				justDefinition := strings.Split(listItem.Text(), "\n")[0]

				current.Definitions = append(current.Definitions, justDefinition)
			})

		}
	}

	// Append last meaning
	if current.Type != "" && len(current.Definitions) > 0 {
		meanings = append(meanings, current)
	}

	word.Meanings = meanings

	// var frenchSection *goquery.Selection
	// doc.Find("h2").Each(func(i int, s *goquery.Selection) {
	// 	// Find the French section
	// 	if s.Text() == "French" {
	// 		frenchSection = s

	// 		return
	// 	}
	// })

	// foundFrench := false

	// doc.Find("div.mw-heading2").Each(func(i int, s *goquery.Selection) {
	// 	s.Find("h2").Each(func(i int, s *goquery.Selection) {
	// 		// Find the French section
	// 		if s.Text() == "French" {
	// 			fmt.Println("FOUND FRENCH")
	// 			fmt.Println(s.Text())

	// 			foundFrench = true

	// 			return
	// 		}
	// 	})

	// 	if foundFrench {
	// 		s.Find("h3").Each(func(i int, s *goquery.Selection) {
	// 			fmt.Println(s.Text())

	// 			switch s.Text() {
	// 			case "Noun":
	// 				m := Meaning{}

	// 				genderSpan := s.Find("span.gender").First()
	// 				m.Gender = genderSpan.Text()

	// 				s.Find("ol.li").Each(func(i int, s *goquery.Selection) {
	// 					m.Definitions = append(m.Definitions, s.Text())
	// 				})

	// 				word.Meanings = append(word.Meanings, m)
	// 			default:
	// 				break
	// 			}
	// 		})
	// 	}
	// })

	// Traverse the document
	// doc.Find("h2, h3, p, ol, ul, li").Each(func(i int, s *goquery.Selection) {
	// 	tag := goquery.NodeName(s)

	// 	// Find the French section
	// 	if tag == "h2" && strings.Contains(s.Text(), "French") {
	// 		foundFrench = true
	// 		return
	// 	}

	// 	// Stop when another language section is reached
	// 	if foundFrench &&
	// 		(tag == "h2" && !strings.Contains(s.Text(), "French")) {
	// 		foundFrench = false
	// 	}

	// 	m := Meaning{}

	// 	// If in the French section, print relevant content
	// 	if foundFrench && (tag == "p" || tag == "li") {
	// 		text := strings.TrimSpace(s.Text())
	// 		if text != "" {
	// 			fmt.Println(text)
	// 		}
	// 	}
	// 	if foundFrench && tag == "li" {
	// 		m.Definitions = append(m.Definitions, s.Text())
	// 	}

	// })

	word.CleanMeanings()

	fmt.Printf("%+v\n", word)
}

func (w *Word) CleanMeanings() {
	var newMeanings []Meaning

	for _, meanings := range w.Meanings {
		if meanings.Type == "Noun" ||
			meanings.Type == "Verb" ||
			meanings.Type == "Adjective" ||
			meanings.Type == "Adverb" {
			newMeanings = append(newMeanings, meanings)
		}
	}

	w.Meanings = newMeanings
}

func getWiktionaryRandomFrenchWordPage() (*goquery.Document, error) {
	// URL to make the HTTP request to
	// url := "https://en.wiktionary.org/wiki/Special:RandomInCategory/French_lemmas#French"
	// url := "https://en.wiktionary.org/wiki/deuxi%C3%A8me_ligne#French"
	url := "https://en.wiktionary.org/wiki/tour#French"

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

	// Read the response body
	// bytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println("couldn't read response body", "err", err)

	// 	return nil, err
	// }

	// Print the body as a string
	// fmt.Println("HTML:\n\n", string(bytes))

	// doc, err := html.Parse(resp.Body)
	// if err != nil {
	// 	fmt.Println("Error:", err)

	// 	return nil, err
	// }

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc, nil
}
