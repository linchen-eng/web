package web

import (
	"net/http"
)

type HandleFunc func(ctx *Context)

// Server web服务接口
// 先注册路由 再启动服务
type Server interface {
	http.Handler

	// Start 启动server
	// addr服务启动地址 如 localhost:8081 监听本地8081端口
	Start(addr string) error

	// addRouter 注册路由
	addRouter(method, path string, handleFunc HandleFunc)
}

// 确保一定实现了 Server 接口
var _ Server = &HTTPServer{}

// HTTPServer web服务结构体
type HTTPServer struct {
	router
	mdls []Middleware
}

// NewHTTPServer 创建新的http服务
func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: NewRouter(),
	}
}

// serve 处理http请求
func (s *HTTPServer) serve(ctx *Context) {
	mi, ok := s.findRoute(ctx.req.Method, ctx.req.URL.Path)
	node := mi.n
	if ok == false || node == nil || node.handleFunc == nil {
		ctx.resp.WriteHeader(404)
		data := []byte("Not Found")
		n, err := ctx.resp.Write(data)
		if err != nil || len(data) != n {
			return
		}
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
	// 最后一个是这个
	root := s.serve
	// 然后这里就是利用最后一个不断往前回溯组装链条
	// 从后往前
	// 把后一个作为前一个的 next 构造好链条
	for i := len(s.mdls) - 1; i >= 0; i-- {
		root = s.mdls[i](root)
	}

	// 这里执行的时候，就是从前往后了
	// 这里，最后一个步骤，就是把 RespData 和 RespStatusCode 刷新到响应里面
	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			// 就设置好了 RespData 和 RespStatusCode
			next(ctx)
			//s.flashResp(ctx)
		}
	}
	root = m(root)
	root(ctx)
}

// Start 启动服务器
func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

// Post http网络协议的post请求方法路由注册
func (s *HTTPServer) Post(path string, handler HandleFunc) {
	s.addRouter(http.MethodPost, path, handler)
}

// Get http网络协议的get请求方法路由注册
func (s *HTTPServer) Get(path string, handler HandleFunc) {
	s.addRouter(http.MethodGet, path, handler)
}
