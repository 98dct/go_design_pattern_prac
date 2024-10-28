package main

import "github.com/gin-gonic/gin"

/**
压缩前缀树/基数树：是一种更节省空间的trie，也是一种多叉树，传承了来自trie的基本设定
每个节点对应一段相对路径，从根节点到某个节点之间沿途所有路径拼接在一起形成的全路径即为key
压缩前缀树的合并：
倘若父节点有且仅有一个子节点，并且不存在单词以这个父节点作为结尾，radix tree会将父节点和子节点进行合并
*/

type radixNode struct {
	path     string       // 当前节点的相对路径
	fullPath string       // 全路径
	indices  string       // 孩子节点的path首字母
	children []*radixNode // 孩子节点
	end      bool         // 是否有路径以当前节点为终点
	passCnt  int          // 有多少路径通过当前节点
}

type Radix struct {
	root *radixNode
}

func NewRadix() *Radix {
	return &Radix{root: &radixNode{}}
}

func (r *Radix) Insert(word string) {
	if r.Search(word) {
		// 不重复插入
		return
	}
	r.root.insert(word)
}

func (rn *radixNode) insert(word string) {
	fullWord := word
	if rn.path == "" && len(rn.children) == 0 {
		rn.insertWord(word, word)
		return
	}
walk:
	for {
		i := commonPrefixLen(word, rn.path)
		if i > 0 {
			rn.passCnt++
		}

		// 如果公共前缀小于当前节点的path
		// split node
		if i < len(rn.path) {
			child := radixNode{
				path:     rn.path[i:],
				fullPath: rn.fullPath,
				// 当前节点的后继节点进行委托
				indices:  rn.indices,
				children: rn.children,
				end:      rn.end,
				passCnt:  rn.passCnt - 1,
			}

			rn.indices = string(rn.path[i])
			rn.children = []*radixNode{&child}
			rn.fullPath = rn.fullPath[:len(rn.fullPath)-len(rn.path)+i]
			rn.path = rn.path[:i]
			rn.end = false
		}

	}
}

func (rn *radixNode) insertWord(path, fullPath string) {
	rn.path, rn.fullPath = path, fullPath
	rn.passCnt = 1
	rn.end = true
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func commonPrefixLen(path1, path2 string) int {
	i := 0
	max := min(len(path1), len(path2))
	for i < max && path1[i] == path2[i] {
		i++
	}
	return i
}

func main() {
	engine := gin.Default()
	engine.POST("/aa", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "hello world"})
	})
}
