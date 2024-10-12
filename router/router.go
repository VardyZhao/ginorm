package router

import (
	"ginorm/controller/api"
	"ginorm/controller/api/user"
	"ginorm/middleware"
	"github.com/gin-gonic/gin"
)

func Load(r *gin.Engine) {
	// 路由
	v1 := r.Group("/api/v1")
	{
		v1.POST("ping", api.Ping)

		// 用户登录
		v1.POST("user/register", user.Register)

		// 用户登录
		v1.POST("user/login", user.Login)

		// 需要校验登录态得
		auth := v1.Group("")
		auth.Use(middleware.AuthRequired())
		{
			// User Routing
			auth.GET("user/profile", user.Profile)
			auth.DELETE("user/logout", user.Logout)
		}
	}
}
