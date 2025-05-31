package wiktionary

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jamesalexatkin/mot-du-jour-api/model"
)

// GetRandomFrenchWordPage fetches a random French word page from Wiktionary.
func GetRandomFrenchWordPage() (*goquery.Document, error) {
	// URL to make the HTTP request to
	url := "https://en.wiktionary.org/wiki/Special:RandomInCategory/French_lemmas#French"

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("couldn't complete request to wiktionary", "err", err)

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

// ParseWordFromPage parses a Wiktionary HTML page into a word struct,
// extracting meanings and definitions.
func ParseWordFromPage(doc *goquery.Document) model.Word {
	word := model.Word{}

	heading := doc.Find("h1").First()
	word.Name = heading.Text()

	var meanings []model.Entry
	var current model.Entry

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
			current = model.Entry{
				Type: model.EntryType(strings.TrimSpace(h3Section.Text())),
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
				// fmt.Printf("found another language '%s', stopping\n", h2Section.Text())

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
				current.Gender = model.MasculineGender
			case "f":
				current.Gender = model.FeminineGender
			case "":
				current.Gender = model.NoGender
			default:
				current.Gender = model.MasculineOrFeminineGender
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

	word.Entries = meanings

	word.CleanEntries()

	return word
}
