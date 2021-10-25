package wiktionary

import "os"

func ProcessWord(word string, langCode string) (LanguageWord, error) {
	nilWord := new(LanguageWord)
	// get the JSON content for the requested word
	wordData, err := getWordDataFromWiktionary(word, langCode)
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
	fileName := langCode + "-" + word + ".wikitext"
	os.WriteFile(fileName, []byte(wikitext), 0666)

	return lw, nil
}
