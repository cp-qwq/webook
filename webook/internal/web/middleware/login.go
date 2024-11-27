package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc  {
	//用 go 的方式编码解码
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		//不需要校验的
		if ctx.Request.URL.Path == "/users/login" ||
			ctx.Request.URL.Path == "/users/signup" {
			return
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			//没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//每隔一段时间刷新
		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		now := time.Now()
		//刚登录，还没刷新过
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
		//updateTime存在
		updateTimeVal, _ := updateTime.(time.Time)
		if now.Sub(updateTimeVal) > time.Minute {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
	}
}
