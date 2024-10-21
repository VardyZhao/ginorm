package app

import (
	"context"
	"errors"
	"ginorm/config"
	"ginorm/db"
	"ginorm/logger"
	"ginorm/middleware"
	"ginorm/router"
	"ginorm/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Init() *gin.Engine {
	// 加载环境变量
	config.LoadEnv()
	// 加载配置
	config.LoadConfig(util.GetAbsPath("/config.yaml"))
	// 设置gin模式
	gin.SetMode(config.Conf.GetString("mode"))
	// 加载日志组件
	logger.Load()
	// 读取翻译文件
	if err := config.LoadLocales(util.GetAbsPath("/locales/zh-cn.yaml")); err != nil {
		logger.Writer.Error("翻译文件加载失败: ", logger.NewError(err))
	}
	// 连接数据库
	db.Load()
	// todo 连接redis
	// 初始化gin
	r := gin.Default()
	// 加载中间件
	middleware.Load(r)
	// 加载路由
	router.Load(r)
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
			logger.Writer.Error("listen: %s\n", logger.NewError(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Writer.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 关闭所有db链接
	defer db.Conn.CloseAll()
	// 程序结束时，保证缓冲区里的所有内容都刷到文件
	defer logger.Flush()
	if shutdownErr := srv.Shutdown(ctx); shutdownErr != nil {
		logger.Writer.Fatal("Server Shutdown:", logger.NewError(shutdownErr))
	}
	logger.Writer.Info("Server exiting")
}
