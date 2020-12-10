package util

import (
	"fmt"
	"testing"
)

func TestTrie(t *testing.T) {
	trie := NewTrie()
	trie.Add("黄色",nil)
	trie.Add("绿色",nil)
	trie.Add("蓝色",nil)

	// 当有两个敏感词有共同的前缀，长的敏感词不会被替换，替换短的
	trie.Add("你妈", nil)
	trie.Add("你妈炸了", nil)

	ret, str := trie.Check("我们这里有个蓝色的天空，后来变黄色了，真是透你妈你妈炸了, 昨日天气不错但你妈炸了，你的口头禅是草你妈", "***")
	fmt.Printf("ret:%#v, str:%v\n", ret, str)
	// ret:true, str:我们这里有个***的天空，后来变***了，真是透******炸了
}
