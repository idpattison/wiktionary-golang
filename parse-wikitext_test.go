package main

import "testing"

func TestParseNoun(t *testing.T) {
	var pos PartOfSpeech
	pos.Attributes = make(map[string]string)
	pos.Name = "Noun"

	// test the basic English form
	pos.Headword = "church (plural churches)"
	testAttributes(t, &pos, "{{en-noun|s}}", "count", "countable")

	// test plurals
	testAttributes(t, &pos, "{{en-noun|s}}", "plural", "churches")

	// test countabilities
	pos.Headword = "red (countable and uncountable, plural reds)"
	testAttributes(t, &pos, "{{en-noun|~}}", "count", "countable and uncountable")

	// test French
	pos.Headword = "chien m (plural chiens, feminine chienne)"
	testAttributes(t, &pos, "{{fr-noun|m|f=chienne}}", "plural", "chiens")
	testAttributes(t, &pos, "{{fr-noun|m|f=chienne}}", "gender", "m")
	testAttributes(t, &pos, "{{fr-noun|m|f=chienne}}", "feminine-form", "chienne")

	// test German
	pos.Headword = "Buch n (genitive Buchs or Buches, plural Bücher, diminutive Büchlein n)"
	testAttributes(t, &pos, "{{de-noun|n|Buchs|gen2=Buches|Bücher|Büchlein}}", "plural", "Bücher")
	testAttributes(t, &pos, "{{de-noun|n|Buchs|gen2=Buches|Bücher|Büchlein}}", "gender", "n")
	testAttributes(t, &pos, "{{de-noun|n|Buchs|gen2=Buches|Bücher|Büchlein}}", "genitive", "Buchs or Buches")
	testAttributes(t, &pos, "{{de-noun|n|Buchs|gen2=Buches|Bücher|Büchlein}}", "diminutive", "Büchlein n")

	// test Danish
	pos.Headword = "stol c (singular definite stolen, plural indefinite stole)"
	testAttributes(t, &pos, "{{da-noun|en|e|e}}", "plural indefinite", "stole")
	testAttributes(t, &pos, "{{da-noun|en|e|e}}", "gender", "c")

	// test Dutch
	pos.Headword = "artikel n (plural artikelen or artikels, diminutive artikeltje n)"
	testAttributes(t, &pos, "{{nl-noun|n|-@en|pl2=-s|artikeltje}}", "plural", "artikelen or artikels")
	testAttributes(t, &pos, "{{nl-noun|n|-@en|pl2=-s|artikeltje}}", "gender", "n")
	testAttributes(t, &pos, "{{nl-noun|n|-@en|pl2=-s|artikeltje}}", "diminutive", "artikeltje n")
}

func testAttributes(t *testing.T, pos *PartOfSpeech, headTag string, attr string, expected string) {
	for k := range pos.Attributes {
		delete(pos.Attributes, k)
	}
	parseNoun(pos, headTag)
	if val, ok := pos.Attributes[attr]; ok {
		if val != expected {
			t.Fatalf(`parseNoun: expected %q, got %q`, expected, val)
		}
	} else {
		t.Fatalf(`parseNoun: no %v attribute for %q`, attr, headTag)
	}
}
