package rhttp

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

type node struct {
	name         string
	childNodes   map[string]*node
	paramNode    *node
	wildcardNode *node
	view         *view
}

func newNode(name string) *node {
	return &node{
		name:       name,
		childNodes: make(map[string]*node),
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
func (r *Router) insertUrl(target_url string) *node {
	parts := getUrlParts(target_url)
	curr := r.rootNode
	for _, part := range parts {
		if part == "" {
			continue
		}
		if next_node, exists := curr.childNodes[part]; exists {
			curr = next_node
		} else {
			var node *node
			if part == "*" {
				if curr.wildcardNode == nil {
					node = newNode("*")
					curr.wildcardNode = node
				} else {
					node = curr.wildcardNode
				}
			} else if after, ok := strings.CutPrefix(part, ":"); ok {
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
func (r *Router) findTrailerNode(target_url string) (*node, Params) {
	parts := getUrlParts(target_url)
	params := NewParams()
	curr := r.rootNode
	for i, part := range parts {
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
			val, _ := url.QueryUnescape(part)
			pNode := curr.paramNode
			params[pNode.name] = val
			curr = pNode
			continue
		}
		// Check wildcard node - captures remaining path
		if curr.wildcardNode != nil {
			remainingParts := parts[i:]
			params["*"] = strings.Join(remainingParts, "/")
			return curr.wildcardNode, params
		}
		return nil, nil
	}
	return curr, params
}
