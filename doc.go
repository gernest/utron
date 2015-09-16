/*
Package utron is a lighweight MVC framework for building fast, scallable and robust web applications

example hello world in utron

	type Hello struct {
		*BaseController
	}

	func (h *Hello) World() {
		h.Ctx.Write([]byte("hello world"))
		h.String(http.StatusOK)
	}

*/
package utron
