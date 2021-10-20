package main

import "testing"

func TestMain(t *testing.T) {
	lw, err := processWord("red", "en")
	if err != nil {
		t.Fatalf(`Error from processWord: %q`, err)
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
	if word.Relationship != expType || word.TargetLanguage != expLang || word.TargetWord != expWord {
		t.Fatalf(`lw.Etymologies[0].Words[0]: expected type %q lang %q word %q, got %q %q %q`,
			expType, expLang, expWord, word.Relationship, word.TargetLanguage, word.TargetWord)
	}
	expType = "inherited"
	expLang = "enm"
	expWord = "red"
	word = lw.Etymologies[0].Words[1]
	if word.Relationship != expType || word.TargetLanguage != expLang || word.TargetWord != expWord {
		t.Fatalf(`lw.Etymologies[0].Words[1]: expected type %q lang %q word %q, got %q %q %q`,
			expType, expLang, expWord, word.Relationship, word.TargetLanguage, word.TargetWord)
	}
	expType = "cognate"
	expLang = "fy"
	expWord = "read"
	word = lw.Etymologies[0].Words[5]
	if word.Relationship != expType || word.TargetLanguage != expLang || word.TargetWord != expWord {
		t.Fatalf(`lw.Etymologies[0].Words[5]: expected type %q lang %q word %q, got %q %q %q`,
			expType, expLang, expWord, word.Relationship, word.TargetLanguage, word.TargetWord)
	}

	// test the function to pull transliterations from the text version
	expLang = "grc"
	expWord = "ἐρυθρός"
	expTrans := "eruthrós"
	word = lw.Etymologies[0].Words[16]
	if word.TargetLanguage != expLang || word.TargetWord != expWord || word.Transliteration != expTrans {
		t.Fatalf(`lw.Etymologies[0].Words[16]: expected lang %q word %q trans %q, got %q %q %q`,
			expLang, expWord, expTrans, word.TargetLanguage, word.TargetWord, word.Transliteration)
	}

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

}

func TestMainFrench(t *testing.T) {
	lw, err := processWord("rouge", "fr")
	if err != nil {
		t.Fatalf(`Error from processWord: %q`, err)
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
