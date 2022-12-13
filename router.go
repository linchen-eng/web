package web

import (
	"regexp"
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
	//参数子节点
	paramChild *node
	//当前节点的参数匹配的正则表达
	regexExpress string
}

// router 路由树的结构
type router struct {
	trees map[string]*node
}

// NewRouter 初始化路由
func NewRouter() router {
	return router{trees: make(map[string]*node)}
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

// paramOrCreate 为当前节点创建参数节点
func (n *node) paramOrCreate(path string) *node {
	//获取正则表达式
	startIndex := strings.Index(path, "(")
	endIndex := strings.Index(path, ")")
	regexExpress := ""
	key := path[1:]
	if startIndex > 0 && startIndex < endIndex {
		regexExpress = path[startIndex+1 : endIndex]
		key = path[1:startIndex]
	}
	n.paramChild = &node{path: key, regexExpress: regexExpress}
	return n.paramChild
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
func (r *router) addRouter(method, path string, handleFunc HandleFunc) {
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
		if subStr == "*" { //通配符路径
			root = root.wildcardOrCreate(subStr)
		} else if subStr[:1] == ":" { //参数路径
			root = root.paramOrCreate(subStr)
		} else { //静态匹配路径
			root = root.childOrCreate(subStr)
		}
	}
	//最后为节点关联路由的回调方法 供用户处理业务逻辑
	root.handleFunc = handleFunc
	root.route = path
}

// 路由查找 获取路由对应的处理方法
func (r *router) findRoute(method, path string) (*matchInfo, bool) {
	if _, ok := r.trees[method]; !ok {
		return nil, false
	}

	root := r.trees[method]
	//路径切割 "/"
	subStrs := strings.Split(path, "/")

	//路由匹配信息
	mi := &matchInfo{}
	for index, subStr := range subStrs {
		if len(subStr) == 0 && (index == 0 || index == len(subStrs)-1) {
			// "/" 将导致第一个元素未空字符 "root/" 将导致最后一个元素为空字符串
			continue
		}
		//通配符路由和子节点路由均为空时退出匹配
		ok := false
		var tmpNode *node
		matchWildcard := false
		if root.child != nil {
			//获取静态匹配路由
			tmpNode, ok = root.child[subStr]
		}
		if !ok && root.paramChild != nil {
			//获取当前参数节点的正则表达式
			var paramRet []string
			if root.paramChild.regexExpress != "" {
				re := regexp.MustCompile(root.paramChild.regexExpress)
				paramRet = re.FindAllString(root.paramChild.path, -1)
			}
			paramValue := subStr
			if len(paramRet) > 0 { //获取路由参数
				paramValue = paramRet[0]
			}
			mi.addValue(root.paramChild.path, paramValue)
			tmpNode = root.paramChild
			ok = true
		}
		if !ok && root.wildcard != nil {
			//获取通配符路由
			tmpNode, ok = root.wildcard["*"]
			matchWildcard = true
		}
		if ok {
			root = tmpNode
		} else if matchWildcard == false {
			return nil, false
		}
	}
	mi.n = root
	return mi, true
}

// 路由匹配信息
type matchInfo struct {
	n *node
	//路由匹配结果参数
	pathParams map[string]string
}

func (m *matchInfo) addValue(key string, value string) {
	if m.pathParams == nil {
		// 大多数情况，参数路径只会有一段
		m.pathParams = map[string]string{key: value}
	}
	m.pathParams[key] = value
}
