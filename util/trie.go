package util

type Node struct {
	char rune
	Data interface{} // 存储附带的数据，没有则为nil
	parent int
	Depth int
	childs map[rune]*Node //当前节点的所有子节点
	term bool // 一个敏感词的结尾
}

type Trie struct {
	root *Node
	size int
}

func NewNode() *Node {
	return &Node{
		childs: make(map[rune]*Node, 32),
	}
}

func NewTrie() *Trie {
	return &Trie{
		root: NewNode(),
	}
}

func (p *Trie) Add(key string, data interface{})(err error) {
	node := p.root
	runes := []rune(key) // 或得key所有的汉字
	for _, r := range runes { // 遍历所有的汉字
		ret ,ok := node.childs[r] // 当前汉字是否在trie当前节点的子节点map中
		if !ok { // 不在，则插入到当前节点的子节点map中
			ret = NewNode()
			ret.Depth = node.Depth + 1
			ret.char = r
			node.childs[r] = ret
		}

		// 在，则切换到该子节点，继续往下遍历汉字
		node = ret
	}

	// 当前敏感词结束了，即遍历到这个term节点，
	//不管后面有没有子节点，就已经检索到了一个敏感词
	node.term = true
	node.Data = data
	return
}

// 找到key汉字中，最后一个汉字所在的节点
func (p *Trie) findNode(key string) (result *Node) {
	node := p.root
	runes := []rune(key)
	for _, v := range runes {
		ret, ok := node.childs[v]
		if !ok { // 该路径深度无敏感词
			return
		}

		node = ret // 在，则切换到该子节点，继续往下遍历汉字
	}

	result = node
	return
}

// 查找当前节点相同前缀的子树（所有子节点），感觉有问题
func (p *Trie) collectNode(node *Node) (result []*Node) {
	if node == nil {
		return
	}
	if node.term {
		result = append(result, node)
		return
	}

	var queue []*Node
	queue = append(queue, node)
	for i := 0; i < len(queue); i++ {
		if queue[i].term {
			result = append(result, queue[i])
			continue
		}
		for _, v1 := range queue[i].childs {
			queue = append(queue, v1)
		}
	}
	return
}

func (p *Trie) PrefixSearch(key string) (result []*Node) {
	node := p.findNode(key)
	if node == nil {
		return
	}
	result = p.collectNode(node)
	return
}

/*
Check 检测敏感词，并替换敏感词
param: text 检测的文本 replace替换敏感词的文字
*/
func (p *Trie) Check(text, replace string) (result bool, str string) {
	runes := []rune(text)
	if p.root == nil {
		return
	}

	var left []rune
	node := p.root
	start := 0
	for index, v := range runes {
		ret, ok := node.childs[v]

		// 当前text中的字不在node子节点中，返回root，遍历下一个字
		if !ok {
			left = append(left, runes[start: index+1]...) // 保存没有敏感词的部分
			start = index + 1
			node = p.root
			continue
		}

		node = ret
		// 有敏感词
		if ret.term {
			result = true

			// 从头开始遍历下一个字
			node = p.root
			left = append(left, ([]rune(replace))...)
			start = index + 1
			continue
		}

	}
	str = string(left)
	return
}