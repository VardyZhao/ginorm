package middleware

import (
	"github.com/gin-gonic/gin"
	"os"
)

func Load(r *gin.Engine) {
	// 中间件, 顺序不能改
	r.Use(Logger())
	r.Use(PanicHandler())
	r.Use(ErrorHandler())
	r.Use(Session(os.Getenv("SESSION_SECRET")))
	r.Use(Cors())
	r.Use(CurrentUser())
}
