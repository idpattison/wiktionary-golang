package main

import (
	"testing"

	"github.com/ianpattison-google/wiktionary-golang"
)

func TestMain(t *testing.T) {
	lw, err := wiktionary.GetWord("red", "en")
	if err != nil {
		t.Fatalf(`Error from GetWord: %q`, err)
	}
	expected := "Having red as its color."
	if lw.Meaning != expected {
		t.Fatalf(`lw.Meaning: expected %q, got %q`, expected, lw.Meaning)
	}
	expected = "English"
	if lw.LanguageName != expected {
		t.Fatalf(`lw.LanguageName: expected %q, got %q`, expected, lw.LanguageName)
	}
	expected = "/ɹɛd/"
	if lw.Ipa != expected {
		t.Fatalf(`lw.Ipa: expected %q, got %q`, expected, lw.Ipa)
	}
	if len(lw.Pronunciations) != 3 {
		t.Fatalf(`lw.Pronunciations: expected length 3, got %v`, len(lw.Pronunciations))
	}
	expected = "Rhymes: -ɛd"
	if lw.Pronunciations[2] != expected {
		t.Fatalf(`lw.Pronunciations[2]: expected %q, got %q`, expected, lw.Pronunciations[2])
	}
	if len(lw.Etymologies) != 3 {
		t.Fatalf(`lw.Etymologies: expected length 3, got %v`, len(lw.Etymologies))
	}
	expected = "Etymology 1"
	if lw.Etymologies[0].Name != expected {
		t.Fatalf(`lw.Etymologies[0].Name: expected %q, got %q`, expected, lw.Etymologies[2].Name)
	}
	if len(lw.Etymologies[0].Words) != 23 {
		t.Fatalf(`lw.Etymologies[0].Words: expected length 23, got %v`, len(lw.Etymologies[0].Words))
	}
	expType := "root"
	expLang := "ine-pro"
	expWord := "*h₁rewdʰ-"
	word := lw.Etymologies[0].Words[0]
	if word.Relationship != expType || word.Language != expLang || word.Word != expWord {
		t.Fatalf(`lw.Etymologies[0].Words[0]: expected type %q lang %q word %q, got %q %q %q`,
			expType, expLang, expWord, word.Relationship, word.Language, word.Word)
	}
	expType = "inherited"
	expLang = "enm"
	expWord = "red"
	word = lw.Etymologies[0].Words[1]
	if word.Relationship != expType || word.Language != expLang || word.Word != expWord {
		t.Fatalf(`lw.Etymologies[0].Words[1]: expected type %q lang %q word %q, got %q %q %q`,
			expType, expLang, expWord, word.Relationship, word.Language, word.Word)
	}
	expType = "cognate"
	expLang = "fy"
	expWord = "read"
	word = lw.Etymologies[0].Words[5]
	if word.Relationship != expType || word.Language != expLang || word.Word != expWord {
		t.Fatalf(`lw.Etymologies[0].Words[5]: expected type %q lang %q word %q, got %q %q %q`,
			expType, expLang, expWord, word.Relationship, word.Language, word.Word)
	}

	// test the function to pull transliterations from the text version
	expLang = "grc"
	expWord = "ἐρυθρός"
	expTrans := "eruthrós"
	word = lw.Etymologies[0].Words[16]
	if word.Language != expLang || word.Word != expWord || word.Transliteration != expTrans {
		t.Fatalf(`lw.Etymologies[0].Words[16]: expected lang %q word %q trans %q, got %q %q %q`,
			expLang, expWord, expTrans, word.Language, word.Word, word.Transliteration)
	}

	// test translations
	expLang = "af"
	expWord = "rooi"
	trans := lw.Etymologies[0].Parts[1].Translations[0]
	if trans.Language != expLang || trans.Word != expWord {
		t.Fatalf(`lw.Etymologies[0].Parts[1].Translations[0]: expected lang %q word %q, got %q %q`,
			expLang, expWord, trans.Language, trans.Word)
	}

	// test parts of speech
	if len(lw.Etymologies[0].Parts) != 2 {
		t.Fatalf(`lw.Etymologies[0].Parts: expected length 2, got %v`, len(lw.Etymologies[0].Parts))
	}
	expPart := "Adjective"
	expHead := "red (comparative redder or more red, superlative reddest or most red)"
	part := lw.Etymologies[0].Parts[0]
	if part.Name != expPart || part.Headword != expHead {
		t.Fatalf(`lw.Etymologies[0].Parts[0]: expected %q - %q, got %q - %q`,
			expPart, expHead, part.Name, part.Headword)
	}
	if len(part.Meanings) != 7 {
		t.Fatalf(`lw.Etymologies[0].Parts[0].Meanings: expected length 7, got %v`, len(part.Meanings))
	}
	if len(part.Attributes) != 2 {
		t.Fatalf(`lw.Etymologies[0].Parts[0].Attributes: expected length 2, got %v`, len(part.Attributes))
	} else {
		if val, ok := part.Attributes["comparative"]; ok {
			expected = "redder or more red"
			if val != expected {
				t.Fatalf(`lw.Etymologies[0].Parts[0].Attributes["comparative"]: expected %q, got %q`, expected, part.Attributes["comparative"])
			}
		}
	}
	expected = "(particle physics) Having a color charge of red."
	if part.Meanings[6] != expected {
		t.Fatalf(`lw.Etymologies[0].Parts[0].Meanings[6]: expected %q, got %q`, expected, part.Meanings[6])
	}
	expPart = "Noun"
	expHead = "red (countable and uncountable, plural reds)"
	part = lw.Etymologies[0].Parts[1]
	if part.Name != expPart || part.Headword != expHead {
		t.Fatalf(`lw.Etymologies[0].Parts[1]: expected %q - %q, got %q - %q`,
			expPart, expHead, part.Name, part.Headword)
	}
	expIpa := "/ɹɛd/"
	if lw.Ipa != expIpa {
		t.Fatalf(`lw.Ipa: expected %v, got %v`, expIpa, lw.Ipa)
	}
	expAntonyms := "(having red as its colour): nonred, unred\n(having red as its colour charge): antired\n"
	part = lw.Etymologies[0].Parts[0]
	if part.Antonyms != expAntonyms {
		t.Fatalf(`lw.Etymologies[0].Parts[0]: expected antonyms %q, got %q`, expAntonyms, part.Antonyms)
	}

}

func TestMainPartial(t *testing.T) {
	// using the core settings, we should get a meaning but not a translation
	var options wiktionary.WiktionaryOptions
	options.RequiredSections = wiktionary.Sec_Core
	options.RequiredLanguages = []string{"all"}

	lw, err := wiktionary.GetWordWithOptions("red", "en", options)
	if err != nil {
		t.Fatalf(`Error from GetWord: %q`, err)
	}
	expected := "Having red as its color."
	if lw.Meaning != expected {
		t.Fatalf(`lw.Meaning: expected %q, got %q`, expected, lw.Meaning)
	}
	part := lw.Etymologies[0].Parts[1]
	if len(part.Translations) != 0 {
		t.Fatalf(`lw.Etymologies[0].Parts[1].Translations: expected length 0, got %v`, len(part.Translations))
	}

}

func TestMainFrench(t *testing.T) {
	lw, err := wiktionary.GetWord("rouge", "fr")
	if err != nil {
		t.Fatalf(`Error from GetWord: %q`, err)
	}
	if len(lw.Etymologies[0].Parts) != 2 {
		t.Fatalf(`lw.Etymologies[0].Parts: expected length 2, got %v`, len(lw.Etymologies[0].Parts))
	}
	expPart := "Adjective"
	expHead := "rouge (plural rouges)"
	part := lw.Etymologies[0].Parts[0]
	if part.Name != expPart || part.Headword != expHead {
		t.Fatalf(`lw.Etymologies[0].Parts[0]: expected %q - %q, got %q - %q`,
			expPart, expHead, part.Name, part.Headword)
	}
}

func TestMainOldEnglish(t *testing.T) {
	lw, err := wiktionary.GetWord("grene", "ang")
	if err != nil {
		t.Fatalf(`Error from GetWord: %q`, err)
	}
	if len(lw.Etymologies[0].Parts) != 1 {
		t.Fatalf(`lw.Etymologies[0].Parts: expected length 1, got %v`, len(lw.Etymologies[0].Parts))
	}
	expPart := "Adjective"
	expHead := "grēne"
	part := lw.Etymologies[0].Parts[0]
	if part.Name != expPart || part.Headword != expHead {
		t.Fatalf(`lw.Etymologies[0].Parts[0]: expected %q - %q, got %q - %q`,
			expPart, expHead, part.Name, part.Headword)
	}
}

func TestMainProtoGermanic(t *testing.T) {
	lw, err := wiktionary.GetWord("*raudaz", "gem-pro")
	if err != nil {
		t.Fatalf(`Error from GetWord: %q`, err)
	}
	if len(lw.Etymologies[0].Parts) != 1 {
		t.Fatalf(`lw.Etymologies[0].Parts: expected length 1, got %v`, len(lw.Etymologies[0].Parts))
	}
	part := lw.Etymologies[0].Parts[0]
	if len(part.Meanings) != 1 {
		t.Fatalf(`lw.Etymologies[0].Parts[0].Meanings: expected length 1, got %v`, len(part.Meanings))
	}
	if len(part.Attributes) != 2 {
		t.Fatalf(`lw.Etymologies[0].Parts[0].Attributes: expected length 2, got %v`, len(part.Attributes))
	} else {
		if val, ok := part.Attributes["superlative"]; ok {
			expected := "*raudōstaz"
			if val != expected {
				t.Fatalf(`lw.Etymologies[0].Parts[0].Attributes["superlative"]: expected %q, got %q`, expected, part.Attributes["comparative"])
			}
		}
	}
}
