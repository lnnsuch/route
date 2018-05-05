package route

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

type ServeConfig struct {
	staticDir string
}

type Serve struct {
	GroupRoute
	config ServeConfig
	routes []route
	logger *log.Logger
}

// 新建http服务
func NewServe() *Serve {
	serve := &Serve{
		config: ServeConfig{staticDir: "/public/"},
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
	serve.GroupRoute.serve = serve
	return serve
}

// 添加get路由
func Get(route string, handle interface{}) {
	mainServe.Get(route, handle)
}

// 添加post路由
func Post(route string, handle interface{}) {
	mainServe.Post(route, handle)
}

// 运行http服务
func Run(addr string) {
	mainServe.Run(addr)
}

// 添加公共函数
func Use(handle ...interface{}) {
	mainServe.Use(handle...)
}

// 用户组
func Group(route string, handle ...interface{}) *GroupRoute {
	return mainServe.Group(route, handle...)
}

// 运行http服务
func (s *Serve) Run(addr string) {
	mux := http.NewServeMux()
	mux.Handle("/", s)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Sprintf("tcp链接创建失败 err:%s", err))
	}
	s.logger.Println("serve start, port:", addr)
	err = http.Serve(l, mux)
	if err != nil {
		panic(fmt.Sprintf("http链接创建失败 err:%s", err))
	}
	s.logger.Println("serve stop, port:", addr)
}

// http服务回调方法
func (s *Serve) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handle(w, r)
}

// 返回response信息
func (s *Serve) Handle(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path
	wd, _ := os.Getwd()
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		if s.tryFileServe(requestPath) {
			http.ServeFile(w, r, path.Join(wd, requestPath))
			return
		}
	}
	s.routeHandle(w, r)
}

// 路由方法处理
func (s *Serve) routeHandle(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path
	context := Context{serve: s, r: r, w: w}
	for _, route := range s.routes {
		if route.method != r.Method {
			continue
		}
		if !route.regUrl.MatchString(requestPath) {
			continue
		}
		match := route.regUrl.FindStringSubmatch(requestPath)
		if match[0] != requestPath {
			continue
		}
		if s.runGroupHandle(route, &context) {
			s.runRouteHandle(route, match[1:], &context)
		}
		return
	}
	context.NotFound()
}

// 执行路由组方法
// 返回值 bool 是否往下执行 true：是 false：否
func (s *Serve) runGroupHandle(route route, context *Context) bool {
	var handleType reflect.Type
	args := make([]reflect.Value, 1)
	for _, handle := range route.groupHandle {
		handleType = handle.Type()
		if requireContext(handleType) {
			args[0] = reflect.ValueOf(context)
		} else {
			args = nil
		}
		if content, ok := s.runCallHandle(handle, args, context); ok {
			context.setContentTypeJson()
			context.Write(content)
			return false
		}
	}
	return true
}

// 执行路由方法
func (s *Serve) runRouteHandle(route route, match []string, context *Context) {
	handleType := route.handle.Type()
	var args []reflect.Value
	if requireContext(handleType) {
		args = append(args, reflect.ValueOf(context))
	}
	for _, v := range match {
		args = append(args, reflect.ValueOf(v))
	}
	if content, ok := s.runCallHandle(route.handle, args, context); ok {
		context.setContentTypeJson()
		context.Write(content)
	}
}

func (s *Serve) runCallHandle(handle reflect.Value, args []reflect.Value, context *Context) ([]byte, bool) {
	res, err := s.callHandle(handle, args)
	if err != nil {
		context.ServeError()
		return nil, true
	}
	var content []byte
	if len(res) > 0 {
		res0 := res[0]
		if res0.Kind() == reflect.String {
			content = []byte(res0.String())
		} else if res0.Kind() == reflect.Slice && res0.Type().Elem().Kind() == reflect.Uint8 {
			content = res0.Interface().([]byte)
		}
		return content, true
	} else {
		return nil, false
	}
}

// 判断是否返回静态文件
func (s *Serve) tryFileServe(path string) bool {
	return strings.HasPrefix(path, s.config.staticDir)
}

// 调用路由方法
func (s *Serve) callHandle(f reflect.Value, args []reflect.Value) (value []reflect.Value, err interface{}) {
	defer func() {
		if e := recover(); e != nil {
			value = nil
			err = e
			for i := 0; ; i++ {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				s.logger.Println(file, line)
			}
		}
	}()
	return f.Call(args), nil
}

// 是否需要写入Context参数
func requireContext(handleType reflect.Type) bool {
	if handleType.NumIn() == 0 {
		return false
	}
	args0 := handleType.In(0)
	if args0.Kind() != reflect.Ptr {
		return false
	}
	if args0.Elem() == reflect.TypeOf(Context{}) {
		return true
	}
	return false
}

type route struct {
	url         string          // url
	regUrl      *regexp.Regexp  // 正则表达式
	method      string          // http方法
	handle      reflect.Value   // 处理函数
	groupHandle []reflect.Value // 是否是路由组
}

var mainServe = NewServe()
