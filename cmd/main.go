package main

import (
	"ckikoo/search/model/index"
	"ckikoo/search/router"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func HttpServer() *http.Server {
	r := gin.Default()

	handle := router.InitRouter(r)
	server := &http.Server{
		Addr:    ":12300",
		Handler: handle,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %s\n", err)
		}
	}()

	return server
}

func Init() {
	index.GetInstance().BuildIndex("./data")
}

func main() {
	server := HttpServer()
	go Init()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("服务关闭错误: ", err)
	}

	log.Println("服务已关闭")
}
