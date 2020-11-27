package gee

import (
	"log"
	"net/http"
	"strings"
)

//HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

//RouterGroup separate URL into different parts, and different parts
//can used middlewares independently
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

//Engine implement the interface of  ServerHTTP
//father of all function
//absolut core of the framework
type Engine struct {
	*RouterGroup
	router      *router
	groups      []*RouterGroup // store all groups
	db          *DB
	accesstoken *AccessTokenJson
}

//Use is defined to add middleware to group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

//New is the construceor of gee.Engine
func New(db *DB, accesstoken *AccessTokenJson) *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	engine.db = db
	engine.accesstoken = accesstoken
	return engine
}

//Group is defined to create a new RouterGroup
//remember all groups share the same Engine instance
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

//jus like its name, add a new router
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Roter %4s - %s", method, pattern)
	group.engine.router.addRouter(method, pattern, handler)
}

//GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

//POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

//Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//to satisfy the requirement of http.ListenAndServe function
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...) //put middlewares that satisfied the requirement into context.handlers
		}
	}
	c := newContext(w, req, engine.db, engine.accesstoken)
	c.handlers = middlewares
	engine.router.handle(c)
}

////create static handler
//func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
//	absolutePath := path.Join(group.prefix, relativePath)
//	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
//	return func(c *Context) {
//		file := c.Param("filepath")
//		//Check if file exists and/or if we have permission to access it
//		if _, err := fs.Open(file); err != nil {
//			c.Status(http.StatusNotFound)
//			return
//		}
//
//		fileServer.ServeHTTP(c.Writer, c.Req)
//	}
//}
//
////serve static files
//func (group *RouterGroup) Static(relativePath string, root string) {
//	handler := group.createStaticHandler(relativePath, http.Dir(root))
//	urlPattern := path.Join(relativePath, "/*filepath")
//	//Register GET handlers
//	group.GET(urlPattern, handler)
//}
