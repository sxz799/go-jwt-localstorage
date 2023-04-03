package middleware

import (
	"encoding/json"
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

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenToken(u string) (tokens Tokens, err error) {
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
	tokens.AccessToken, err = accessToken.SignedString(JwtKey)
	tokens.RefreshToken, err = refreshToken.SignedString(JwtKey)
	if err != nil {
		// 创建Token失败
		return tokens, err
	}
	return tokens, nil
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokens Tokens
		tokensStr := c.GetHeader("tokens")
		log.Println(tokensStr)
		err := json.Unmarshal([]byte(tokensStr), &tokens)
		if err != nil {
			// 获取access-token失败
			c.Abort()
			return
		}

		accessClaims := &Claims{}
		refreshClaims := &Claims{}
		accessToken, err := jwt.ParseWithClaims(tokens.AccessToken, accessClaims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil || !accessToken.Valid {
			refreshToken, err2 := jwt.ParseWithClaims(tokens.RefreshToken, refreshClaims, func(token *jwt.Token) (interface{}, error) {
				return JwtKey, nil
			})
			if err2 != nil || !refreshToken.Valid {
				c.String(200, "refreshToken不合法或已过期,需要重新登录")
				c.Abort()
				return
			}
			if refreshToken.Valid {
				tokens, err = GenToken(accessClaims.Username)
				if err != nil {
					c.String(200, "更新token失败！")
					c.Abort()
					return
				}
			}
		}
		log.Println("accessToken过期时间:", accessClaims.ExpiresAt)
		log.Println("refreshToken过期时间:", refreshClaims.ExpiresAt)
		c.Set("accessClaims", accessClaims)
		marshal, _ := json.Marshal(tokens)
		c.Header("tokens", string(marshal))
		c.Next()
	}

}
