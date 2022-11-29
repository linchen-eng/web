package web

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// TestAddRouter 测试路由添加 & 查找 方法
func TestAddAndFindRouter(t *testing.T) {
	type mockRouter struct {
		method     string
		path       string
		handleFunc HandleFunc
	}
	mockHandleFunc := func(ctx *Context) string {
		return ctx.route
	}
	cases := make([]*mockRouter, 0)
	cases = append(cases, &mockRouter{method: "POST", path: "/user/add", handleFunc: mockHandleFunc})
	cases = append(cases, &mockRouter{method: "POST", path: "/user/del", handleFunc: mockHandleFunc})
	cases = append(cases, &mockRouter{method: "POST", path: "/user/save", handleFunc: mockHandleFunc})
	cases = append(cases, &mockRouter{method: "GET", path: "/user/getList", handleFunc: mockHandleFunc})
	cases = append(cases, &mockRouter{method: "GET", path: "user/getList", handleFunc: mockHandleFunc})
	cases = append(cases, &mockRouter{method: "GET", path: "/user/*", handleFunc: mockHandleFunc})

	r := NewRouter()
	for _, v := range cases {
		r.addRouter(v.method, v.path, v.handleFunc)
		node, ok := r.findRoute(v.method, v.path)
		assert.True(t, ok, "路由注册错误未找到")
		ctx := &Context{
			req:   nil,
			resp:  nil,
			route: v.path,
		}
		retPath := node.handleFunc(ctx)
		assert.Equal(t, v.path, retPath)
	}
	//通配符路由测试
	casePath := "/user/update/V1"
	_, ok := r.findRoute("GET", casePath)
	if !ok {
		t.Error("通配符匹配路由失败:", casePath)
	}
	//panic 测试
	assert.PanicsWithValue(t, "路由格式不支持:路径中不能包含'//'这样的字符串", func() {
		r.addRouter(http.MethodGet, "/user//add", mockHandleFunc)
	})
	assert.PanicsWithValue(t, "路由格式不支持:路径中不能包含'//'这样的字符串", func() {
		r.addRouter(http.MethodGet, "//////", mockHandleFunc)
	})
	assert.PanicsWithValue(t, "路由格式不支持:路径中不能包含'//'这样的字符串", func() {
		r.addRouter(http.MethodGet, "//user//add//", mockHandleFunc)
	})
}

func BenchmarkAddAndFindRouter(b *testing.B) {
	type mockRouter struct {
		method     string
		path       string
		handleFunc HandleFunc
	}
	mockHandleFunc := func(ctx *Context) string {
		return ctx.route
	}
	cases := make([]*mockRouter, 0)
	cases = append(cases, &mockRouter{method: "POST", path: "/user/add", handleFunc: mockHandleFunc})
	cases = append(cases, &mockRouter{method: "POST", path: "/user/del", handleFunc: mockHandleFunc})
	cases = append(cases, &mockRouter{method: "POST", path: "/user/save", handleFunc: mockHandleFunc})
	cases = append(cases, &mockRouter{method: "GET", path: "/user/getList", handleFunc: mockHandleFunc})
	cases = append(cases, &mockRouter{method: "GET", path: "user/getList", handleFunc: mockHandleFunc})
	cases = append(cases, &mockRouter{method: "GET", path: "/user/*", handleFunc: mockHandleFunc})
	cases = append(cases, &mockRouter{method: "GET", path: "/message/*", handleFunc: mockHandleFunc})
	cases = append(cases, &mockRouter{method: "GET", path: "/message/send", handleFunc: mockHandleFunc})
	cases = append(cases, &mockRouter{method: "GET", path: "/message/del", handleFunc: mockHandleFunc})

	r := NewRouter()
	for i := 0; i < b.N; i++ {
		for _, v := range cases {
			r.addRouter(v.method, v.path, v.handleFunc)
			node, ok := r.findRoute(v.method, v.path)
			assert.True(b, ok, "路由注册错误未找到")
			ctx := &Context{
				req:   nil,
				resp:  nil,
				route: v.path,
			}
			retPath := node.handleFunc(ctx)
			assert.Equal(b, v.path, retPath)
		}
		//panic 测试
		assert.PanicsWithValue(b, "路由格式不支持:路径中不能包含'//'这样的字符串", func() {
			r.addRouter(http.MethodGet, "/user//add", mockHandleFunc)
		})
		assert.PanicsWithValue(b, "路由格式不支持:路径中不能包含'//'这样的字符串", func() {
			r.addRouter(http.MethodGet, "//////", mockHandleFunc)
		})
		assert.PanicsWithValue(b, "路由格式不支持:路径中不能包含'//'这样的字符串", func() {
			r.addRouter(http.MethodGet, "//user//add//", mockHandleFunc)
		})
	}
}
