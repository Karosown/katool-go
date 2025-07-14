package markdown

import (
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/karosown/katool-go/helper/jsonhp"
)

type Node struct {
	Level    int
	Title    string
	Content  string
	Children []*Node
}

func (n *Node) String() string {
	return jsonhp.ToJSONIndent(n)
}

func ToTree(md string) []*Node {
	var roots []*Node
	var stack []*Node
	lines := strings.Split(md, "\n")
	reHeader := regexp.MustCompile(`^(#{1,6})\s*(.+)`)
	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
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
			if curr.Content != "" {
				curr.Content += "\n"
			}
			curr.Content += line
		}
	}
	return roots
}

func ToHTML(md string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	return string(markdown.ToHTML([]byte(md), p, nil))
}
