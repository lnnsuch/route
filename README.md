# route

route

group1 := route.Group("/aa", public1, public2)
{
	group1.Get("/zz", index)
	group1.Get("/zz/([0-9]+)", index2)
	group1.Get("/zz/([\\w]+)", index3)
}
route.Run("0.0.0.0:3000")

路由组函数的参数可以为空或者为 *route.Context

路由函数的第一个参数可以为 * route.Context
其他参数由正则结果值决定,例如index2匹配1个数字,则参数为一个
func index2(string string) {
}

函数可以有返回值,如果返回值为[]byte或者string,则会直接输出
路由组如果存在返回值,则后面的函数不会执行
