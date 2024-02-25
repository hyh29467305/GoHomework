package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"
	"webook/internal/web"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginMiddlewareJWTBuilder struct {
}

func (m *LoginMiddlewareJWTBuilder) CheckLoginJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" || path == "/users/login_sms/code/send" || path == "/users/login_sms" {
			return
		}
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segs := strings.Split(authCode, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			// token解析错误，token是伪造的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			// token虽然解析出来了，但可能是非法的，或者已过期
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// if uc.UserAgent != ctx.GetHeader("User-Agent") {
		// 	ctx.AbortWithStatus(http.StatusUnauthorized)
		// 	return
		// }
		expireTime := uc.ExpiresAt
		if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
			tokenStr, err := token.SignedString(web.JWTKey)
			ctx.Header("x-jwt-header", tokenStr)
			if err != nil {
				log.Println(err)
			}
		}
		ctx.Set("user", uc)
	}
}
