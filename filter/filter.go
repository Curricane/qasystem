package filter

func Replace(text, replace string)(ret string, isReplaced bool) {
	isReplaced, ret = trie.Check(text, replace)
	return 
}
