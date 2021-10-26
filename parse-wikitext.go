package wiktionary

import (
	"fmt"
	"regexp"
	"strings"
)

func parseSections(word string, langCode string, sections []Section, options WiktionaryOptions) LanguageWord {
	// define the LanguageWord
	lw := LanguageWord{
		Word:         word,
		LanguageCode: langCode,
		LanguageName: getLanguageFromCode(langCode),
	}

	// iterate over the sections - ignore the first as it's the language header
	for i := 1; i < len(sections); i++ {
		parseSection(&lw, sections[i], options)
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

func parseSection(lw *LanguageWord, section Section, options WiktionaryOptions) {
	// determine the section type
	sectionType := strings.Trim(section.header, "=")

	// process each type separately
	// etymology requires special handling as it may have numbers after it
	if strings.HasPrefix(sectionType, "Etymology") {
		parseEtymologySection(lw, section, options)
	} else {

		// process others
		switch sectionType {
		case "Pronunciation":
			if sectionRequired(options, Sec_IPA) || sectionRequired(options, Sec_Extended_Pronunciation) {
				parsePronunciationSection(lw, section, options)
			}
		case "Noun", "Verb", "Adjective":
			if sectionRequired(options, Sec_Parts) {
				parsePartofSpeechSection(lw, section, options)
			}
		case "Translations":
			if sectionRequired(options, Sec_Translations) {
				parseTranslationSection(lw, section, options)
			}
		case "Descendants":
			if sectionRequired(options, Sec_Etymology_Words) {
				parseDescendantSection(lw, section, options)
			}
		case "Synonyms", "Antonyms", "Anagrams", "Alternative forms":
			if sectionRequired(options, Sec_Synonyms) || sectionRequired(options, Sec_Antonyms) {
				parseOtherSections(lw, section, options)
			}
		default:
			// ignore all others
		}

	}
}

func parsePronunciationSection(lw *LanguageWord, section Section, options WiktionaryOptions) {
	var pr []string
	var ipa string
	var text string
	// read each line - it should begin with a * - into the slice
	for _, line := range section.lines {
		if strings.HasPrefix(line, "*") {
			// process the pronunciation line
			// special handling for the audio line
			prTag := searchForTag(line, "audio")
			if prTag != "" {
				elems := splitTag(prTag)
				if val, ok := elems["2"]; ok {
					pr = append(pr, "Audio: https://en.wiktionary.org/wiki/File:"+val)
				}
				continue
			}

			text, _ = getConvertedTextFromWiktionary(line[2:], lw.Word, lw.LanguageCode)
			if sectionRequired(options, Sec_Extended_Pronunciation) {
				if text != "" {
					pr = append(pr, text)
				}
			}
			// find the first occurence of an IPA tag and record that separately
			// NB some languages have an automatically generated IPA - for simplicity
			// we will always use the text version rther than the tag
			if sectionRequired(options, Sec_IPA) {
				if ipa == "" {
					re := regexp.MustCompile(`.*(\/.*?\/)`)
					match := re.FindStringSubmatch(text)
					if len(match) > 0 {
						ipa = match[1]
					}
				}
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

func parseEtymologySection(lw *LanguageWord, section Section, options WiktionaryOptions) {
	var etym Etymology
	etym.Name = strings.Trim(section.header, "=")

	// read each line and process tags
	// we will concatenate the text with new lines as the etymology may span multiple paragraphs
	for _, line := range section.lines {

		// get the etymology text
		if sectionRequired(options, Sec_Etymology_Text) || sectionRequired(options, Sec_Etymology_Words) {
			text, _ := getConvertedTextFromWiktionary(line, lw.Word, lw.LanguageCode)
			if sectionRequired(options, Sec_Etymology_Text) {
				etym.Text += text
			}
			if text != "" {
				etym.Text += "\n"
			}

			// get the word link tags from the etymology
			if sectionRequired(options, Sec_Etymology_Words) {
				tags := getAllTags(line)
				for _, tag := range tags {
					var link LinkedWord

					elems := splitTag(tag[1])
					// ignore the m tag, it's sometimes used in etymologies, and it's ambiguous
					if elems["0"] == "m" || elems["0"] == "mention" {
						continue
					}

					// depending on the tag type, process accordingly
					switch elems["0"] {
					case "root":
						link.Relationship = Root
						link.Word = elems["3"]
						link.Language = elems["2"]
					case "inh", "inherited":
						link.Relationship = Inherited
						link.Word = elems["3"]
						link.Language = elems["2"]
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
						link.Word = elems["2"]
						link.Language = elems["1"]
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
					if len(link.Word) > 0 && !latinRe.MatchString(link.Word) {
						re := regexp.MustCompile(link.Word + ` *\((.*?)[\),]`)
						match := re.FindStringSubmatch(text)
						if len(match) > 1 {
							link.Transliteration = match[1]
						}
					}

					// if the target word exists, save it
					if link.Word != "" && link.Word != "-" {
						etym.Words = append(etym.Words, link)
					}
				}
			}
		}
	}

	lw.Etymologies = append(lw.Etymologies, etym)
}

func parsePartofSpeechSection(lw *LanguageWord, section Section, options WiktionaryOptions) {
	var pos PartOfSpeech
	pos.Attributes = make(map[string]string)
	pos.Name = strings.Trim(section.header, "=")
	var headTag string

	// read each line and process tags
	for _, line := range section.lines {
		// the headword line will have tags
		if strings.HasPrefix(line, "{{") && pos.Headword == "" {
			text, _ := getConvertedTextFromWiktionary(line, lw.Word, lw.LanguageCode)
			headTag = line
			pos.Headword = text
		}
		// find meaning lines (but not quotations - maybe later)
		if sectionRequired(options, Sec_Meanings) {
			if strings.HasPrefix(line, "# ") {
				text, _ := getConvertedTextFromWiktionary(line[2:], lw.Word, lw.LanguageCode)
				pos.Meanings = append(pos.Meanings, text)
			}
		}
	}

	// process attributes - depends on the part of speech type
	if sectionRequired(options, Sec_Part_Attributes) {
		switch pos.Name {
		case "Noun":
			parseNoun(&pos, headTag)
		case "Adjective", "Adverb": // adverbs are largely treated the same as adjectives
			parseAdjective(&pos, headTag)
		case "Verb":
			parseVerb(&pos, headTag)
			// other parts of speech suchas conjunctions are generally simpler
			// they may have attributes but we will rely on the headword text
		}
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

func parseAdjective(pos *PartOfSpeech, headTag string) {
	// get the headword forms from the text
	femSing := getHeadwordForm(pos, "feminine singular")
	mascSing := getHeadwordForm(pos, "masculine singular")
	femPlural := getHeadwordForm(pos, "feminine plural")
	mascPlural := getHeadwordForm(pos, "masculine plural")
	// if we have one of these specific adjective forms, don't check for a normal plural
	if !femSing && !mascSing && !femPlural && !mascPlural {
		getHeadwordForm(pos, "plural")
	}
	getHeadwordForm(pos, "comparative")
	getHeadwordForm(pos, "superlative")
}

func parseVerb(pos *PartOfSpeech, headTag string) {
	tagMap := splitTag(headTag)
	// get the headword forms from the text
	sppp := getHeadwordForm(pos, "simple past and past participle")
	// if we have this specific combined form, don't check for separate forms
	if !sppp {
		getHeadwordForm(pos, "simple past")
		getHeadwordForm(pos, "past participle")
	}
	tpssp := getHeadwordForm(pos, "third-person singular simple present")
	tpsps := getHeadwordForm(pos, "third-person singular present")
	fpsp := getHeadwordForm(pos, "first-person singular present")
	pt := getHeadwordForm(pos, "present tense")
	if !tpssp && !tpsps && !fpsp && !pt {
		getHeadwordForm(pos, "present")
	}
	fpspt := getHeadwordForm(pos, "first-person singular preterite")
	if !fpspt {
		getHeadwordForm(pos, "preterite")
	}
	getHeadwordForm(pos, "present participle")
	getHeadwordForm(pos, "past tense")
	getHeadwordForm(pos, "past subjunctive")
	getHeadwordForm(pos, "perfect tense")
	getHeadwordForm(pos, "imperitive")
	getHeadwordForm(pos, "infinitive")
	getHeadwordForm(pos, "auxiliary")
	getHeadwordForm(pos, "type")  // French verbs use this for auxiliary, defective etc
	if tagMap["0"] == "de-verb" { // first headword item in German verbs is the type
		getHeadwordItem(pos, "type", 0)
	}

}

func getHeadwordForm(pos *PartOfSpeech, form string) bool {
	// get the required form from the text
	re := regexp.MustCompile(form + ` *(.*?)[\),]`)
	match := re.FindStringSubmatch(pos.Headword)
	if len(match) > 1 {
		if len(match[1]) > 0 {
			pos.Attributes[form] = match[1]
			return true
		}
	}
	return false
}

func getHeadwordItem(pos *PartOfSpeech, form string, index int) bool {
	// get the part of the headword in braces then split by commas
	re := regexp.MustCompile(`.*\((.*?)\)$`)
	match := re.FindStringSubmatch(pos.Headword)
	if len(match) != 0 {
		items := strings.Split(match[1], ",")
		if len(items) > index {
			if len(items[index]) > 0 {
				pos.Attributes[form] = items[index]
				return true
			}
		}
	}
	return false
}

func parseTranslationSection(lw *LanguageWord, section Section, options WiktionaryOptions) {
	var tr []TranslatedWord

	// NB we will only record the first translation block as this will be the principal meaning
	// many words hav translations of colloquial meanings which are not relevant
	for _, line := range section.lines {
		// read lines on until we reach the first {{trans-bottom}} tag, then exit
		if strings.HasPrefix(line, "{{trans-bottom") {
			break
		}
		if strings.HasPrefix(line, "*") {
			tags := getAllTags(line)
			for _, tag := range tags {
				elems := splitTag(tag[1])
				// if we have a translation tag
				if elems["0"] == "t" || elems["0"] == "t+" || elems["0"] == "tt" || elems["0"] == "tt+" {
					var tw TranslatedWord
					if val, ok := elems["1"]; ok {
						tw.Language = val
					}
					if val, ok := elems["2"]; ok {
						tw.Word = val
						if val, ok := elems["tr"]; ok {
							tw.Transliteration = val
						}

						tr = append(tr, tw)
					}
				}
			}
		}
	}

	// add this section to the current part in the current etymology
	if len(lw.Etymologies) > 0 && len(tr) > 0 {
		currentEtym := len(lw.Etymologies) - 1
		if len(lw.Etymologies[currentEtym].Parts) > 0 {
			currentPart := len(lw.Etymologies[len(lw.Etymologies)-1].Parts) - 1
			lw.Etymologies[currentEtym].Parts[currentPart].Translations = tr
		}
	}
}

func parseDescendantSection(lw *LanguageWord, section Section, options WiktionaryOptions) {
	// read in all descendant words and add them to LinkedWords in the current Etymology
	for _, line := range section.lines {
		if strings.HasPrefix(line, "*") {
			descTag := searchForTag(line, "desc") // NB this will also pick up desctree
			if descTag != "" {
				// decode the tag
				var link LinkedWord
				elems := splitTag(descTag)
				if val, ok := elems["1"]; ok {
					link.Language = val
					if val, ok := elems["2"]; ok {
						link.Word = val
						if _, ok := elems["bor"]; ok {
							link.Relationship = Descendant
						}
						// add to the current etymology
						if len(lw.Etymologies) > 0 {
							lw.Etymologies[len(lw.Etymologies)-1].Words = append(lw.Etymologies[len(lw.Etymologies)-1].Words, link)
						}
					}
				}
			}
		}
	}

}

func parseOtherSections(lw *LanguageWord, section Section, options WiktionaryOptions) {
	// used for synonyms, antonyms and other sections where we just return the text
	var secText string
	header := strings.Trim(section.header, "=")
	for _, line := range section.lines {
		if strings.HasPrefix(line, "*") {
			if (sectionRequired(options, Sec_Synonyms) && header == "Synonyms") ||
				(sectionRequired(options, Sec_Antonyms) && header == "Antonyms") ||
				(sectionRequired(options, Sec_Anagrams) && header == "Anagrams") ||
				(sectionRequired(options, Sec_Alternatives) && header == "Alternative forms") {
				text, _ := getConvertedTextFromWiktionary(line[2:], lw.Word, lw.LanguageCode)
				if text != "" {
					secText += text + "\n"
				}
			}
		}
	}

	// sections that live at the LanguageWord level
	if header == "Anagrams" {
		if sectionRequired(options, Sec_Anagrams) && secText != "" {
			lw.Anagrams = secText
		}
		return
	}

	// sections that live at the Etymology level
	currentEtym := len(lw.Etymologies) - 1
	if currentEtym < 0 {
		return
	}
	if header == "Alternative forms" {
		if sectionRequired(options, Sec_Alternatives) && secText != "" {
			lw.Etymologies[currentEtym].AlternativeForms = secText
		}
		return
	}

	// sections that live at the PartOfSpeech level
	currentPart := len(lw.Etymologies[currentEtym].Parts) - 1
	if currentPart < 0 {
		return
	}
	if secText != "" {
		switch header {
		case "Synonyms":
			lw.Etymologies[currentEtym].Parts[currentPart].Synonyms = secText
		case "Antonyms":
			lw.Etymologies[currentEtym].Parts[currentPart].Antonyms = secText
		}
	}
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

func sectionRequired(options WiktionaryOptions, section int16) bool {
	b := options.RequiredSections&section > 0
	return b
}
