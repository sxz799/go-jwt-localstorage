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

func GenToken(u string) (accessTokenStr, refreshTokenStr string, err error) {
	// 创建jwt accessClaims 设置过期时间15s
	accessClaims := &Claims{
		Username: u,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Second)),
		},
	}
	refreshClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Second)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	accessTokenStr, err = accessToken.SignedString(JwtKey)
	refreshTokenStr, err = refreshToken.SignedString(JwtKey)
	if err != nil {
		// 创建Token失败
		return "", "", err
	}
	return accessTokenStr, refreshTokenStr, nil
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessTokenStr, err := c.Cookie("access-token")
		if err != nil {
			// 获取access-token失败
			c.String(200, "获取token失败")
			c.Abort()
			return
		}

		accessClaims := &Claims{}
		refreshClaims := &Claims{}
		accessToken, err := jwt.ParseWithClaims(accessTokenStr, accessClaims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		log.Println("accessToken过期时间:", accessClaims.ExpiresAt)
		if err != nil || !accessToken.Valid {
			refreshTokenStr, err1 := c.Cookie("refresh-token")
			refreshToken, err2 := jwt.ParseWithClaims(refreshTokenStr, refreshClaims, func(token *jwt.Token) (interface{}, error) {
				return JwtKey, nil
			})
			log.Println("refreshToken过期时间:", refreshClaims.ExpiresAt)
			if err1 != nil || err2 != nil || !refreshToken.Valid {
				c.String(200, "refreshToken不合法或已过期,需要重新登录")
				c.Abort()
				return
			}
			if refreshToken.Valid {
				accessTokenStr, refreshTokenStr, err := GenToken(accessClaims.Username)
				if err != nil {
					c.String(200, "更新token失败！")
					c.Abort()
					return
				}
				c.SetCookie("access-token", accessTokenStr, 3600, "", "", false, true)
				c.SetCookie("refresh-token", refreshTokenStr, 3600, "", "", false, true)
			}
		}
		c.Set("accessClaims", accessClaims)
		c.Next()
	}

}
