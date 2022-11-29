package web

import "net/http"

type HandleFunc func(ctx *Context) string

// web服务接口
type server interface {
	http.Handler

	// Start 启动server
	// addr服务启动地址 如 localhost:8081 监听本地8081端口
	Start(addr string) error

	// addRouter 注册路由
	addRouter(method, path string, handleFunc HandleFunc)
}

// HTTPServer web服务结构体
type HTTPServer struct {
	r *Router
}

// NewHTTPServer 创建新的http服务
func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		r: NewRouter(),
	}
}

// serve 处理http请求
func (s *HTTPServer) serve(ctx *Context) {
	node, ok := s.r.findRoute(ctx.req.Method, ctx.req.URL.Path)
	if ok == false || node == nil || node.handleFunc == nil {
		ctx.resp.WriteHeader(404)
		ctx.resp.Write([]byte("Not Found"))
		return
	}
	ctx.route = node.path
	node.handleFunc(ctx)
}

// ServeHTTP HTTPServer 处理请求的入口
func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		req:  request,
		resp: writer,
	}
	s.serve(ctx)
}

// Start 启动服务器
func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

// Post http网络协议的post请求方法路由注册
func (s *HTTPServer) Post(path string, handler HandleFunc) {
	s.r.addRouter(http.MethodPost, path, handler)
}

// Get http网络协议的get请求方法路由注册
func (s *HTTPServer) Get(path string, handler HandleFunc) {
	s.r.addRouter(http.MethodGet, path, handler)
}
