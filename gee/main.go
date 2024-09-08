package main

import (
	"fmt"
	"gee"
	"log"
	"net/http"
	"text/template"
	"time"
)

// // Engine is the uni handler for all requests
// type Engine struct{}

// func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	switch req.URL.Path {
// 	case "/":
// 		fmt.Fprintf(w, "URL.Path=%q\n", req.URL.Path)
// 	case "/hello-world":
// 		for k, v := range req.Header {
// 			fmt.Fprintf(w, "Header [%q]=%q\n", k, v)
// 		}
// 	default:
// 		fmt.Fprintf(w, "404 NOT FOUND:%s\n", req.URL)
// 	}
// }

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.New()
	r.Use(gee.Recovery())
	r.Use(gee.Logger()) // 全局middleware

	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "jack", Age: 20}
	stu2 := &student{Name: "jack2", Age: 22}
	r.GET("/stu", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "index.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
			"now":    time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.GET("/date", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "date.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.GET("/trace", func(ctx *gee.Context) {
		names := []string{"1", "2", "3"}
		ctx.String(http.StatusOK, names[4])
	})
	// r.GET("/", func(w http.ResponseWriter, req *http.Request) {
	// 	fmt.Fprintf(w, "URL.Path=%q\n", req.URL.Path)
	// })
	// r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
	// 	for k, v := range req.Header {
	// 		fmt.Fprintf(w, "Header[%q]=%q\n", k, v)
	// 	}
	// })
	r.GET("/", func(ctx *gee.Context) {
		ctx.String(200, "hello world")
	})
	r.GET("/json", func(ctx *gee.Context) {
		ctx.JSON(200, gee.H{
			"Name": "jack",
			"Age":  10,
		})
	})
	r.GET("/user/:id", func(ctx *gee.Context) {
		ctx.String(200, "userId is::"+ctx.Params["id"])
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/user", func(ctx *gee.Context) {
			ctx.String(200, "this is v1 user")
		})
	}

	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello v2")
		})
	}
	r.Run(":9999")
	// http.HandleFunc("/", indexHandler)
	// http.HandleFunc("/hello", helloHandler)

	// engine := new(Engine)
	// // 使用了engine，则上面单独声明的HandleFunc则会失效
	// log.Fatal(http.ListenAndServe(":8009", engine))
}

// func indexHandler(w http.ResponseWriter, req *http.Request) {
// 	fmt.Fprintf(w, "URL.Path=%q\n", req.URL.Path)
// }

// func helloHandler(w http.ResponseWriter, req *http.Request) {
// 	for k, v := range req.Header {
// 		fmt.Fprintf(w, "Header[%q]=%q\n", k, v)
// 	}
// }
