package gee

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"runtime"
	"strings"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

// HandleFunc defines the request handler used by gee
// type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)

// ServeHTTP引擎的接口
type Engine struct {
	// 路由映射表
	// router map[string]HandlerFunc
	// router *router
	*RouterGroup
	router        *router
	groups        []*RouterGroup
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(message + "\nTrace back:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// 新建一个引擎对象
func New() *Engine {
	// return &Engine{router: make(map[string]HandlerFunc)}
	// return &Engine{router: newRouter()}
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	// key := method + "-" + pattern
	// engine.router[key] = handler
	// engine.router.addRoute(method, pattern, handler)
	pattern := group.prefix + comp
	log.Printf("Route %4s -%s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (engine *RouterGroup) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *RouterGroup) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	fmt.Println("use middlewares://", middlewares)
	group.middlewares = append(group.middlewares, middlewares...)
}

func (engine Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// c := newContext(w, req)
	// engine.router.handle(c)
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	fmt.Println("all-middlewares", len(middlewares))
	for k, v := range middlewares {
		fmt.Println(k, v)
	}
	c := newContext(w, req)
	// 读取当前路由匹配的中间件
	c.handlers = middlewares
	c.engine = &engine
	engine.router.handle(c)
}

// func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	key := req.Method + "-" + req.URL.Path
// 	if handler, ok := engine.router[key]; ok {
// 		handler(w, req)
// 	} else {
// 		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
// 	}
// }

// 静态资源处理器
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}
