package wiktionary

type TreeNode struct {
	Word         string       `json:"word"`
	LanguageCode string       `json:"lang-code"`
	Data         LanguageWord `json:"data"`
	Parent       *TreeNode    `json:"parent"`
}

func getEtymologyTree(word string, langCode string, languages []string) (TreeNode, error) {

	rootNode := TreeNode{}

	// first fetch the available translations for the specified word
	// this is likely to only include extant languages, so we'll need to pick up extinct ones later
	translations, err := GetTranslations(word, langCode, languages)
	if err != nil {
		return rootNode, err
	}

	// build a list of tree nodes, one per translation - we will build this into a structured tree later
	// as we build each one, fetch the LanguageWord data
	nodeList := make([]TreeNode, len(translations))

	for _, tr := range translations {

		// fetch the LanguageWord data
		var options WiktionaryOptions
		options.RequiredSections = Sec_All
		options.RequiredLanguages = AllLanguages
		lw, _ := processWord(word, langCode, options) // if this returns an error, we will just have a nil entry in the data

		node := TreeNode{
			Word:         tr.Word,
			LanguageCode: tr.Language,
			Data:         lw,
		}
		nodeList = append(nodeList, node)

	}

	// now scan the tree node list - for each one, read the etymology and get direct ancestors, one at a time
	// for _, node := range nodeList {
	// 	if node.Data != nil {
	// 		ancestors := getAncestors(node.Data)
	// 		for _, ancestor := range ancestors {

	// 		}
	// 	}
	// }
	return rootNode, err
}

func getAncestors(lw LanguageWord) []LinkedWord {
	var ancestors []LinkedWord

	// if there are no etymologies then there are no ancestors to return
	if len(lw.Etymologies) == 0 {
		return nil
	}

	// for now, use the default Etymology (the first one)
	for _, word := range lw.Etymologies[0].Words {
		if word.Relationship == "inherited" {
			ancestors = append(ancestors, word)
		}
	}

	return ancestors
}
