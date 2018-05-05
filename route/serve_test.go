package route

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func index() string {
	return "aaa"
}

func index2(string string) []byte {
	fmt.Println(string)
	fmt.Println("index2")
	return []byte("1111111")
}

func index3(context *Context, string string) string {
	fmt.Println("index3")
	return "2222222"
}

func public1(context *Context) string {
	fmt.Println("public1")
	return "public1"
}
func public2() {
	fmt.Println("public2")
}

func index4() {
	fmt.Println("index4")
}
func runHttp() {
	group1 := Group("/aa", public1, public2)
	{
		group1.Get("/zz", index)
		group1.Get("/zz/([0-9]+)", index2)
		group1.Get("/zz/([\\w]+)", index3)
	}
	Get("/bb", index4)
	Run("0.0.0.0:3000")
}

func TestRun1(t *testing.T) {
	go runHttp()
	response, _ := http.Get("http://127.0.0.1:3000")
	fmt.Println(response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}
func TestRun2(t *testing.T) {
	go runHttp()
	response, _ := http.Get("http://127.0.0.1:3000/bb")
	fmt.Println(response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}

func TestRun3(t *testing.T) {
	go runHttp()
	response, _ := http.Get("http://127.0.0.1:3000/aa/zz")
	fmt.Println(response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}

func TestRun4(t *testing.T) {
	go runHttp()
	response, _ := http.Get("http://127.0.0.1:3000/aa/zz/123")
	fmt.Println(response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}

func TestRun5(t *testing.T) {
	go runHttp()
	response, _ := http.Get("http://127.0.0.1:3000/aa/zz/abc")
	fmt.Println(response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}
