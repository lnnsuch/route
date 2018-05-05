package route

import (
	"path"
	"reflect"
	"regexp"
)

type GroupRoute struct {
	path   string
	handle []reflect.Value
	serve  *Serve
}

// 用户组
func (g *GroupRoute) Group(path string, handle ...interface{}) *GroupRoute {
	group := &GroupRoute{
		path:  g.combinePath(path),
		serve: g.serve,
	}
	group.combineHandle(handle...)
	return group
}

// 添加公共函数
func (g *GroupRoute) Use(handle ...interface{}) {
	g.combineHandle(handle...)
}

// 拼接url path
func (g *GroupRoute) combinePath(p string) string {
	return path.Join(g.path, p)
}

// 拼接组方法
func (g *GroupRoute) combineHandle(handle ...interface{}) {
	for _, v := range handle {
		g.handle = append(g.handle, reflect.ValueOf(v))
	}
}

// 添加get路由方法
func (g *GroupRoute) Get(route string, handle interface{}) {
	g.addRoute("GET", route, handle)
}

// 添加post路由方法
func (g *GroupRoute) Post(route string, handle interface{}) {
	g.addRoute("POST", route, handle)
}

// 添加路由方法
func (g *GroupRoute) addRoute(method, r string, handle interface{}) {
	r = path.Join(g.path, r)
	reg, err := regexp.Compile(r)
	if err != nil {
		g.serve.logger.Println("路由解析错误", err)
		return
	}
	h := reflect.ValueOf(handle)
	route := route{url: r, regUrl: reg, method: method, handle: h, groupHandle: g.handle}
	g.serve.routes = append(g.serve.routes, route)
}
