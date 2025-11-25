package router

import (
	"fmt"
	"path"
	"strings"

	"github.com/rnium/rhttp/internal/request"
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

func newParams() request.Params {
	return make(request.Params)
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
		if next_node, exists := curr.childNodes[part]; exists {
			curr = next_node
		} else {
			var node *Node
			if strings.HasPrefix(part, ":") {
				if curr.paramNode == nil {
					node = newNode(strings.Trim(part, ":"))
					curr.paramNode = node
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
	params := newParams()
	curr := r.rootNode
	for _, part := range parts {
		if curr.name == "products" && curr.paramNode != nil {
			fmt.Println(curr.paramNode.name, part)
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
