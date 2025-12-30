package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"mockmate/config"
	"mockmate/handler"
)

// RegisterEndpoints 注册所有 mock 接口
func RegisterEndpoints(router *gin.Engine, endpoints []config.MockEndpoint) {
	for _, endpoint := range endpoints {
		registerEndpoint(router, endpoint)
	}
}

// registerEndpoint 注册单个 mock 接口
func registerEndpoint(router *gin.Engine, endpoint config.MockEndpoint) {
	handlerFunc := handler.CreateMockHandler(endpoint)

	// 根据 HTTP 方法注册路由
	switch endpoint.Method {
	case "GET":
		router.GET(endpoint.Path, handlerFunc)
		log.Printf("注册 GET %s", endpoint.Path)
	case "POST":
		router.POST(endpoint.Path, handlerFunc)
		log.Printf("注册 POST %s", endpoint.Path)
	case "PUT":
		router.PUT(endpoint.Path, handlerFunc)
		log.Printf("注册 PUT %s", endpoint.Path)
	case "DELETE":
		router.DELETE(endpoint.Path, handlerFunc)
		log.Printf("注册 DELETE %s", endpoint.Path)
	case "PATCH":
		router.PATCH(endpoint.Path, handlerFunc)
		log.Printf("注册 PATCH %s", endpoint.Path)
	default:
		// 支持任意方法
		router.Any(endpoint.Path, handlerFunc)
		log.Printf("注册 ANY %s", endpoint.Path)
	}
}
