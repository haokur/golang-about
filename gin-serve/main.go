package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type student struct {
	Name string
	Age  int8
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("请求前---")
		token := c.Query("token")
		if token == "" {
			// panic("请登录，携带token")
			c.String(http.StatusForbidden, "请先登录")
		} else {
			c.Next()
		}
		fmt.Println("请求后---")
	}
}

func main() {
	fmt.Println("serve is running on http://localhost:8080")
	r := gin.Default()
	// 路由方法有 GET,POST,PUT,PATCH,DELETE,OPTIONS,Any
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World34")
	})
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "hello %s", name)
	})
	// name为路径参数，必选，role为可选
	r.GET("/user/:name/*role", func(c *gin.Context) {
		name := c.Param("name")
		role := c.Param("role")

		id := c.Query("id")
		age := c.DefaultQuery("age", "18")
		c.String(http.StatusOK, "hello name is %s,role is %s,id is %s,age is %s", name, role, id, age)
	})

	// POST请求
	// 测试：curl http://localhost:8090/form  -X POST -d 'username=geektutu&password=1234'
	r.POST("/form", func(c *gin.Context) {
		id := c.Query("id")
		age := c.DefaultQuery("age", "18")
		username := c.PostForm("username")
		password := c.DefaultPostForm("password", "123456")

		c.JSON(http.StatusOK, gin.H{
			"id":       id,
			"age":      age,
			"username": username,
			"password": password,
		})
	})

	// 字典参数
	// 测试：curl -g "http://localhost:8090/post?ids[Jack]=001&ids[Tom]=002" -X POST -d 'names[a]=Sam&names[b]=David'
	r.POST("/post", func(c *gin.Context) {
		ids := c.QueryMap("ids")
		names := c.PostFormMap("names")
		// {"ids":{"Jack":"001","Tom":"002"},"names":{"a":"Sam","b":"David"}}
		c.JSON(http.StatusOK, gin.H{
			"ids":   ids,
			"names": names,
		})
	})

	// 重定向
	// 地址会变成/
	r.GET("/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	// 地址不变，内容定位到 /
	r.GET("/go_index", func(c *gin.Context) {
		c.Request.URL.Path = "/"
		r.HandleContext(c)
	})

	// 分组路由
	defaultHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"path": c.FullPath(),
		})
	}

	// group:v1
	v1 := r.Group("/v1")
	{
		v1.GET("/posts", defaultHandler)
		v1.GET("/series", defaultHandler)
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/posts", defaultHandler)
		v2.GET("/series", defaultHandler)
	}

	// 上传文件
	// 单个文件
	r.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		c.String(http.StatusOK, "%s uploaded", file.Filename)
	})
	// 多个文件
	r.POST("/upload-multi", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)
		}
		c.String(http.StatusOK, "%d files uploaded", len(files))
	})

	// HTML模板
	r.LoadHTMLGlob("./templates/*")
	stu1 := &student{Name: "tutu", Age: 18}
	stu2 := &student{Name: "jack", Age: 22}
	r.GET("/user", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user.html", gin.H{
			"title":  "Gin",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	// 中间件
	// 作用于全局
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 作用于某个组
	authorized := r.Group("/auth")
	authorized.Use(AuthRequired())
	{
		authorized.GET("/login", func(c *gin.Context) {
			c.String(http.StatusOK, "login now!")
		})
	}
	r.Run(":8080")
}
