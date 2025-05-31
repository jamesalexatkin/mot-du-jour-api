package model

import "fmt"

// Word represents the full structured encapsulation of a dictionary word.
type Word struct {
	Name    string
	Entries []Entry
}

// CleanEntries removes any other Wiktionary elements parsed into the word's entries.
// Sections like Etymology and Derived Terms are present on the same level with h2 headings
// so can easily end up here.
//
// Useful reference: https://en.wiktionary.org/wiki/Category:French_lemmas
func (w *Word) CleanEntries() {
	var newEntries []Entry

	for _, entries := range w.Entries {
		switch entries.Type {
		case AdjectiveEntryType, AdverbEntryType, ArticleEntryType,
			ConjunctionEntryType, ContractionEntryType, DeterminerEntryType,
			InterfixEntryType, InterjectionEntryType, MorphemeEntryType,
			MultiwordTermEntryType, LetterEntryType, NounEntryType,
			NumeralEntryType, ParticipleEntryType, PhraseEntryType,
			PrefixEntryType, PrepositionEntryType, PrepositionalPhraseEntryType,
			PostpositionEntryType, ProverbEntryType, ProperNounEntryType,
			SuffixEntryType, VerbEntryType:

			newEntries = append(newEntries, entries)
		default:
			fmt.Printf("Unrecognised entry type '%s'\n", entries.Type)
		}
	}

	w.Entries = newEntries
}

// Entry represents a distinct form of a given word.
// The same spelling of a word may be used in different forms (e.g. noun or verb), as
// well as having different genders (e.g. le tour vs la tour).
type Entry struct {
	Type        EntryType
	Gender      Gender
	Definitions []string
}

// EntryType represents the part of speech or lexical category of a word entry.
type EntryType string

const (
	AdjectiveEntryType           EntryType = "Adjective"
	AdverbEntryType              EntryType = "Adverb"
	ArticleEntryType             EntryType = "Article"
	ConjunctionEntryType         EntryType = "Conjunction"
	ContractionEntryType         EntryType = "Contraction"
	DeterminerEntryType          EntryType = "Determiner"
	InterfixEntryType            EntryType = "Interfix"
	InterjectionEntryType        EntryType = "Interjection"
	MorphemeEntryType            EntryType = "Morpheme"
	MultiwordTermEntryType       EntryType = "Multiword term"
	LetterEntryType              EntryType = "Letter"
	NounEntryType                EntryType = "Noun"
	NumeralEntryType             EntryType = "Numeral"
	ParticipleEntryType          EntryType = "Participle"
	PhraseEntryType              EntryType = "Phrase"
	PrefixEntryType              EntryType = "Prefix"
	PrepositionEntryType         EntryType = "Preposition"
	PrepositionalPhraseEntryType EntryType = "Prepositional phrase"
	PostpositionEntryType        EntryType = "Postposition"
	ProverbEntryType             EntryType = "Proverb"
	ProperNounEntryType          EntryType = "Proper noun"
	SuffixEntryType              EntryType = "Suffix"
	VerbEntryType                EntryType = "Verb"
)

// Gender represents grammatical gender.
type Gender string

const (
	// MasculineGender represents the masculine grammatical gender.
	MasculineGender Gender = "masc."
	// FeminineGender represents the feminine grammatical gender.
	FeminineGender Gender = "fem."
	// NoGender represents genderless words (e.g. verbs, adjectives).
	NoGender Gender = ""
	// MasculineOrFeminineGender represents grammatical gender for words that can be both masculine and feminine.
	MasculineOrFeminineGender Gender = "masc. ou fem."
)
