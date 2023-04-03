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

	accessTokenStr, refreshTokenStr, err := middleware.GenToken(u.Username)
	if err != nil {
		log.Println(err)
		c.String(200, "token生成失败")
	}
	c.SetCookie("access-token", accessTokenStr, 3600, "", "", false, true)
	c.SetCookie("refresh-token", refreshTokenStr, 3600, "", "", false, true)
	c.JSON(200, gin.H{
		"access-token":           accessTokenStr,
		"access-refreshTokenStr": refreshTokenStr,
		"user":                   u.Username,
	})

}
func logout(c *gin.Context) {

	c.SetCookie("access-token", "-", 0, "", "", false, true)
	c.SetCookie("refresh-token", "-", 0, "", "", false, true)
	c.String(200, "退出成功")
}

func home(c *gin.Context) {
	c.HTML(200, "home.html", nil)

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
