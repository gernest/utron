/*
Package utron is a lightweight MVC framework for building fast, scalable and robust web applications

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
