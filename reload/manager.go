package reload

import (
	"log"
	"sync"

	"mockmate/config"
	"mockmate/handler"
)

// Manager 热加载管理器
type Manager struct {
	mu             sync.RWMutex
	dynamicHandler *handler.DynamicHandler
	endpoints      []config.MockEndpoint
	endpointsDir   string
}

// NewManager 创建热加载管理器
func NewManager(dynamicHandler *handler.DynamicHandler, endpointsDir string) *Manager {
	return &Manager{
		dynamicHandler: dynamicHandler,
		endpointsDir:   endpointsDir,
	}
}

// ReloadEndpoints 重新加载接口配置
func (m *Manager) ReloadEndpoints() {
	log.Println("检测到配置文件变化，正在重新加载...")

	// 加载新的接口配置
	newEndpoints, err := config.LoadEndpointsFromDir(m.endpointsDir)
	if err != nil {
		log.Printf("重新加载配置失败: %v", err)
		return
	}

	m.mu.Lock()
	m.endpoints = newEndpoints
	m.mu.Unlock()

	// 更新动态路由表
	m.dynamicHandler.UpdateEndpoints(newEndpoints)

	log.Printf("配置重载成功，已加载 %d 个 mock 接口", len(newEndpoints))
}

// GetEndpoints 获取当前接口列表（用于测试）
func (m *Manager) GetEndpoints() []config.MockEndpoint {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.endpoints
}
