package web

import (
	"strings"
)

// 路由树的节点结构
type node struct {
	//当前路由
	route string
	//当前节点的路径
	path string
	//当前节点的下一级子节点
	child map[string]*node
	//当前节点处理的回调方法
	handleFunc HandleFunc
	//通配符路由
	wildcard map[string]*node
}

// Router 路由树的结构
type Router struct {
	trees map[string]*node
}

// NewRouter 初始化路由
func NewRouter() *Router {
	return &Router{trees: make(map[string]*node)}
}

// wildcardOrCreate 为当前节点创建通配符节点
func (n *node) wildcardOrCreate(path string) *node {
	if n.wildcard == nil {
		n.wildcard = make(map[string]*node)
	}

	if _, ok := n.wildcard[path]; !ok {
		//当前路径不存在子节点需要创建子节点组
		n.wildcard[path] = &node{path: path, child: nil, handleFunc: nil, wildcard: nil}
	}
	return n.wildcard[path]
}

// childOrCreate 为当前节点创建子节点
func (n *node) childOrCreate(path string) *node {
	if n.child == nil {
		n.child = make(map[string]*node)
	}

	if _, ok := n.child[path]; !ok {
		//当前路径不存在子节点需要创建子节点组
		n.child[path] = &node{path: path, child: nil, handleFunc: nil}
	}
	return n.child[path]
}

// 注册路由到路有树
func (r *Router) addRouter(method, path string, handleFunc HandleFunc) {
	if _, ok := r.trees[method]; !ok {
		//当前http方法未注册，需要初始化路由树
		r.trees[method] = &node{path: "/", child: nil, handleFunc: nil}
	}

	//路径切割 "/"
	subStrs := strings.Split(path, "/")

	//获取当前路由的根节点
	root := r.trees[method]
	for index, subStr := range subStrs {
		if len(subStr) == 0 && (index == 0 || index == len(subStrs)-1) {
			// "/" 将导致第一个元素未空字符 "root/" 将导致最后一个元素为空字符串
			continue
		}
		if len(subStr) == 0 {
			panic("路由格式不支持:路径中不能包含'//'这样的字符串")
		}
		//获取匹配的子节点 or 创建子节点
		switch subStr {
		case "*":
			root = root.wildcardOrCreate(subStr)
			break
		default:
			root = root.childOrCreate(subStr)
			break
		}
	}
	//最后为节点关联路由的回调方法 供用户处理业务逻辑
	root.handleFunc = handleFunc
	root.route = path
}

// 路由查找 获取路由对应的处理方法
func (r *Router) findRoute(method, path string) (*node, bool) {
	if _, ok := r.trees[method]; !ok {
		return nil, false
	}

	root := r.trees[method]
	//路径切割 "/"
	subStrs := strings.Split(path, "/")
	for index, subStr := range subStrs {
		if len(subStr) == 0 && (index == 0 || index == len(subStrs)-1) {
			// "/" 将导致第一个元素未空字符 "root/" 将导致最后一个元素为空字符串
			continue
		}
		//通配符路由和子节点路由均为空时退出匹配
		ok := false
		var tmpNode *node
		if root.child != nil {
			//获取静态匹配路由
			tmpNode, ok = root.child[subStr]
		}
		if !ok && root.wildcard != nil {
			//获取通配符路由
			tmpNode, ok = root.wildcard["*"]
		}
		if ok {
			root = tmpNode
		}
	}
	return root, true
}
