package wiktionary

import "testing"

func TestExtendedParts(t *testing.T) {
	lw, err := GetWord("mando", "la")
	if err != nil {
		t.Fatalf(`Error from GetWord: %q`, err)
	}
	if len(lw.Etymologies[0].Parts) != 1 {
		t.Fatalf(`lw.Etymologies[0].Parts: expected length 1, got %v`, len(lw.Etymologies[0].Parts))
	}
	part := lw.Etymologies[0].Parts[0]
	if len(part.Attributes) < 150 {
		t.Fatalf(`lw.Etymologies[0].Parts[0].Attributes: expected length 50+, got %v`, len(part.Attributes))
	} else {
		if val, ok := part.Attributes["1|s|pres|act|ind"]; ok {
			expected := "mandō"
			if val != expected {
				t.Fatalf(`lw.Etymologies[0].Parts[0].Attributes["1|s|pres|act|ind"]: expected %q, got %q`, expected, part.Attributes["1|s|pres|act|ind"])
			}
		}
		if val, ok := part.Attributes["1|p|plup|act|sub|alt2"]; ok {
			expected := "mandāssēmus"
			if val != expected {
				t.Fatalf(`lw.Etymologies[0].Parts[0].Attributes["1|p|plup|act|sub|alt2"]: expected %q, got %q`, expected, part.Attributes["1|p|plup|act|sub|alt2"])
			}
		}
	}

}

func TestGetMeaning(t *testing.T) {
	meaning, err := GetMeaning("green", "en")
	if err != nil {
		t.Fatalf(`Error from GetMeaning: %q`, err)
	}
	expMeaning := "Having green as its color."
	if meaning != expMeaning {
		t.Fatalf(`GetMeaning: expected %q, got %q`,
			expMeaning, meaning)
	}

}
