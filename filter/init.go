package filter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"qasystem/util"
)

var (
	trie *util.Trie
)

func Init(filename string) (err error){
	trie = util.NewTrie()
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("failed to open:", filename)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	for {
		word, _, e := reader.ReadLine() // 读取一行文件，并自动去掉'\n'
		if e == io.EOF {
			err = nil
			return
		}
		if e != nil {
			err = e
			break
		}
		err = trie.Add(string(word), nil)
		if err != nil {
			return
		}
	}
	return
}