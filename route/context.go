package route

import "net/http"

type Context struct {
	serve *Serve
	r     *http.Request
	w     http.ResponseWriter
}

// 服务错误 500
func (c *Context) ServeError() {
	c.w.WriteHeader(http.StatusInternalServerError)
	c.Write([]byte("server error"))
}

// 资源未找到 404
func (c *Context) NotFound() {
	c.w.WriteHeader(http.StatusNotFound)
	c.Write([]byte("page not found"))
}

// 设置header bool:header包中是否存在,不存在添加
func (c *Context) SetHeader(key, value string, bool bool) {
	if bool {
		c.w.Header().Set(key, value)
	} else {
		c.w.Header().Add(key, value)
	}
}

// 设置content-type为文本格式
func (c *Context) SetContentTypeText() {
	c.SetHeader("Content-Type", "text/html; charset=utf-8", true)
}

// 设置content-type为json格式
func (c *Context) setContentTypeJson() {
	c.SetHeader("Content-Type", "application/json; charset=utf-8", true)
}

// 返回response包
func (c *Context) Write(b []byte) {
	_, err := c.w.Write(b)
	if err != nil {
		c.serve.logger.Println("Error http write", err)
	}
}
