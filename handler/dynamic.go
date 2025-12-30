package handler

import (
	"log"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"mockmate/config"
)

// DynamicHandler 动态路由处理器
type DynamicHandler struct {
	mu        sync.RWMutex
	endpoints map[string]*config.MockEndpoint // key: method:path
}

// NewDynamicHandler 创建动态处理器
func NewDynamicHandler() *DynamicHandler {
	return &DynamicHandler{
		endpoints: make(map[string]*config.MockEndpoint),
	}
}

// UpdateEndpoints 更新端点列表
func (h *DynamicHandler) UpdateEndpoints(endpoints []config.MockEndpoint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 清空旧的映射
	h.endpoints = make(map[string]*config.MockEndpoint)

	// 重建映射
	for i := range endpoints {
		key := endpoints[i].Method + ":" + endpoints[i].Path
		h.endpoints[key] = &endpoints[i]
		log.Printf("  注册路由: %s %s", endpoints[i].Method, endpoints[i].Path)
	}

	log.Printf("动态路由表已更新，共 %d 个接口", len(h.endpoints))
}

// Handle 通用处理器
func (h *DynamicHandler) Handle(c *gin.Context) {
	method := c.Request.Method
	path := c.Request.URL.Path

	// 首先尝试精确匹配
	key := method + ":" + path
	h.mu.RLock()
	endpoint, ok := h.endpoints[key]
	h.mu.RUnlock()

	if ok {
		h.executeEndpoint(c, endpoint)
		return
	}

	// 如果没有精确匹配，尝试路径参数模式匹配
	h.mu.RLock()
	defer h.mu.RUnlock()

	var matchedEndpoint *config.MockEndpoint
	for _, ep := range h.endpoints {
		if ep.Method == method && h.matchPathPattern(ep.Path, path) {
			matchedEndpoint = ep
			break
		}
	}

	if matchedEndpoint != nil {
		h.executeEndpoint(c, matchedEndpoint)
		return
	}

	// 没有找到匹配的路由
	c.JSON(404, gin.H{
		"code":    404,
		"message": "接口不存在",
		"path":    path,
		"method":  method,
	})
}

// executeEndpoint 执行端点处理
func (h *DynamicHandler) executeEndpoint(c *gin.Context, endpoint *config.MockEndpoint) {
	handler := CreateMockHandler(*endpoint)
	handler(c)
}

// matchPathPattern 匹配路径模式（支持 :param 参数）
func (h *DynamicHandler) matchPathPattern(pattern, path string) bool {
	// 完全相等
	if pattern == path {
		return true
	}

	// 分割路径
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	// 段数不同，不匹配
	if len(patternParts) != len(pathParts) {
		return false
	}

	// 逐段比较
	for i := 0; i < len(patternParts); i++ {
		patternPart := patternParts[i]
		pathPart := pathParts[i]

		// 如果模式段以 : 开头，表示参数，匹配任意值
		if strings.HasPrefix(patternPart, ":") {
			continue
		}

		// 普通段需要精确匹配
		if patternPart != pathPart {
			return false
		}
	}

	return true
}
