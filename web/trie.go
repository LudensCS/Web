package web

import "strings"

// Trie树
type node struct {
	pattern  string
	part     string
	children []*node
	iswild   bool
}

// 寻找一个匹配节点
func (nd *node) matchchild(part string) *node {
	for _, child := range nd.children {
		if child.iswild || child.part == part {
			return child
		}
	}
	return nil
}

// 寻找所有匹配节点
func (nd *node) matchchildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range nd.children {
		if child.iswild || child.part == part {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 插入路由
func (nd *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		nd.pattern = pattern
		return
	}
	part := parts[height]
	child := nd.matchchild(part)
	if child == nil {
		child = &node{
			part:   part,
			iswild: part[0] == ':' || part[0] == '*',
		}
		nd.children = append(nd.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 匹配路由
func (nd *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(nd.part, "*") {
		if nd.pattern == "" {
			return nil
		}
		return nd
	}
	part := parts[height]
	children := nd.matchchildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
