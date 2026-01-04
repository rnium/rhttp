package router

import (
	"fmt"
	"path"
	"strings"

	"github.com/rnium/rhttp/internal/http/request"
)

type Node struct {
	name       string
	childNodes map[string]*Node
	paramNode  *Node
	view       *View
}

func newNode(name string) *Node {
	return &Node{
		name:       name,
		childNodes: make(map[string]*Node),
	}
}

func getUrlParts(url string) []string {
	if url == "" {
		return []string{}
	}
	url = path.Clean(url)
	url_trimmed := strings.Trim(url, "/")
	return strings.Split(url_trimmed, "/")
}

// Splits the parts of an URL separated by forward slash (/) and inserts into the Trie as node.
// Returns the trailing node of the url
func (r *Router) insertUrl(target_url string) *Node {
	parts := getUrlParts(target_url)
	curr := r.rootNode
	for _, part := range parts {
		if part == "" {
			continue
		}
		if next_node, exists := curr.childNodes[part]; exists {
			curr = next_node
		} else {
			var node *Node
			if after, ok := strings.CutPrefix(part, ":"); ok {
				if curr.paramNode == nil {
					node = newNode(after)
					curr.paramNode = node
				} else {
					node = curr.paramNode
					if node.name != after {
						panic(fmt.Sprintf("conflict with parameterized node name. %s != %s", node.name, after))
					}
				}
			} else {
				node = newNode(part)
				curr.childNodes[part] = node
			}
			curr = node
		}
	}
	return curr
}

// Find for the trailer node of an url
func (r *Router) findTrailerNode(target_url string) (*Node, request.Params) {
	parts := getUrlParts(target_url)
	params := request.NewParams()
	curr := r.rootNode
	for _, part := range parts {
		if part == "" {
			continue
		}
		// First check the static parts
		if node, exists := curr.childNodes[part]; exists {
			curr = node
			continue
		}
		// Check parameterized node
		if curr.paramNode != nil {
			pNode := curr.paramNode
			params[pNode.name] = part
			curr = pNode
			continue
		}
		return nil, nil
	}
	return curr, params
}
