package main

import (
	"log"
	"os"

	"github.com/ianpattison-google/wiktionary-golang"
)

func main() {

	// get the required word - if no language given then we assume English
	var word string
	langCode := "en"
	etymTree := false
	if len(os.Args) < 2 {
		word = "red" // if no args (arg 0 is the program name) then assume word = red for testing
	} else {
		word = os.Args[1]
	}
	if len(os.Args) > 2 {
		langCode = os.Args[2] // capture the language if specified
	}
	if len(os.Args) > 3 && os.Args[3] == "-e" {
		etymTree = true // the -e flag will produce an etymology tree for the given word
	}

	if etymTree {

	} else {

		_, err := wiktionary.GetWord(word, langCode)
		if err != nil {
			log.Fatalln(err)
		}

	}

}
