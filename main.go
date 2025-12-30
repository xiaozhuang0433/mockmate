package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"mockmate/config"
	"mockmate/handler"
	"mockmate/reload"
	"mockmate/watcher"
)

func main() {
	// 解析命令行参数
	systemConfigFile := flag.String("config", "configs/system.yml", "系统配置文件路径")
	endpointsDir := flag.String("endpoints", "configs/endpoints", "接口配置目录路径")
	flag.Parse()

	// 加载系统配置
	sysCfg, err := config.LoadSystemConfig(*systemConfigFile)
	if err != nil {
		log.Fatalf("加载系统配置失败: %v", err)
	}

	// 加载接口配置
	endpoints, err := config.LoadEndpointsFromDir(*endpointsDir)
	if err != nil {
		log.Fatalf("加载接口配置失败: %v", err)
	}

	// 设置 Gin 运行模式
	gin.SetMode(sysCfg.Mode)
	r := gin.Default()

	// 创建动态处理器
	dynamicHandler := handler.NewDynamicHandler()
	dynamicHandler.UpdateEndpoints(endpoints)

	// 注册通配符路由，处理所有请求
	r.Any("/*path", dynamicHandler.Handle)

	// 如果启用热加载，启动文件监听
	if sysCfg.HotReload {
		reloadManager := reload.NewManager(dynamicHandler, *endpointsDir)

		// 创建文件监听器
		w, err := watcher.NewWatcher()
		if err != nil {
			log.Fatalf("创建文件监听器失败: %v", err)
		}
		defer w.Close()

		// 监听接口配置目录
		err = w.WatchDir(*endpointsDir, reloadManager.ReloadEndpoints)
		if err != nil {
			log.Fatalf("监听接口配置目录失败: %v", err)
		}

		// 启动监听
		w.Start()
		log.Println("热加载已启用，配置文件修改后将自动重载")
	}

	// 启动服务
	addr := fmt.Sprintf(":%d", sysCfg.Port)
	log.Printf("Mock 服务启动成功，监听端口: %d", sysCfg.Port)
	log.Printf("已加载 %d 个 mock 接口", len(endpoints))

	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
