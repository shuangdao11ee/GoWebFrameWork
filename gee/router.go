package gee

import (
	"net/http"
	"strings"
)

//separate URL node by different methods of request
type router struct {
	roots map[string]*node
}

// roots key, eg, roots['GET'] roots['POST']

//return a  router pointer
func newRouter() *router {
	return &router{
		roots: make(map[string]*node),
	}
}

//Only one * is allowed
//separate url by /
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

//add a new router to node
func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	_, ok := r.roots[method] //GET OR POST OR PUT, roots is used to separate these node
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, handler, 0)
}

//get router by pattern, if router not exist, return nil, nil
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method] //get method root node , if not register , return nil

	if !ok {
		return nil, nil
	} //return nil if method and path no register

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern) //base on node.pattern to decide what params are
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(parts) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

//append handler into c.handlers if handler exists
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path) //if request method and path exist, return pattern of node and params
	if n != nil {
		c.Params = params
		c.handlers = append(c.handlers, n.handler) //insert handler after middleware
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
