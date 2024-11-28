package middleware

import (
	"basic-go/webook/internal/web"
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

// JWT 登录见面
type LoginJWTMiddlewareBuilder struct {

}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc  {
	//用 go 的方式编码解码
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		//不需要校验的
		if ctx.Request.URL.Path == "/users/login" ||
			ctx.Request.URL.Path == "/users/signup" {
			return
		}
		// 利用 JWT校验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			//没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segs := strings.Split(tokenHeader, "")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := segs[1]
		claims := &web.UserClaims{}
		//这里一定要传指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("HRx0ToImlZtakubRzKfJ2NCSNGdRik6z"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			// 严重的安全问题
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
