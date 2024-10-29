package main

import (
	"github.com/gin-gonic/gin"
	"strings"
)

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

		// 如果公共前缀小于插入word的长度
		if i < len(word) {
			word = word[i:]
			c := word[0]
			for i := 0; i < len(rn.indices); i++ {
				// 如果与后继节点还有公共前缀，将rn指向子节点，递归流程
				if rn.indices[i] == c {
					rn = rn.children[i]
					continue walk
				}
			}

			// word剩余部分也没有公共前缀了
			// 构造新的节点进行插入
			rn.indices += string(c)
			child := &radixNode{}
			child.insertWord(word, fullWord)
			rn.children = append(rn.children, child)
			return
		}

		// 公共前缀恰好是rn.path,
		rn.end = true
		return
	}
}

// 传入相对路径和绝对路径，补充一个新生成的节点信息
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

// 两个单词的公共前缀
func commonPrefixLen(path1, path2 string) int {
	i := 0
	max := min(len(path1), len(path2))
	for i < max && path1[i] == path2[i] {
		i++
	}
	return i
}

// 查询流程
func (r *Radix) Search(word string) bool {
	node := r.root.search(word)
	return node != nil && node.fullPath == word && node.end
}

func (rn *radixNode) search(word string) *radixNode {

walk:
	for {
		prefix := rn.path
		// word 长于 path
		if len(word) > len(prefix) {
			// 没匹配上，直接返回nil
			if word[:len(prefix)] != prefix {
				return nil
			}
			// 扣除公共部分后的剩余部分
			word = word[len(prefix):]
			c := word[0]
			for i := 0; i < len(rn.indices); i++ {
				// 后继节点还有公共前缀，继续匹配
				if rn.indices[i] == c {
					rn = rn.children[i]
					continue walk
				}
			}
			return nil
		}

		// 精准匹配上了
		if word == prefix {
			return rn
		}

		// 走到这里意味着，len(word) <= len(prefix) && word != prefix
		return rn
	}

}

// 前缀匹配
func (r *Radix) StartWith(prefix string) bool {

	node := r.root.search(prefix)
	return node != nil && strings.HasPrefix(node.fullPath, prefix)
}

// 前缀统计
func (r *Radix) PassCnt(prefix string) int {

	node := r.root.search(prefix)
	if node != nil && strings.HasPrefix(node.fullPath, prefix) {
		return node.passCnt
	}

	return 0
}

// 删除流程
//func (r *Radix) Erase(word string) bool {
//
//	if !r.Search(word) {
//		return false
//	}
//
//	// root直接精准命中了
//	if r.root.fullPath == word {
//		// 如果一个孩子都没有
//		if len(r.root.indices) == 0 {
//			r.root.path = ""
//			r.root.fullPath = ""
//			r.root.end = false
//			r.root.passCnt = 0
//			return true
//		}
//
//		// 如果有一个孩子
//		if len(r.root.indices) == 1 {
//			r.root.children[0].path = r.root.path + r.root.children[0].path
//			r.root = r.root.children[0]
//			return true
//		}
//
//		// 如果有多个孩子
//		for i := 0; i < len(r.root.indices); i++ {
//			r.root.children[i].path = r.root.path + r.root.children[i].path
//		}
//
//		newRoot := radixNode{
//			indices:  r.root.indices,
//			children: r.root.children,
//			passCnt:  r.root.passCnt - 1,
//		}
//		r.root = &newRoot
//		return true
//	}
//	return false
//}

func main() {
	engine := gin.Default()
	engine.POST("/aa", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "hello world"})
	})
}
