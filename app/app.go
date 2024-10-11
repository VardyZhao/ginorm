package app

import (
	"context"
	"errors"
	"ginorm/config"
	"ginorm/db"
	"ginorm/router"
	"ginorm/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Init() *gin.Engine {
	// 加载环境变量
	config.LoadEnv()
	// 加载配置
	config.LoadConfig(config.Env.RootDir + config.Env.Separate + "config.yaml")
	// 设置gin模式
	gin.SetMode(config.Conf.GetString("mode"))
	// 设置日志级别
	util.BuildLogger(config.Conf.GetString("log_level"))
	// 读取翻译文件
	if err := config.LoadLocales(config.Env.RootDir + config.Env.Separate + "locales" + config.Env.Separate + "zh-cn.yaml"); err != nil {
		util.Log().Panic("翻译文件加载失败", err)
	}
	// 连接数据库
	db.Load()
	defer db.Conn.CloseAll()
	// todo 连接redis

	// 加载中间件和路由
	r := router.NewRouter()
	return r
}

func Run(r *gin.Engine) {
	port := config.Conf.GetString("server.port")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if shutdownErr := srv.Shutdown(ctx); shutdownErr != nil {
		log.Fatal("Server Shutdown:", shutdownErr)
	}
	log.Println("Server exiting")

}
