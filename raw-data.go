package wiktionary

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
)

type Section struct {
	header string
	lines  []string
}

func getWordDataFromWiktionary(word string, langCode string) ([]byte, error) {
	// for a given word, retrieve the word's JSON data from Wiktionary
	urlHead := "https://en.wiktionary.org/w/api.php?action=parse&page="
	urlTail := "&prop=wikitext&format=json"

	// make an HHTP request to Wiktionary
	url := urlHead + url.QueryEscape(getPageTitle(word, langCode)) + urlTail
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// process the response to retrieve the body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// return the byte string
	return body, nil
}

func getWikitext(wordData []byte, word string) (string, error) {
	// use regex to locate the wikitext within the JSON data
	// NB we could process the JSON itself, but as we don't need anything else, this would be inefficient
	re := regexp.MustCompile(`\"wikitext\":\{\"\*\":\"(.*?)\"\}\}\}$`)
	match := re.FindStringSubmatch(string(wordData))
	if len(match) == 0 {
		msg := fmt.Sprintf("No wikitext for word '%s'", word)
		return "", errors.New(msg)
	}

	wikitext, err := convertWikitext(match[1])
	if err != nil {
		msg := fmt.Sprintf("Wikitext decode error for word '%s'", word)
		return "", errors.New(msg)
	}

	return wikitext, nil
}

func convertWikitext(text string) (string, error) {
	// decode any scripts represented by surrogate pairs - this includes Avestan, for example
	re := regexp.MustCompile(`\\ud([8-9a-fA-F][0-9a-fA-F]{2})`) // matches the format \udxxx
	match := re.FindAllStringSubmatch(text, -1)
	// this returns a slice of matches - we need to pair them up
	for i := 0; i < len(match); i += 2 {
		hiByte, _ := strconv.ParseInt(match[i][0][2:], 16, 32)
		loByte, _ := strconv.ParseInt(match[i+1][0][2:], 16, 32)
		r := utf16.DecodeRune(rune(hiByte), rune(loByte))
		textToReplace := match[i][0] + match[i+1][0]
		replaceWith := string(r)
		text = strings.ReplaceAll(text, textToReplace, replaceWith)
	}

	// now decode the more usual backslash-encoded characters
	// NB strconv.Unquote requires and returns a string enclosed in quotes - we need to handle that
	convertedText, err := strconv.Unquote(`"` + text + `"`)
	if err != nil {
		return "", nil
	}
	convertedText = strings.Trim(convertedText, `"`)
	return convertedText, nil

}

func processWikitext(wikitext string) []Section {
	// split the wikitext up into a slice of strings
	lines := strings.Split(wikitext, "\n")

	// now group the lines into sections
	var sections []Section
	var currentSection Section
	for _, line := range lines {
		if strings.HasPrefix(line, "==") {
			// commit the existing section and start a new section
			sections = append(sections, currentSection)
			currentSection.header = line
			currentSection.lines = nil
		} else {
			// ignore the following:
			// '----' (end of language marker) as this is inferred by a new language line, e.g. '==English=='
			// lines starting with '<!--' as these are HTML comments
			// lines starting with '[[Category' as we won't process these
			// blank lines
			if !strings.HasPrefix(line, "----") &&
				!strings.HasPrefix(line, "<!--") &&
				!strings.HasPrefix(line, "[[Category") &&
				len(line) > 0 {
				// write the line to the current section
				currentSection.lines = append(currentSection.lines, line)
			}
		}
	}
	// end of data - commit the last section
	sections = append(sections, currentSection)
	return sections
}

func extractLanguageSections(word string, langCode string, sections []Section) ([]Section, error) {
	// find the relevant language sections - this will be all of the sections starting with
	// ==Language== and up to (but not including) the next ==????== block, or the end of the data
	languageName := getLanguageFromCode(langCode)
	languageHeader := "==" + languageName + "=="
	startIndex, endIndex := 0, 0

	// find the start index
	for i := 0; i < len(sections); i++ {
		if strings.HasPrefix(sections[i].header, languageHeader) {
			startIndex = i
			break
		}
	}

	// if there was no start index, return an error
	if startIndex == 0 {
		msg := fmt.Sprintf("Word '%s' exists on Wiktionary, but not for %s", word, languageName)
		return nil, errors.New(msg)
	}

	// find the start of the next language, or the end of the file
	for i := startIndex + 1; i < len(sections); i++ {
		// match the next language-level heading
		re := regexp.MustCompile(`^==[^=]+==$`)
		if re.MatchString(sections[i].header) {
			endIndex = i
			break
		}
	}

	// if there was no next language block, set end index to the end of the array
	if endIndex == 0 {
		endIndex = len(sections)
	}

	langSections := make([]Section, endIndex-startIndex)
	copy(langSections, sections[startIndex:endIndex])

	return langSections, nil
}

func getLanguageFromCode(code string) string {
	// convert a language code to the full name, e.g. for "en" return "English"
	return languageCodes[code]
}

func getPageTitle(word string, langCode string) string {
	var title string
	// reconstructed words will have an asterisk as the first character and need special handling
	if strings.HasPrefix(word, "*") {
		reconLang := "Reconstruction:" + getLanguageFromCode(langCode) + "/"
		title = strings.Replace(word, "*", reconLang, 1)
	} else {
		// but normallly this is just the word
		title = word
	}
	return title

}

func getConvertedTextFromWiktionary(text string, word string, langCode string) (string, error) {
	// for the given text with tags, retrieve the equivalent text from Wiktionary
	urlHead := "https://en.wiktionary.org/w/api.php?action=parse&text="
	urlTail := "&prop=text&title=" + getPageTitle(word, langCode) + "&formatversion=2&format=json"

	// make an HHTP request to Wiktionary
	url := urlHead + url.QueryEscape(text) + urlTail
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// process the response to retrieve the body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	returnedText := string(body)

	// strip out the part of interest
	re := regexp.MustCompile(`text":"(.*?)</p.*>\\n<!--`)
	match := re.FindStringSubmatch(returnedText)
	if len(match) == 0 {
		// if there is no relevant text, return the original tags, so we at least have something
		return text, nil
	}

	// remove anything in <...> HTML braces
	re = regexp.MustCompile(`<(.*?)>`)
	convertedText := re.ReplaceAllString(match[1], "")

	// replace explicit spaces with actual spaces
	convertedText = strings.ReplaceAll(convertedText, "&#32;", " ")
	convertedText = strings.ReplaceAll(convertedText, "&nbsp", " ")
	convertedText = strings.ReplaceAll(convertedText, "&#160;", " ")
	convertedText = strings.ReplaceAll(convertedText, "\u0026#160;", " ")
	convertedText = strings.ReplaceAll(convertedText, "\u0026#8206;", " ") // left-to-right marker
	convertedText = strings.ReplaceAll(convertedText, "\u0026lt;", "<")
	convertedText = strings.ReplaceAll(convertedText, "\u003c", "<")

	// strip any newlines at the beginning or end
	convertedText = strings.Trim(convertedText, "\\n")

	return convertedText, nil

}
