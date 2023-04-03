package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

var JwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenToken(u string) (tokenStr string, err error) {
	// 创建jwt accessClaims 设置过期时间15s
	claims := &Claims{
		Username: u,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString(JwtKey)
	return
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokensStr := c.GetHeader("token")
		log.Println("获取到的token", tokensStr)

		if tokensStr == "" {
			// 获取access-token失败
			c.String(200, "未登录状态！")
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokensStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		log.Println(token.Valid)
		if err == nil && token.Valid {
			if claims.ExpiresAt.Unix()-time.Now().Unix() < 15 {
				log.Println("生成了新Token")
				str, _ := GenToken(claims.Username)
				c.Header("new-token", str)
			}
			c.Set("claims", claims)
			c.Next()
		} else {
			c.String(200, "token已过期，请重新登录！")
			c.Abort()
			return
		}

	}

}
