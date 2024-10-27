package main

import (
	"math/rand"
	"time"
)

/*
跳表：
以空间换时间，在有序链表的基础上，维持了多层级的索引结构，基于二分查找的方式
实现了Olog(n)级别时间复杂度的增删改查
跳表的性质：
1.跳表由多层索引结构组成
2.每层节点个数接近于相邻下层节点的一半
3.对于一个m层存在的节点，在1-m-1层，这个节点也一定存在
4.为了保证连贯性，跳表额外补充了本身不存储数据的头结点和尾结点
5.头结点和尾结点的高度是动态扩缩的，其高度取决于当前跳表内数据节点的最大高度
6.跳表的“跳”字体现在，在高层检索时，每跳过一个节点，实际上都略过了底层的大量数据，从而实现了检索加速
*/

type skipList struct {
	head *node
}

type node struct {
	next     []*node // 长度对应当前节点的高度
	key, val int
}

func (s *skipList) Get(key int) (int, bool) {
	// 根据key检索对应的node，如果node存在，则返回对应的val
	if _node := s.search(key); _node != nil {
		return _node.val, true
	}

	return -1, false
}

// 从跳表中检索key对应的node
func (s *skipList) search(key int) *node {
	// 每次检索从头部出发
	move := s.head
	// 每次检索从最高层出发，直到来到首层
	for level := len(s.head.next) - 1; level >= 0; level-- {
		// 在每一层中持续向右遍历，直到下一个节点不存在或者key大于等于key
		for s.head.next[level] != nil && s.head.next[level].key < key {
			move = move.next[level]
		}
		// 如果key相等，找到了，直接返回
		if move.next[level] != nil && move.next[level].key == key {
			return move.next[level]
		}
		// 当前层没找到，层数减一，继续向下
		// 没找到的原因是：move.next[level] == nil 或者
		// move.next[level] != nil && move.next[level].key > key
	}
	// 遍历完所有的层数,都没有找到目标，返回nil
	return nil
}

// 随机决定待插入节点所在的层数
func (s *skipList) roll() int {
	rand.Seed(time.Now().UnixNano())
	level := 0
	for rand.Intn(2) > 0 { // level是0的概率：1/2(第一次投出0) ；level1的概率： 1/2 * 1/2(第一次投出1，第二次投出0)
		level++
	}
	return level
}

// 将key-val对加入链表
func (s *skipList) put(key, val int) {
	// 加入key-val对已经存在，则直接对其更新
	if _node := s.search(key); _node != nil {
		_node.val = val
		return
	}

	// 计算出插入节点的高度
	newLevel := s.roll()

	// 新高度超过跳表的最大高度，需要对高度进行补齐
	for len(s.head.next)-1 < newLevel {
		s.head.next = append(s.head.next, nil)
	}

	// 创建出新的节点
	newNode := &node{
		next: make([]*node, newLevel+1),
		key:  key,
		val:  val,
	}

	// 从头结点最高层出发
	move := s.head
	for level := newLevel; level >= 0; level-- {
		// 向右遍历，直到右侧节点不存在或者key大于key
		for move.next[level] != nil && move.next[level].key < key {
			move = move.next[level]
		}
		// 调整指针关系，完成新节点的插入
		newNode.next[level] = move.next[level]
		move.next[level] = newNode
	}
}

// 删除流程
// 根据key删除对应的节点
func (s *skipList) del(key int) {
	// 如果key-val对不存在，无需删除，直接返回
	if _node := s.search(key); _node == nil {
		return
	}

	// 从头节点最高层出发
	move := s.head
	for level := len(s.head.next) - 1; level >= 0; level-- {
		// 向右遍历，直到右侧节点不存在或者key大于等于key
		for move.next[level] != nil && move.next[level].key < key {
			move = move.next[level]
		}

		// 右侧节点不存在或者key 大于 target, 直接跳过
		if move.next[level] == nil || move.next[level].key > key {
			continue
		}

		// 右侧节点的key必然等于key
		// 调整引用关系
		move.next[level] = move.next[level].next[level]
	}

	// 更细跳表的最大高度
	var dif int
	// 倘若某一层已经不存在数据节点，高度需要递减
	for level := len(s.head.next) - 1; level >= 0 && s.head.next[level] == nil; level++ {
		dif++
	}
	s.head.next = s.head.next[:len(s.head.next)-dif]
}

// range
func (s *skipList) Range(start, end int) [][2]int {
	// 通过ceiling方法，找到skipList中key大于等于start, 且最接近于start的节点ceilNode
	ceilNode := s.ceiling(start)
	if ceilNode == nil {
		return [][2]int{}
	}

	// 从ceilNode首层出发，向右遍历，把位于【start, node】区间的节点统统返回
	var res [][2]int
	for move := ceilNode; move != nil && move.key <= end; move = move.next[0] {
		res = append(res, [2]int{move.key, move.val})
	}

	return res
}

// 找到key >= target 且最接近于target的节点
func (s *skipList) ceiling(target int) *node {
	move := s.head
	for level := len(s.head.next) - 1; level >= 0; level-- {
		for move.next[level] != nil && move.next[level].key < target {
			move = move.next[level]
		}
		// 找到了直接返回
		if move.next[level] != nil && move.next[level].key == target {
			return move.next[level]
		}
	}
	// 没找到，就返回最接近的node
	return move.next[0]
}

// 天花板，大于等于target,且最接近于target的key-val键值对
func (s *skipList) Ceiling(target int) ([2]int, bool) {

	if ceilingNode := s.ceiling(target); ceilingNode != nil {
		return [2]int{ceilingNode.key, ceilingNode.val}, true
	}
	return [2]int{}, false
}

// 地板，小于等于target,且最接近于target的key-val键值对
func (s *skipList) floor(target int) *node {
	move := s.head
	for level := len(s.head.next) - 1; level >= 0; level-- {
		for move.next[level] != nil && move.next[level].key < target {
			move = move.next[level]
		}

		if move.next[level] != nil && move.next[level].key == target {
			return move.next[level]
		}

	}

	return move
}

func (s *skipList) Floor(target int) ([2]int, bool) {
	if floorNode := s.floor(target); floorNode != nil {
		return [2]int{floorNode.key, floorNode.val}, true
	}
	return [2]int{}, false
}

func main() {

}
