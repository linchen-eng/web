package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"regexp"
	"strings"
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
	cases = append(cases, &mockRouter{method: "GET", path: "/user/:userId/:action", handleFunc: mockHandleFunc})
	//正则参数路由
	cases = append(cases, &mockRouter{method: "GET", path: "/room/:userId(:id^[0-9]+$)", handleFunc: mockHandleFunc})

	r := NewRouter()
	for _, v := range cases {
		r.addRouter(v.method, v.path, v.handleFunc)
		mi, ok := r.findRoute(v.method, v.path)
		assert.True(t, ok, "路由注册错误未找到")
		ctx := &Context{
			req:   nil,
			resp:  nil,
			route: v.path,
		}
		retPath := mi.n.handleFunc(ctx)
		assert.Equal(t, v.path, retPath)
	}
	//通配符路由测试 /user/*
	casePath := "/user/update/V1"
	_, ok := r.findRoute("GET", casePath)
	if !ok {
		t.Error("通配符匹配路由失败:", casePath)
	}

	//通配符路由测试 /user/*
	casePath2 := "/user/update"
	_, ok = r.findRoute("GET", casePath2)
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
	//参数路由匹配测试
	caseParamPath := "/user/123/add"
	mi, ok := r.findRoute("GET", caseParamPath)
	assert.Equal(t, "123", mi.pathParams["userId"])
	assert.Equal(t, "add", mi.pathParams["action"])

	//参数路由匹配测试
	caseParamPath2 := "/user/123/add/ass"
	mi, ok = r.findRoute("GET", caseParamPath2)
	if ok {
		t.Error("参数路由匹配失败")
	}

	//参数路由匹配测试 正则表达式
	caseParamPath3 := "/room/xx"
	mi, ok = r.findRoute("GET", caseParamPath3)
	fmt.Println(mi)
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
			mi, ok := r.findRoute(v.method, v.path)
			assert.True(b, ok, "路由注册错误未找到")
			ctx := &Context{
				req:   nil,
				resp:  nil,
				route: v.path,
			}
			retPath := mi.n.handleFunc(ctx)
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

func TestRegexExpress(t *testing.T) {
	test := "(:id^[0-9]+$"
	startIndex := strings.Index(test, "(")
	endIndex := strings.Index(test, ")")
	fmt.Println("start>>", startIndex, " end>>", endIndex)
	test = test[startIndex+1 : endIndex]
	assert.Equal(t, "^[0-9]+$", test)

	re := regexp.MustCompile(test)
	ss := re.FindAllString("123", -1)
	fmt.Println(ss)
}
