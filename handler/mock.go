package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"mockmate/config"
)

// CreateMockHandler 创建 mock 接口的处理器
func CreateMockHandler(endpoint config.MockEndpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果设置了延迟，先等待
		if endpoint.Delay > 0 {
			time.Sleep(time.Duration(endpoint.Delay) * time.Millisecond)
		}

		// 设置自定义响应头
		for key, value := range endpoint.Headers {
			c.Header(key, value)
		}

		// 返回响应
		c.JSON(endpoint.StatusCode, endpoint.Response)
	}
}
