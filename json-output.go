package main

import (
	"encoding/json"
	"os"
)

type LanguageWord struct {
	Word             string      `json:"word"`
	Meaning          string      `json:"meaning,omitempty"`
	LanguageName     string      `json:"lang"`
	LanguageCode     string      `json:"lang-code"`
	AlternativeForms string      `json:"alts,omitempty"`
	Pronunciations   []string    `json:"pron,omitempty"`
	Ipa              string      `json:"ipa,omitempty"`
	Etymologies      []Etymology `json:"etym,omitempty"`
}

type Etymology struct {
	Name           string         `json:"name"`
	Text           string         `json:"text,omitempty"`
	Words          []LinkedWord   `json:"words,omitempty"`
	Parts          []PartOfSpeech `json:"parts,omitempty"`
	Pronunciations []string       `json:"pron,omitempty"`
	Ipa            string         `json:"ipa,omitempty"`
}

// relationships for LinkedWord
const (
	Root       string = "root"
	Inherited  string = "inherited"
	Cognate    string = "cognate"
	Borrowed   string = "borrowed"
	Descendant string = "descendant"
)

type LinkedWord struct {
	SourceLanguage  string `json:"-"`
	SourceWord      string `json:"-"`
	Relationship    string `json:"type"`
	TargetLanguage  string `json:"lang"`
	TargetWord      string `json:"word"`
	Transliteration string `json:"translit,omitempty"`
	Meaning         string `json:"meaning,omitempty"`
}

type PartOfSpeech struct {
	Name       string            `json:"name"`
	Headword   string            `json:"head,omitempty"`
	Attributes map[string]string `json:"attrs,omitempty"`
	Meanings   []string          `json:"meanings,omitempty"`
}

func writeJson(word string, langCode string, lw *LanguageWord) error {
	b, err := json.Marshal(lw)
	if err != nil {
		return err
	}
	// write file
	fileName := langCode + "-" + word + ".json"
	errf := os.WriteFile(fileName, b, 0666)
	if errf != nil {
		return errf
	}

	return nil
}
