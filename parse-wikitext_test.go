package wiktionary

import "testing"

func TestParseNoun(t *testing.T) {
	var pos PartOfSpeech
	pos.Attributes = make(map[string]string)
	pos.Name = "Noun"

	// test the basic English form
	pos.Headword = "church (plural churches)"
	testAttributes(t, &pos, "{{en-noun|s}}", "count", "countable", true)

	// test plurals
	testAttributes(t, &pos, "{{en-noun|s}}", "plural", "churches", true)

	// test countabilities
	pos.Headword = "red (countable and uncountable, plural reds)"
	testAttributes(t, &pos, "{{en-noun|~}}", "count", "countable and uncountable", true)

	// test French
	pos.Headword = "chien m (plural chiens, feminine chienne)"
	testAttributes(t, &pos, "{{fr-noun|m|f=chienne}}", "plural", "chiens", true)
	testAttributes(t, &pos, "{{fr-noun|m|f=chienne}}", "gender", "m", true)
	testAttributes(t, &pos, "{{fr-noun|m|f=chienne}}", "feminine-form", "chienne", true)

	// test German
	pos.Headword = "Buch n (genitive Buchs or Buches, plural Bücher, diminutive Büchlein n)"
	testAttributes(t, &pos, "{{de-noun|n|Buchs|gen2=Buches|Bücher|Büchlein}}", "plural", "Bücher", true)
	testAttributes(t, &pos, "{{de-noun|n|Buchs|gen2=Buches|Bücher|Büchlein}}", "gender", "n", true)
	testAttributes(t, &pos, "{{de-noun|n|Buchs|gen2=Buches|Bücher|Büchlein}}", "genitive", "Buchs or Buches", true)
	testAttributes(t, &pos, "{{de-noun|n|Buchs|gen2=Buches|Bücher|Büchlein}}", "diminutive", "Büchlein n", true)

	// test Danish
	pos.Headword = "stol c (singular definite stolen, plural indefinite stole)"
	testAttributes(t, &pos, "{{da-noun|en|e|e}}", "plural indefinite", "stole", true)
	testAttributes(t, &pos, "{{da-noun|en|e|e}}", "gender", "c", true)

	// test Dutch
	pos.Headword = "artikel n (plural artikelen or artikels, diminutive artikeltje n)"
	testAttributes(t, &pos, "{{nl-noun|n|-@en|pl2=-s|artikeltje}}", "plural", "artikelen or artikels", true)
	testAttributes(t, &pos, "{{nl-noun|n|-@en|pl2=-s|artikeltje}}", "gender", "n", true)
	testAttributes(t, &pos, "{{nl-noun|n|-@en|pl2=-s|artikeltje}}", "diminutive", "artikeltje n", true)
}

func TestParseNounDeclension(t *testing.T) {
	var pos PartOfSpeech
	pos.Attributes = make(map[string]string)
	pos.Name = "Noun"

	// test the basic English form
	pos.Headword = "church (plural churches)"
	testAttributes(t, &pos, "{{en-noun|s}}", "count", "countable", true)
}

func TestParseAdjective(t *testing.T) {
	var pos PartOfSpeech
	pos.Attributes = make(map[string]string)
	pos.Name = "Adjective"

	// test the basic English form
	pos.Headword = "late (comparative later, superlative latest)"
	testAttributes(t, &pos, "{{en-adj|more}}", "comparative", "later", true)

	// test incomparable forms
	pos.Headword = "annual (not comparable)"
	testAttributes(t, &pos, "{{en-adj|-}}", "comparative", "", false)

	// test French
	pos.Headword = "compact (feminine singular compacte, masculine plural compacts, feminine plural compactes)"
	testAttributes(t, &pos, "{{fr-adj}}", "masculine plural", "compacts", true)

	// test German
	pos.Headword = "entenartig (comparative entenartiger, superlative am entenartigsten)"
	testAttributes(t, &pos, "{{de-adj|er|sten}}", "superlative", "am entenartigsten", true)
}

func TestParseVerb(t *testing.T) {
	var pos PartOfSpeech
	pos.Attributes = make(map[string]string)
	pos.Name = "Verb"

	// test the basic English form
	pos.Headword = "flip (third-person singular simple present flips, present participle flipping, simple past and past participle flipped)"
	testAttributes(t, &pos, "{{en-verb}}", "third-person singular simple present", "flips", true)
	testAttributes(t, &pos, "{{en-verb}}", "present participle", "flipping", true)
	testAttributes(t, &pos, "{{en-verb}}", "simple past and past participle", "flipped", true)

	// test irregular forms with missing attributes
	pos.Headword = "can (third-person singular simple present can, no present participle, simple past could, no past participle)"
	testAttributes(t, &pos, "{{en-verb}}", "past participle", "", false)

	// // test German
	pos.Headword = "laufen (class 7 strong, third-person singular present läuft, past tense lief, past participle gelaufen, auxiliary sein)"
	testAttributes(t, &pos, "{{de-verb|laufen<läuft#lief,gelaufen.sein>}}", "type", "class 7 strong", true)
	testAttributes(t, &pos, "{{de-verb|laufen<läuft#lief,gelaufen.sein>}}", "past tense", "lief", true)
}

func testAttributes(t *testing.T, pos *PartOfSpeech, headTag string, attr string, expected string, shouldExist bool) {
	for k := range pos.Attributes {
		delete(pos.Attributes, k)
	}
	switch pos.Name {
	case "Noun":
		parseNoun(pos, headTag)
	case "Adjective":
		parseAdjective(pos, headTag)
	case "Verb":
		parseVerb(pos, headTag)
	}
	if val, ok := pos.Attributes[attr]; ok {
		if shouldExist {
			if val != expected {
				t.Fatalf(`parse%s: expected %q, got %q`, pos.Name, expected, val)
			}
		} else {
			t.Fatalf(`parse%s: %v attribute found for %q; should not exist`, pos.Name, attr, headTag)
		}
	} else {
		if shouldExist {
			t.Fatalf(`parse%s: no %v attribute for %q`, pos.Name, attr, headTag)
		}
	}
}

// func testExtendedAttributes(t *testing.T, pos *PartOfSpeech, headTag string, attr string, expected string, shouldExist bool) {
// 	for k := range pos.Attributes {
// 		delete(pos.Attributes, k)
// 	}
// 	parseExtendedPartSection()
// 	if val, ok := pos.Attributes[attr]; ok {
// 		if shouldExist {
// 			if val != expected {
// 				t.Fatalf(`parse%s: expected %q, got %q`, pos.Name, expected, val)
// 			}
// 		} else {
// 			t.Fatalf(`parse%s: %v attribute found for %q; should not exist`, pos.Name, attr, headTag)
// 		}
// 	} else {
// 		if shouldExist {
// 			t.Fatalf(`parse%s: no %v attribute for %q`, pos.Name, attr, headTag)
// 		}
// 	}
// }

func TestHtmlTableProcessing(t *testing.T) {
	var pos PartOfSpeech
	pos.Attributes = make(map[string]string)
	pos.Name = "Noun"

	pos.Headword = "lūna f (genitive lūnae); first declension"
	testAttributes(t, &pos, "{{la-ndecl|l\u016bna<1>}}", "genitive", "lūnae", true)
}
