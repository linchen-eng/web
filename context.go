package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Context struct {
	req  *http.Request
	resp http.ResponseWriter
	//路由GET参数 缓存
	cacheQueryValues map[string][]string
	//路由
	route string
}

// ResponseJson 响应客户端json数据
func (ctx *Context) ResponseJson(code int, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	ctx.resp.WriteHeader(code)
	_, err = ctx.resp.Write(bs)
	return err
}

// BindJson 绑定json
func (ctx *Context) BindJson(val any) error {
	if ctx.req.Body == nil {
		return errors.New("web: body 为 nil")
	}
	decoder := json.NewDecoder(ctx.req.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}

// FormValue 从表单获取数据
func (ctx *Context) FormValue(key string) StringValue {
	if err := ctx.req.ParseForm(); err != nil {
		return StringValue{err: err}
	}
	return StringValue{val: ctx.req.FormValue(key)}
}

// QueryValue 路由GET参数
func (ctx *Context) QueryValue(key string) StringValue {
	if ctx.cacheQueryValues == nil {
		ctx.cacheQueryValues = ctx.req.URL.Query()
	}
	vals, ok := ctx.cacheQueryValues[key]
	if !ok {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}
	return StringValue{val: vals[0]}
}

// StringValue 参数返回值结构体
type StringValue struct {
	val string
	err error
}

// AsInt64 强转64位 int
func (sv *StringValue) AsInt64() (int64, error) {
	if sv.err != nil {
		return 0, sv.err
	}
	return strconv.ParseInt(sv.val, 10, 64)
}

// AsInt32 强转32位 int
func (sv *StringValue) AsInt32() (int32, error) {
	if sv.err != nil {
		return 0, sv.err
	}
	val, err := strconv.Atoi(sv.val)
	return int32(val), err
}

// AsInt16 强转16位 int
func (sv *StringValue) AsInt16() (int16, error) {
	if sv.err != nil {
		return 0, sv.err
	}
	val, err := strconv.Atoi(sv.val)
	return int16(val), err
}

// AsInt8 强转8位 int
func (sv *StringValue) AsInt8() (int8, error) {
	if sv.err != nil {
		return 0, sv.err
	}
	val, err := strconv.Atoi(sv.val)
	return int8(val), err
}
