package main

import (
	"context"
	"demo_webapp/controller"
	"demo_webapp/dao/mysql"
	"demo_webapp/dao/redis"
	"demo_webapp/logger"
	"demo_webapp/pkg/snowflake"
	"demo_webapp/router"
	"demo_webapp/setting"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//加载配置文件
	if err := setting.Init(); err != nil {
		fmt.Printf("setting Init failed,err:%v\n", err)
		return
	}
	//初始化zap日志

	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("logger init failer:%v\n", err)
		return
	}

	//初始化mysql
	if err := mysql.Init(setting.Conf.MysqlConfig); err != nil {
		fmt.Printf("mysql init failed,err:%v\n", err)
		return
	}
	defer mysql.Close()

	//初始化redis
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		fmt.Printf("redis Init failed,err:%v\n", err)
		return
	}
	defer redis.Close()

	//加载翻译
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("InitTrans failed,err:%v\n", err)
		return
	}

	//雪花算法
	if err := snowflake.Init(setting.Conf.StartTime, setting.Conf.MachineID); err != nil {
		fmt.Printf("snowflake Init failed,err:%v\n", err)
		return
	}

	//优雅关机
	r := router.SetupRouter()
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", setting.Conf.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("ListenAndServe failed,err:%v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit //当从一个队列中取值，并且没有收到close()通知时，会进入一个阻塞状态

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server Shutdown,err:%v\n", err)
		return
	}
	fmt.Println("Server exiting")

}
