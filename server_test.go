package web

import (
	"fmt"
	"testing"
)

func TestHTTPServer_Get(t *testing.T) {
	t.Parallel()

	s := NewHTTPServer()
	s.mdls = []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第一个before")
				next(ctx)
				fmt.Println("第一个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第二个before")
				next(ctx)
				fmt.Println("第二个after")
			}
		},
	}
	s.Get("/user", func(ctx *Context) {
		ctx.resp.Write([]byte("hello world!"))
	})
	s.Start(":8081")
}
