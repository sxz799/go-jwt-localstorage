package main

import (
	"github.com/gin-gonic/gin"
	"go-jwt-study/router"
	"log"
)

func main() {
	gin.SetMode("debug")
	r := gin.Default()
	frontMode := true
	if frontMode {
		log.Println("已开启前后端整合模式！")
		r.LoadHTMLGlob("static/*")
	}
	router.RegRouter(r)
	log.Println("路由注册完成！当前端口为:8000")
	r.Run(":8000")
}
