package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SystemConfig 系统配置
type SystemConfig struct {
	Port      int    `json:"port" yaml:"port"`             // 服务端口
	Mode      string `json:"mode" yaml:"mode"`             // Gin 运行模式: debug, release, test
	LogLevel  string `json:"log_level" yaml:"log_level"`   // 日志级别
	HotReload bool   `json:"hot_reload" yaml:"hot_reload"` // 是否启用热加载
}

// MockEndpoint 定义单个 mock 接口
type MockEndpoint struct {
	Method     string            `json:"method"`      // HTTP 方法: GET, POST, PUT, DELETE 等
	Path       string            `json:"path"`        // 接口路径
	StatusCode int               `json:"status_code"` // 返回的状态码
	Response   interface{}       `json:"response"`    // 返回的数据
	Headers    map[string]string `json:"headers"`     // 自定义响应头
	Delay      int               `json:"delay"`       // 延迟时间(毫秒)
}

// EndpointsConfig Mock 接口配置
type EndpointsConfig struct {
	Endpoints []MockEndpoint `json:"endpoints"` // mock 接口列表
}

// LoadSystemConfig 加载系统配置（支持 YAML 和 JSON）
func LoadSystemConfig(filePath string) (*SystemConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config SystemConfig

	// 根据文件扩展名选择解析器
	if strings.HasSuffix(filePath, ".yml") || strings.HasSuffix(filePath, ".yaml") {
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			return nil, err
		}
	} else {
		// 默认使用 JSON
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			return nil, err
		}
	}

	// 设置默认值
	if config.Port == 0 {
		config.Port = 8080
	}
	if config.Mode == "" {
		config.Mode = "release"
	}
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	return &config, nil
}

// LoadEndpointsConfig 加载单个接口配置文件
func LoadEndpointsConfig(filePath string) (*EndpointsConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config EndpointsConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	// 为每个端点设置默认值
	for i := range config.Endpoints {
		if config.Endpoints[i].Method == "" {
			config.Endpoints[i].Method = "GET"
		}
		if config.Endpoints[i].StatusCode == 0 {
			config.Endpoints[i].StatusCode = 200
		}
	}

	return &config, nil
}

// LoadEndpointsFromDir 从目录加载所有接口配置
func LoadEndpointsFromDir(dirPath string) ([]MockEndpoint, error) {
	var allEndpoints []MockEndpoint

	// 读取目录下的所有 JSON 文件
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// 只处理 .json 文件
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			filePath := filepath.Join(dirPath, file.Name())
			config, err := LoadEndpointsConfig(filePath)
			if err != nil {
				return nil, err
			}
			allEndpoints = append(allEndpoints, config.Endpoints...)
		}
	}

	return allEndpoints, nil
}
