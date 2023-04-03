package router

import (
	"github.com/gin-gonic/gin"
	"go-jwt-study/middleware"
	"go-jwt-study/model"
	"log"
)

func login(c *gin.Context) {
	var u model.User
	err := c.BindJSON(&u)
	if err != nil {
		c.String(200, "参数有误！")
		return
	}
	if u.Username == "1" && u.Password == "1" {

	} else {
		c.String(200, "用户名或密码有误！")
		return
	}
	tokens, err := middleware.GenToken(u.Username)
	if err != nil {
		log.Println(err)
		c.String(200, "token生成失败")
	} else {
		c.JSON(200, gin.H{
			"tokens": tokens,
			"user":   u.Username,
		})
	}

}
func logout(c *gin.Context) {

	c.String(200, "退出成功")
}

func home(c *gin.Context) {

	c.JSON(200, gin.H{
		"tUser": '1',
	})
}

func index(c *gin.Context) {
	c.HTML(200, "index.html", "")
}

func User(e *gin.Engine) {
	e.GET("/", index)
	e.POST("/login", login)
	e.GET("/home", middleware.JWTAuth(), home)
	e.GET("/logout", logout)

}
