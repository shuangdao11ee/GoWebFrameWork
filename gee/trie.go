package gee

import (
	"strings"
)

//node is used to build a tree like data struct
//and at the end of the tree replace different URL that have been registered
type node struct {
	//why there are chinese words below? i cant type chinese on this god damn IDE
	pattern  string      // 待匹配路由， 例如/p/:lang
	part     string      // 路由中的一部分， 例如:lang
	children []*node     // 子节点， 例如[doc, tutorial, intro]
	isWild   bool        //是否hu匹配，part 含有 ： 或 *时为true
	handler  HandlerFunc //if pattern exist, handler exist
}

//第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//crate node tree by the pattern input
func (n *node) insert(pattern string, parts []string, handler HandlerFunc, height int) {
	if len(parts) == height {
		n.pattern = pattern
		n.handler = handler
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, handler, height+1)
}

//search node by pattern input
//and pattern has already been handled into a []string struct
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part) //if cant match, return nil

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
