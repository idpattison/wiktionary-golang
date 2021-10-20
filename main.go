package main

import (
	"log"
	"os"
)

func main() {
	// fetch the word or list of words
	// if the word json file exists locally, use that
	// otherwise fetch the wikitext from wiktionary
	// process the word and write the outputs

	// get the required word - if no language given then we assume English
	var word string
	langCode := "en"
	if len(os.Args) < 2 {
		word = "red" // if no args (arg 0 is the program name) then assume word = red for testing
	} else {
		word = os.Args[1]
	}
	if len(os.Args) > 2 {
		langCode = os.Args[2] // capture the language if specified
	}

	// var text string
	// text = "{{IPA|en|/\u0279\u025bd/|[\u027b\u02b7\u025b\u02d1d\u0325]}}"
	// println(text)

	_, err := processWord(word, langCode)
	if err != nil {
		log.Fatalln(err)
	}

}

func processWord(word string, langCode string) (LanguageWord, error) {
	nilWord := new(LanguageWord)
	// get the JSON content for the requested word
	wordData, err := getWordDataFromWiktionary(word)
	if err != nil {
		return *nilWord, err
	}

	// extract the wikitext from the JSON content
	wikitext, err := getWikitext(wordData, word)
	if err != nil {
		return *nilWord, err
	}

	// process the wikitext into sections
	sections := processWikitext(wikitext)

	// get the relevant sections for the language
	languageSections, err := extractLanguageSections(word, langCode, sections)
	if err != nil {
		return *nilWord, err
	}

	// parse the language sections and build a Language struct
	lw := parseSections(word, langCode, languageSections)

	// write the word data to a JSON file
	errw := writeJson(word, langCode, &lw)
	if errw != nil {
		return lw, errw
	}

	// for debug purposes, write the wikitext to a file
	fileName := langCode + "-" + word + "-wiki.txt"
	os.WriteFile(fileName, []byte(wikitext), 0666)

	return lw, nil
}
