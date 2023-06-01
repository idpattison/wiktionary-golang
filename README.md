# Golang library to parse Wiktionary data

The aim of [this library](https://github.com/ianpattison-google/wiktionary-golang) is to convert Wiktionary data from its raw Wikitext form to JSON.

## Usage

~~~
GetWord(word string, langCode string) (LanguageWord, error)
GetWordWithOptions(word string, langCode string, options WiktionaryOptions) (LanguageWord, error)
~~~
- word = the required word to be parsed
- langCode = the code for the language of the word (see languages.go)
- options = options for more control (see core.go)

## Example
~~~
import (
	"github.com/ianpattison-google/wiktionary-golang"
)

func main() {
		lw, err := wiktionary.GetWord("red", "en)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
~~~
lw is a structure containing the parsed word, and a JSON version will be written to the working directory if required.

The JSON output looks like this (example is the English word 'red') - this will vary dependent on the options chosen - NB a large number of cognates and translations have been omitted for brevity.
~~~
{
    "word": "red",
    "meaning": "Having red as its color.",
    "lang": "English",
    "lang-code": "en",
    "pron": [
        "enPR: rĕd, IPA(key): /ɹɛd/, [ɻʷɛˑd̥]",
        "Homophone: read (past tense/participle)",
        "Rhymes: -ɛd"
    ],
    "ipa": "/ɹɛd/",
    "etym": [
        {
            "name": "Etymology 1",
            "text": "From Middle English red, from Old English rēad, from Proto-West Germanic *raud, from Proto-Germanic *raudaz (compare West Frisian read, Low German root, rod, Dutch rood, German rot, Danish and Norwegian Bokmål rød, Norwegian Nynorsk raud), from Proto-Indo-European *h₁rowdʰós, from the root *h₁rewdʰ- (compare Old Armenian արոյդ/արուրդ (aroyd/arurd), Welsh rhudd, Latin ruber, rufus, Tocharian A rtär, Tocharian B ratre, Ancient Greek ἐρυθρός (eruthrós), Albanian pruth (“redhead”), Old Church Slavonic рудъ (rudŭ), Czech rudý, Lithuanian raúdas, Avestan 𐬭𐬀𐬊𐬌𐬛𐬌𐬙𐬀  (raoidita), Sanskrit रुधिर (rudhirá, “red, bloody”)).\n",
            "words": [
                {
                    "type": "root",
                    "lang": "ine-pro",
                    "word": "*h₁rewdʰ-"
                },
                {
                    "type": "inherited",
                    "lang": "enm",
                    "word": "red"
                },
                {
                    "type": "inherited",
                    "lang": "ang",
                    "word": "rēad"
                },
                {
                    "type": "inherited",
                    "lang": "gmw-pro",
                    "word": "*raud"
                },
                {
                    "type": "inherited",
                    "lang": "gem-pro",
                    "word": "*raudaz"
                },
                {
                    "type": "cognate",
                    "lang": "fy",
                    "word": "read"
                },
                {
                    "type": "cognate",
                    "lang": "grc",
                    "word": "ἐρυθρός",
                    "translit": "eruthrós"
                },
                {
                    "type": "cognate",
                    "lang": "sq",
                    "word": "pruth",
                    "meaning": "redhead"
                },
                {
                    "type": "cognate",
                    "lang": "sa",
                    "word": "रुधिर",
                    "meaning": "red, bloody",
                    "translit": "rudhirá"
                }
            ],
            "parts": [
                {
                    "name": "Adjective",
                    "head": "red (comparative redder or more red, superlative reddest or most red)",
                    "attrs": {
                        "comparative": "redder or more red",
                        "superlative": "reddest or most red"
                    },
                    "meanings": [
                        "Having red as its color.",
                        "(of hair) Having an orange-brown or orange-blond colour; ginger.",
                        "(card games, of a card) Of the hearts or diamonds suits. Compare black (“of the spades or clubs suits”)",
                        "(often capitalized) Supportive of, related to, or dominated by a political party or movement represented by the color red:",
                        "(chiefly derogatory, offensive) Amerind; relating to Amerindians or First Nations",
                        "(astronomy) Of the lower-frequency region of the (typically visible) part of the electromagnetic spectrum which is relevant in the specific observation.",
                        "(particle physics) Having a color charge of red."
                    ]
                },
                {
                    "name": "Noun",
                    "head": "red (countable and uncountable, plural reds)",
                    "attrs": {
                        "count": "countable and uncountable",
                        "plural": "reds"
                    },
                    "meanings": [
                        "(countable and uncountable) Any of a range of colours having the longest wavelengths, 670 nm, of the visible spectrum; a primary additive colour for transmitted light: the colour obtained by subtracting green and blue from white light using magenta and yellow filters; the colour of blood, ripe strawberries, etc.",
                        "(countable) A revolutionary socialist or (most commonly) a Communist; (usually capitalized) a Bolshevik, a supporter of the Bolsheviks in the Russian Civil War.",
                        "(countable, snooker) One of the 15 red balls used in snooker, distinguished from the colours.",
                        "(countable and uncountable) Red wine.",
                        "(countable, informal, Britain, birdwatching) A redshank.",
                        "(derogatory, offensive) An Amerind.",
                        "(slang) The drug secobarbital; a capsule of this drug.",
                        "(informal) A red light (a traffic signal)",
                        "(Ireland, Britain, beverages, informal) red lemonade",
                        "(particle physics) One of the three color charges for quarks.",
                        "(US, colloquial, uncountable) chili con carne (usually in the phrase \\\"bowl of red\\\")",
                        "(informal) The redfish or red drum, Sciaenops ocellatus, a fish with reddish fins and scales."
                    ],
                    "trans": [
                        {
                            "lang": "af",
                            "word": "rooi"
                        },
                        {
                            "lang": "ain",
                            "word": "フレ",
                            "translit": "hure"
                        },
                        {
                            "lang": "akl",
                            "word": "puea"
                        },
                        {
                            "lang": "sq",
                            "word": "kuq"
                        },
                        {
                            "lang": "am",
                            "word": "ቀይ"
                        },
                    ]
                }
            ]
        },
    ]
}
~~~


## core.go

- Define constants which show the sections we are interested in - you can choose to include or omit sections such as anagrams, synonyms etc

- Define constants which represent the language codes

- GetWord - given a word and language, return a LanguageWord object (an in-memory representation of our target JSON structure) with all sections, and etymology for all languages

  - This function is mainly a wrapper for the internal function processWord

- GetWordWithOptions - as GetWord, but giving more control over the required sections and languages

- GetMeaning - call processWord but only get the Meanings section, and return the meaning

  - NB we want to change this to get meanings from Wordnet - this will help us to cluster similar words for etymological analysis
  - The challenge here is around homonyms - for example if we find a word ‘bear’ - does it mean the animal, or to carry something?

- GetTranslations - call processWord but only get the translations for the languages specified

- GetIpaPronunciation - call processWord but only get the IPA translation

- GetEtymologyTree - wrapper for the internal function getEtymologyTree

- GetLanguageFromCode - wrapper for the internal function getLanguageFromCode

- processWord - the main controlling function

  - Get the JSON content from Wiktionary (getWordDataFromWiktionary)
  - Extract the Wikitext (getWikitext)
  - Process the Wikitext into sections (processWikitext)
  - Get the relevant sections for the specified language (extractLanguageSections)
  - Parse the language sections and build a LanguageWord structure (parseSections)
  - For debug purposes we also write a JSON file (writeJson) and a Wikitext file


## raw-data.go

- Define a Section struct as a header plus an array of lines

- getWordDataFromWiktionary - construct a URL and call the Wiktionary API, returning the body of text

- getWikitext - use regex to find the Wikitext block in the JSON - call convertWikitext when done

- convertWikitext

  - Decode any surrogate pairs - these would be found in languages such as Avestan - and replace with an equivalent rune

    - NB this is clunky, we need to see if there is a better way to do this

  - Decode backslash encoded characters

    - NB can we parse the HTML text which matches the WIkitext to get this?

- processWikitext

  - Split the Wikitext into a slice of strings

  - Group the lines into sections by finding a section header and then appending each line until we find another section - we can ignore

    - The end of language marker (four dashes) as this is implied by finding a new language header
    - HTML comments
    - Categories (we don’t process them)
    - Blank lines

- extractLanguageSections - iterate through all the sections and return only the ones from the language header to the next language header

  - NB can we do this with indexes rather than copying line by line?

- getLanguageFromCode - return the language name

- getPageTitle - work out the page title, which will be different if we’re dealing with a reconstructed word

- getConvertedTextFromWiktionary

  - Make a call to Wiktionary’s API with an isolated tag - this is usually to get a human-readable version of Wikitext for etymology purposes
  - API to be called (this example is for the text in the page _red_) https&#x3A;//en.wiktionary.org/w/api.php?action=**parse**&text=**{{cog|nds|root}}**&prop=**text**&title=**red**&formatversion=**2**&format=**json**
  - Select the part from text”:” up to \\n&lt;/p>
  - Remove everything in HTML braces &lt;>
  - Replace explicit spaces (&amp;#32, &amp;nbsp, etc) with actual spaces
  - NB we changed to using the HTML5-compliant parser in _golang.org/x/net/html_


## process-wikitext.go

- parseSections

  - Create the LanguageWord (LW) object which all subsequent functions will write into

  - Iterate over the sections and call parseSection

  - Assign a meaning based on the etymology

    - NB can we do this based on Wordnet?

- parseSection - based on the section type, call the relevant section parsing function

- parsePronunciationSection

  - Process each line beginning with \* and add it to the LW
  - Audio lines need some special handling to add a link to the audio file
  - Get the first IPA tag and also store that separately
  - NB we usually attach this to the LW at the word level, but if there are homographs (spelled the same, pronounced differently) then we add it to the etymology level

- parseEtymologySection

  - Process each line (which represents a paragraph of text)
  - Get the text so we have a human-readable representation
  - Call parseLinkedWord to get the individual linked words

- parseLinkedWord

  - Process each tag in the line - ignore **_m_** tags for now
  - Process root word, inherited, cognate and descendant words
  - Handle non-standard inheritance - borrowings, calques and semantic loans
  - Handle any provided transliteration
  - If we have non-Latin script but no transliteration, get one from the etymology text

- parseDescendantSection

  - Process each line and call parseLinkedWord to get the individual linked words

- parsePartofSpeechSection

  - Read the headword line and get the tags
  - Read the meaning lines, but ignore quotations for now
  - Nouns, verbs, adjectives and adverbs call a specific handling function at this point
  - For simpler words (usually in non-English languages) we may need to create an etymology if one doesn’t exist in the WIkitext

- parseNoun

  - Check for noun gender - usually parameter 1, or explicitly defined
  - Get feminine and masculine forms if they exist
  - Countability will be param 1, or param 2 for gendered languages
  - Get the headword forms (plural, genitive, diminutive etc)

- parseAdjective

  - Get the headword forms (plural, masc/fem, comparative, superlative, etc)

- parseVerb

  - Get the headword forms (there are potentially many - simple past, particples, past tense, imperative, infinitive, etc)

- getHeadwordForm

  - Get the headword from the text
  - NB many of the headword forms are computed by Wiktionary based on minimal Wikitext input - so we need to parse the resultant text to find them

- getHeadwordItem

  - Get the part of the headword text in braces and split by commas

- parseTranslationSection

  - We only translate the first block of translations to get the core meaning of the word
  - For each line, read the language and word
  - Add these to the current etymology

- parseOtherSections

  - For other sections such as synonyms, antonym and anagrams we just get a list of the items

- Helper functions

  - getAllTags - returns all WIkitext tags in the given text
  - searchForTag - find a specific tag in the given text
  - splitTag - given a complete tag, splits it into a map of position and parameters
  - sectionRequired - returns true if the given section is specified in the options
  - languageRequired - returns true if the given language is specified in the options


## etym-tree.go 

- This is experimental, we are trying to build a tree of words across languages and history


## json-output.go 

- Defines structs used in the JSON output
- writeJson - writes the JSON to an output file


## languages.go 

- Defines a map of language codes and the equivalent language name
- Defines helper variables for specified language subsets
