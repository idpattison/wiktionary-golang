package main

import "testing"

func TestGetWikitext(t *testing.T) {
	rawData := `{"parse":{"title":"red","pageid":3654,"wikitext":{"*":"{{also|-red|red-|Red|RED|r\u011bd}}"}}}`
	bytes := []byte(rawData)
	expected := "{{also|-red|red-|Red|RED|r\u011bd}}"
	wikitext, err := getWikitext(bytes, "red")
	if err != nil {
		t.Fatalf(`Error from getWikitext: %q`, err)
	}
	if wikitext != expected {
		t.Fatalf(`getWikitext: expected %q, got %q`, expected, wikitext)
	}
}

func TestConvertWikitext(t *testing.T) {
	inputData := `{{IPA|en|/\u0279\u025bd/|[\u027b\u02b7\u025b\u02d1d\u0325]}}`
	expected := "{{IPA|en|/ɹɛd/|[ɻʷɛˑd̥]}}"
	outputData, err := convertWikitext(inputData)
	if err != nil {
		t.Fatalf(`Error from convertWikitext: %q`, err)
	}
	if outputData != expected {
		t.Fatalf(`convertWikitext: expected %q, got %q`, expected, outputData)
	}
}

func TestGetConvertedTextFromWiktionary(t *testing.T) {
	inputData := `{{en-adj|redder|more}}`
	expected := "red (comparative redder or more red, superlative reddest or most red)"
	outputData, err := getConvertedTextFromWiktionary(inputData, "red", "en")
	if err != nil {
		t.Fatalf(`Error from getConvertedTextFromWiktionary: %q`, err)
	}
	if outputData != expected {
		t.Fatalf(`getConvertedTextFromWiktionary: expected %q, got %q`, expected, outputData)
	}
}

func TestGetLanguageFromCode(t *testing.T) {
	lang := getLanguageFromCode("de")
	if lang != "German" {
		t.Fatalf(`getLanguageFromCode: expected %q, got %q`, "German", lang)
	}
	lang = getLanguageFromCode("ine-pro")
	if lang != "Proto-Indo-European" {
		t.Fatalf(`getLanguageFromCode: expected %q, got %q`, "Proto-Indo-European", lang)
	}
}
