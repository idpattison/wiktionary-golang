package wiktionary

import (
	"errors"
	"os"
)

type WiktionaryOptions struct {
	RequiredSections  int16
	RequiredLanguages []string
}

const (
	Sec_Etymology_Text         int16 = 0x0001
	Sec_Etymology_Words        int16 = 0x0002
	Sec_IPA                    int16 = 0x0004
	Sec_Extended_Pronunciation int16 = 0x0008
	Sec_Parts                  int16 = 0x0010 // NB this is required if part attrs, part extended, meanings, translations, synonyms, antonyms are required
	Sec_Part_Attributes        int16 = 0x0020
	Sec_Part_Extended          int16 = 0x0040 // declensions for nouns, conjugations for verbs (if provided)
	Sec_Meanings               int16 = 0x0080
	Sec_Translations           int16 = 0x0100
	Sec_Synonyms               int16 = 0x0200
	Sec_Antonyms               int16 = 0x0400
	Sec_Anagrams               int16 = 0x0800
	Sec_Alternatives           int16 = 0x1000
)
const Sec_Core = Sec_Etymology_Text | Sec_Etymology_Words | Sec_IPA | Sec_Parts | Sec_Meanings
const Sec_All = 0x0FFF

func GetWord(word string, langCode string) (LanguageWord, error) {
	var options WiktionaryOptions
	options.RequiredSections = Sec_All
	options.RequiredLanguages = AllLanguages
	lw, err := processWord(word, langCode, options)
	return lw, err
}

func GetWordWithOptions(word string, langCode string, options WiktionaryOptions) (LanguageWord, error) {
	lw, err := processWord(word, langCode, options)
	return lw, err
}

func GetMeaning(word string, langCode string) (string, error) {
	var options WiktionaryOptions
	options.RequiredSections = Sec_Parts | Sec_Meanings
	options.RequiredLanguages = AllLanguages
	lw, err := processWord(word, langCode, options)
	return lw.Meaning, err
}

func GetTranslations(word string, langCode string, requiredLanguages []string) ([]TranslatedWord, error) {
	var options WiktionaryOptions
	options.RequiredSections = Sec_Parts | Sec_Translations
	options.RequiredLanguages = requiredLanguages
	lw, err := processWord(word, langCode, options)
	if err != nil {
		return nil, err
	}
	var tr []TranslatedWord
	// iterate across all etymologies and parts to find the translations
	for _, etym := range lw.Etymologies {
		for _, part := range etym.Parts {
			tr = append(tr, part.Translations...)
		}
	}
	return tr, err
}

func GetIpaPronunciation(word string, langCode string) (string, error) {
	var options WiktionaryOptions
	options.RequiredSections = Sec_IPA
	options.RequiredLanguages = AllLanguages
	lw, err := processWord(word, langCode, options)
	if err != nil {
		return "", err
	}
	if lw.Ipa == "" {
		return "", errors.New("no IPA found")
	}
	return lw.Ipa, nil
}

func GetEtymologyTree(word string, langCode string, languages []string) TreeNode {
	// get a full etymology tree containing ancestor words
	etym, err := getEtymologyTree(word, langCode, languages)
	if err != nil {
		return etym
	}
	return etym // TODO
}

func GetLanguageFromCode(code string) string {
	// convert a language code to the full name, e.g. for "en" return "English"
	return getLanguageFromCode(code)
}

func processWord(word string, langCode string, options WiktionaryOptions) (LanguageWord, error) {
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
	lw := parseSections(word, langCode, languageSections, options)

	// for debug purposes, write the word data to a JSON file and a wikitext file
	// TODO - remove these once we are done
	errw := writeJson(word, langCode, &lw)
	if errw != nil {
		return lw, errw
	}
	fileName := langCode + "-" + word + ".wikitext"
	os.WriteFile(fileName, []byte(wikitext), 0666)

	return lw, nil
}
