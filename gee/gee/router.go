package gee

import (
	"fmt"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	// return &router{handlers: make(map[string]HandlerFunc)}
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

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

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	// log.Printf("Route %4s-%s", method, pattern)
	// key := method + "-" + pattern
	// r.handlers[key] = handler

	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	fmt.Println("searchParts:////", searchParts)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	// key := c.Method + "-" + c.Path
	// if handler, ok := r.handlers[key]; ok {
	// 	handler(c)
	// } else {
	// 	c.String(http.StatusNotFound, "404 NOT FOUND:%s\n", c.Path)
	// }

	// n, params := r.getRoute(c.Method, c.Path)
	// fmt.Println("router.go handle:///", n, params)
	// if n != nil {
	// 	c.Params = params
	// 	key := c.Method + "-" + n.pattern
	// 	r.handlers[key](c)
	// } else {
	// 	c.String(http.StatusNotFound, "404 NOT FOUND:%s\n", c.Path)
	// }

	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = params
		fmt.Println("路由上的所有handlers:///", r.handlers)
		for k, v := range r.handlers {
			fmt.Println(k, v)
		}
		fmt.Println("所有handlers:///", key, c.handlers, r.handlers[key])
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	fmt.Println(len(c.handlers))
	// for _, handler := range c.handlers {
	// 	fmt.Println(handler)
	// 	if handler != nil {
	// 		handler(c)
	// 	}
	// }
	c.Next()
}
