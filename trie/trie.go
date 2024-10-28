package main

import "fmt"

/*
trie: 前缀树，是一种多叉树，每个节点存储一个字符，数据不单独存放于每个节点中，而是由根节点root出发，
直至来到目标节点的路径组成
优点: 可以节省存储空间、前缀频率统计和前缀模糊匹配
*/

// 前缀树的节点
type trieNode struct {
	next    [26]*trieNode // 默认只能存储26个字母组成的单词
	end     bool          // 是否存在单词以当前节点为结尾
	passCnt int           // 以“从根节点到当前节点”形成的path为前缀的数量
}

// 前缀树
type Trie struct {
	root *trieNode
}

func NewTrie() *Trie {
	return &Trie{root: &trieNode{}}
}

// 是否存在某个单词
func (t *Trie) Search(word string) bool {
	node := t.search(word)
	return node != nil && node.end
}

func (t *Trie) search(target string) *trieNode {

	move := t.root
	// 依次遍历target中的每个字符
	for _, ch := range target {
		if move.next[ch-'a'] == nil {
			return nil
		}
		move = move.next[ch-'a']
	}
	// 此时move是target最后一个字母所在的trieNode
	// 目标单词不一定存在于这个树，例如：之前插入了apple，我们查找app，可以找到，还要在判断end是否为true
	return move
}

// 是否存在包含指定前缀的单词
func (t *Trie) StartWith(prefix string) bool {
	return t.search(prefix) != nil
}

// 包含指定前缀的单词(从根节点出发到当前节点的路径形成的单词)出现的次数
func (t *Trie) PassCnt(prefix string) int {
	node := t.search(prefix)
	if node == nil {
		return 0
	}
	return node.passCnt
}

// 插入流程
func (t *Trie) insert(word string) {
	if t.Search(word) {
		return
	}

	move := t.root

	for _, ch := range word {
		if move.next[ch-'a'] == nil {
			move.next[ch-'a'] = &trieNode{}
		}
		move.next[ch-'a'].passCnt++
		move = move.next[ch-'a']
	}
	move.end = true
}

// 删除流程
func (t *Trie) erase(word string) bool {
	if !t.Search(word) {
		return false
	}

	move := t.root
	for _, ch := range word {
		move.next[ch-'a'].passCnt--
		if move.next[ch-'a'].passCnt == 0 {
			move.next[ch-'a'] = nil
			return true
		}
		move = move.next[ch-'a']
	}
	move.end = false
	return true
}

func main() {

	trie := NewTrie()
	trie.insert("hello")
	trie.insert("hell")
	trie.insert("hel")
	trie.insert("world")
	trie.insert("worl")
	trie.insert("wor")

	fmt.Println(trie.Search("wo"))     // false
	fmt.Println(trie.Search("wor"))    // true
	fmt.Println(trie.Search("worl"))   // true
	fmt.Println(trie.StartWith("wo"))  // true
	fmt.Println(trie.PassCnt("wo"))    // 3
	fmt.Println(trie.PassCnt("wor"))   // 3
	fmt.Println(trie.PassCnt("worl"))  // 2
	fmt.Println(trie.PassCnt("world")) // 1
	fmt.Println(trie.erase("world"))   // true
	fmt.Println(trie.Search("world"))  // false
	fmt.Println(trie.Search("worl"))   // true
}
