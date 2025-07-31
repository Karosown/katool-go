package markdown

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/helper/jsonhp"
)

type Node struct {
	Level    int
	Title    string
	Content  string
	Children Tree
}

func (n *Node) String() string {
	return jsonhp.ToJSONIndent(n)
}

func (n *Node) ToHtml() string {
	md := n.ToMarkDown()
	return ToHtml(md)
}
func (n *Node) ToMarkDown() string {
	var sb strings.Builder
	n.toMarkdown(&sb, 0)
	return sb.String()
}
func (n *Node) toMarkdown(sb *strings.Builder, indent int) {
	// 标题行
	header := strings.Repeat("#", n.Level)
	sb.WriteString(fmt.Sprintf("%s %s\n", header, n.Title))
	// 内容
	if n.Content != "" {
		sb.WriteString(n.Content)
		// 避免内容连在下一个标题上，补个换行
		if !strings.HasSuffix(n.Content, "\n") {
			sb.WriteString("\n")
		}
	}
	// 子节点
	for _, c := range n.Children {
		c.toMarkdown(sb, indent+1)
	}
}
func ToTree(md string) Tree {
	var roots Tree
	var stack Tree
	lines := strings.Split(md, "\n")
	reHeader := regexp.MustCompile(`^(#{1,6})\s*(.+)`)
	for _, line := range lines {
		if matches := reHeader.FindStringSubmatch(line); matches != nil {
			level := len(matches[1])
			title := matches[2]
			node := &Node{
				Level:   level,
				Title:   title,
				Content: "",
			}
			for len(stack) > 0 && stack[len(stack)-1].Level >= level {
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				roots = append(roots, node)
			} else {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, node)
			}
			stack = append(stack, node)
		} else if len(stack) > 0 {
			curr := stack[len(stack)-1]
			// 与空行无关，每一行都加，空行就多一个\n
			if curr.Content != "" {
				curr.Content += "\n"
			}
			curr.Content += line
		}
	}
	return roots
}
func ToHtml(md string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	return string(markdown.ToHTML([]byte(md), p, nil))
}

type Tree []*Node

// ToMarkdown 递归将 Node 树转回 Markdown 字符串
func (t *Tree) ToMarkdown() string {
	var sb strings.Builder
	for _, n := range *t {
		n.toMarkdown(&sb, 0)
	}
	return sb.String()
}
func (t *Tree) ToHtml() string {
	return ToHtml(t.ToMarkdown())
}
func (t *Tree) ToStream() *stream.Stream[*Node, Tree] {
	return stream.Of(t)
}
