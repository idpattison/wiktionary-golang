package main

import (
	"fmt"
	"regexp"
	"strings"
)

func parseSections(word string, langCode string, sections []Section) LanguageWord {
	// define the LanguageWord
	lw := LanguageWord{
		Word:         word,
		LanguageCode: langCode,
		LanguageName: getLanguageFromCode(langCode),
	}

	// iterate over the sections - ignore the first as it's the language header
	for i := 1; i < len(sections); i++ {
		parseSection(&lw, sections[i])
	}

	// assign a meaning - take the first entry in the first part of the first etymology
	if len(lw.Etymologies) > 0 {
		if len(lw.Etymologies[0].Parts) > 0 {
			if len(lw.Etymologies[0].Parts[0].Meanings) > 0 {
				lw.Meaning = lw.Etymologies[0].Parts[0].Meanings[0]
			}
		}
	}

	return lw
}

func parseSection(lw *LanguageWord, section Section) {
	// determine the section type
	sectionType := strings.Trim(section.header, "=")

	// process each type separately
	// etymology requires special handling as it may have numbers after it
	if strings.HasPrefix(sectionType, "Etymology") {
		parseEtymologySection(lw, section)
	} else {

		// process others
		switch sectionType {
		case "Pronunciation":
			parsePronunciationSection(lw, section)
		case "Noun", "Verb", "Adjective":
			parsePartofSpeechSection(lw, section)
		default:
			// ignore all others
		}

	}
}

func parsePronunciationSection(lw *LanguageWord, section Section) {
	var pr []string
	var ipa string
	// read each line - it should begin with a * - into the slice
	for _, line := range section.lines {
		if strings.HasPrefix(line, "*") {
			// find the first occurence of an IPA tag and record that separately
			if ipa == "" {
				ipaTag := searchForTag(line, "IPA")
				if ipaTag != "" {
					ipa = splitTag(ipaTag)["2"]
				}
			}
			// now process the pronunciation line
			text, _ := getConvertedTextFromWiktionary(line[2:], lw.Word)
			if text != "" {
				pr = append(pr, text)
			}
		}
	}

	// this section is usually added to the main language word
	// however if there are homographs (words spelled the same but pronounced differently)
	// then the pronunciation will be tied to the etymology
	// the way to tell is to see if there are any etymologies added yet - if so add to the latest one
	if len(lw.Etymologies) == 0 {
		lw.Pronunciations = pr
		lw.Ipa = ipa
	} else {
		lw.Etymologies[len(lw.Etymologies)-1].Pronunciations = pr
		lw.Etymologies[len(lw.Etymologies)-1].Ipa = ipa
	}
}

func parseEtymologySection(lw *LanguageWord, section Section) {
	var etym Etymology
	etym.Name = strings.Trim(section.header, "=")

	// read each line and process tags
	// we will concatenate the text with new lines as the etymology may span multiple paragraphs
	for _, line := range section.lines {

		// get the etymology text
		text, _ := getConvertedTextFromWiktionary(line, lw.Word)
		etym.Text += text
		if text != "" {
			etym.Text += "\n"
		}

		// get the word link tags from the etymology
		tags := getAllTags(line)
		for _, tag := range tags {
			var link LinkedWord
			link.SourceWord = lw.Word
			link.SourceLanguage = lw.LanguageCode

			elems := splitTag(tag[1])
			// ignore the m tag, it's sometimes used in etymologies, and it's ambiguous
			if elems["0"] == "m" || elems["0"] == "mention" {
				continue
			}

			// depending on the tag type, process accordingly
			switch elems["0"] {
			case "root":
				link.Relationship = Root
				link.TargetWord = elems["3"]
				link.TargetLanguage = elems["2"]
			case "inh", "inherited":
				link.Relationship = Inherited
				link.TargetWord = elems["3"]
				link.TargetLanguage = elems["2"]
				if val, ok := elems["5"]; ok {
					link.Meaning = val
				}
				// the meaning should appear at slot 4, but sometimes it's at slot 5
				// this is non-standard, but happens
				if val, ok := elems["6"]; ok && link.Meaning == "" {
					link.Meaning = val
				}
				if val, ok := elems["tr"]; ok {
					link.Transliteration = val
				}
			case "cog", "cognate":
				link.Relationship = Cognate
				link.TargetWord = elems["2"]
				link.TargetLanguage = elems["1"]
				if val, ok := elems["4"]; ok {
					link.Meaning = val
				}
				// the meaning should appear at slot 4, but sometimes it's at slot 5
				// this is non-standard, but happens
				if val, ok := elems["5"]; ok && link.Meaning == "" {
					link.Meaning = val
				}
				if val, ok := elems["tr"]; ok {
					link.Transliteration = val
				}
			}

			// if we have a word in a non-Latin script but no transliteration
			latinRe := regexp.MustCompile(`\p{Latin}`)
			if len(link.TargetWord) > 0 && !latinRe.MatchString(link.TargetWord) {
				re := regexp.MustCompile(link.TargetWord + ` *\((.*?)[\),]`)
				match := re.FindStringSubmatch(text)
				if len(match) > 1 {
					link.Transliteration = match[1]
				}
			}

			// if the target word exists, save it
			if link.TargetWord != "" && link.TargetWord != "-" {
				etym.Words = append(etym.Words, link)
			}
		}
	}

	lw.Etymologies = append(lw.Etymologies, etym)
}

func parsePartofSpeechSection(lw *LanguageWord, section Section) {
	var pos PartOfSpeech
	pos.Attributes = make(map[string]string)
	pos.Name = strings.Trim(section.header, "=")
	var headTag string

	// read each line and process tags
	for _, line := range section.lines {
		// the headword line will have tags
		if strings.HasPrefix(line, "{{") && pos.Headword == "" {
			text, _ := getConvertedTextFromWiktionary(line, lw.Word)
			headTag = line
			pos.Headword = text
		}
		// find meaning lines (but not quotations - maybe later)
		if strings.HasPrefix(line, "# ") {
			text, _ := getConvertedTextFromWiktionary(line[2:], lw.Word)
			pos.Meanings = append(pos.Meanings, text)
		}
	}

	// process attributes - depends on the part of speech type
	switch pos.Name {
	case "Noun":
		parseNoun(&pos, headTag)
	}

	if len(lw.Etymologies) > 0 {
		lw.Etymologies[len(lw.Etymologies)-1].Parts = append(lw.Etymologies[len(lw.Etymologies)-1].Parts, pos)
	}
}

func parseNoun(pos *PartOfSpeech, headTag string) {
	tagMap := splitTag(headTag)
	gendered := false

	// NB gender, if it exists, will be param 1 - it will be one of m f n c m-p f-p n-p c-p mf p
	if val, ok := tagMap["1"]; ok {
		valid := map[string]bool{
			"m": true, "f": true, "n": true, "c": true, "p": true, "m-p": true,
			"f-p": true, "n-p": true, "c-p": true, "mf": true, "m-f": true,
			"m-f-p": true, "mfp": true,
		}
		if valid[val] { // we have a gendered language with a gender as param 1
			gendered = true
			pos.Attributes["gender"] = val
		} else {
			gendered = false
		}
	}
	// there may also be an explicit gender as g= - check for this as well
	if val, ok := tagMap["g"]; ok {
		gendered = true
		pos.Attributes["gender"] = val
	}
	// handle the slightly quirky Danish situation
	if val, ok := tagMap["1"]; ok {
		if tagMap["0"] == "da-noun" {
			switch val {
			case "en", "n":
				pos.Attributes["gender"] = "c"
			case "et", "t":
				pos.Attributes["gender"] = "n"
			}
			gendered = true
		}
	}

	// get feminine and masculine forms
	if val, ok := tagMap["f"]; ok {
		pos.Attributes["feminine-form"] = val
	}
	if val, ok := tagMap["m"]; ok {
		pos.Attributes["masculine-form"] = val
	}

	// countability will be param 1 (or 2 for gendered languages) - it will be one of +, -, ~
	// it may also be marked as -|+
	firstCountableTag := "1"
	secondCountableTag := "2"
	if gendered {
		firstCountableTag = "2"
		secondCountableTag = "3"
	}
	if val, ok := tagMap[firstCountableTag]; ok {
		switch val {
		case "+", "s", "es": // takes care of an anomaly in English
			pos.Attributes["count"] = "countable"
		case "-":
			pos.Attributes["count"] = "uncountable"
			if val2, ok2 := tagMap[secondCountableTag]; ok2 {
				if val2 == "+" {
					pos.Attributes["count"] = "usually uncountable"
				}
			}
		case "~":
			pos.Attributes["count"] = "countable and uncountable"
		}
	}

	// get the headword forms from the text
	getHeadwordForm(pos, "singular definite")
	getHeadwordForm(pos, "singular indefinite")
	pluralDef := getHeadwordForm(pos, "plural definite")
	pluralInd := getHeadwordForm(pos, "plural indefinite")
	// if we have one of these specific plural forms, don't check for a normal plural
	if !pluralDef && !pluralInd {
		getHeadwordForm(pos, "plural")
	}
	getHeadwordForm(pos, "genitive")
	getHeadwordForm(pos, "diminutive")

}

func getHeadwordForm(pos *PartOfSpeech, form string) bool {
	// get the required form from the text
	re := regexp.MustCompile(form + ` *(.*?)[\),]`)
	match := re.FindStringSubmatch(pos.Headword)
	if len(match) > 1 {
		pos.Attributes[form] = match[1]
		return true
	}
	return false
}

func getAllTags(text string) [][]string {
	// return all wikitext tags in the text
	re := regexp.MustCompile(`\{\{(.*?)\}\}`)
	match := re.FindAllStringSubmatch(text, -1)
	return match
}

func searchForTag(text string, tag string) string {
	// return a tag of form {{head|param1|param2}} if it exists in the given text, otherwise ""
	re := regexp.MustCompile(`\{\{` + tag + `(.*?)\}\}`)
	match := re.FindStringSubmatch(text)
	if len(match) != 0 {
		return match[0]
	}
	return ""
}

func splitTag(tag string) map[string]string {
	// given a tag of form {{head|param1|param2}}, return a map of the components of the tag
	tagMap := make(map[string]string)
	text := strings.Trim(tag, "{}")
	tags := strings.Split(text, "|")
	for i, elem := range tags {
		if strings.Contains(elem, "=") {
			keyval := strings.Split(elem, "=")
			tagMap[keyval[0]] = keyval[1]
		} else {
			tagMap[fmt.Sprint(i)] = elem
		}
	}
	return tagMap
}
