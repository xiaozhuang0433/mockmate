# MockMate

一个基于 Go 和 Gin 框架的轻量级 Mock 服务，可以通过配置文件快速 mock 任何 HTTP 接口。

## 功能特性

- 支持所有 HTTP 方法 (GET, POST, PUT, DELETE, PATCH 等)
- 支持动态路由参数 (如 `/api/user/:id`)
- 支持自定义状态码
- 支持自定义响应头
- 支持模拟延迟响应
- **系统配置与接口配置分离**
- **支持按模块分文件配置接口**
- **✨ 热加载 - 修改配置文件后自动重载，无需重启**

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置服务

**系统配置** (`configs/system.yml`)：
```yaml
# 服务监听端口
port: 9090

# Gin 运行模式：debug | release | test
mode: release

# 日志级别：debug | info | warn | error
log_level: info

# 是否启用热加载（修改配置文件后自动重载，无需重启）
hot_reload: true
```

**接口配置** (`configs/endpoints/*.json`)：
```json
{
  "endpoints": [
    {
      "method": "GET",
      "path": "/api/user/:id",
      "status_code": 200,
      "response": {
        "code": 0,
        "message": "success",
        "data": {
          "id": 1,
          "name": "张三"
        }
      }
    }
  ]
}
```

### 3. 启动服务

```bash
# 使用默认配置
go run .

# 指定系统配置文件
go run . -config configs/system.yml

# 指定接口配置目录
go run . -endpoints configs/endpoints
```

### 4. 编译运行

```bash
# 编译
go build -o mock-server

# 运行
./mock-server
```

## 配置说明

### 目录结构

```
configs/
├── system.yml              # 系统配置（端口、运行模式等）
└── endpoints/              # Mock 接口配置目录
    ├── user-api.json       # 用户相关接口
    ├── device-api.json     # 设备相关接口
    └── *.json              # 其他接口配置
```

### 系统配置项 (system.yml)

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `port` | int | 8080 | 服务监听端口 |
| `mode` | string | release | Gin 运行模式：`debug`（调试）、`release`（生产）、`test`（测试） |
| `log_level` | string | info | 日志级别 |
| `hot_reload` | bool | false | 是否启用热加载 |

### 接口配置项 (endpoints/*.json)

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `method` | string | GET | HTTP 方法：GET、POST、PUT、DELETE、PATCH 等 |
| `path` | string | - | 接口路径，支持路由参数如 `/api/user/:id` |
| `status_code` | int | 200 | HTTP 响应状态码 |
| `delay` | int | 0 | 响应延迟时间（毫秒） |
| `headers` | object | {} | 自定义响应头 |
| `response` | any | - | 响应数据，可以是任意 JSON 结构 |

## 配置示例

### 基础接口

```json
{
  "method": "GET",
  "path": "/api/users",
  "response": {
    "code": 0,
    "data": [
      {"id": 1, "name": "张三"},
      {"id": 2, "name": "李四"}
    ]
  }
}
```

### 带路径参数

```json
{
  "method": "GET",
  "path": "/api/user/:id",
  "response": {
    "code": 0,
    "data": {"id": 1, "name": "张三"}
  }
}
```

### 模拟延迟

```json
{
  "method": "GET",
  "path": "/api/slow",
  "delay": 2000,
  "response": {"message": "延迟2秒返回"}
}
```

### 自定义响应头

```json
{
  "method": "GET",
  "path": "/api/custom",
  "headers": {
    "X-Custom-Header": "custom-value"
  },
  "response": {"message": "带自定义响应头"}
}
```

## 测试接口

```bash
# GET 请求
curl http://localhost:9090/api/users

# POST 请求
curl -X POST http://localhost:9090/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'

# 带参数的请求
curl http://localhost:9090/api/user/123
```

## 项目结构

```
mock/
├── main.go                  # 主程序入口
├── config/
│   └── config.go           # 配置加载
├── handler/
│   ├── mock.go             # 请求处理器
│   └── dynamic.go          # 动态路由处理器
├── router/
│   └── router.go           # 路由注册
├── watcher/
│   └── watcher.go          # 文件监听器
├── reload/
│   └── manager.go          # 热加载管理器
├── configs/
│   ├── system.yml          # 系统配置
│   └── endpoints/          # 接口配置目录
├── go.mod
└── README.md
```

## 注意事项

1. 系统配置使用 YAML 格式（支持注释）
2. 接口配置使用 JSON 格式（API 标准格式）
3. 路径参数使用 `:` 前缀，如 `/api/user/:id`
4. 启用热加载后，修改配置文件自动生效，无需重启

## 对比其他工具

| 工具 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| **Postman Mock** | 界面友好，与 Postman 集成 | 需要账号，免费版有限制 | 个人/小团队 |
| **json-server** | 极简，30秒启动 | 只支持 CRUD，不支持自定义路由 | 快速原型 |
| **WireMock** | 功能强大，支持请求匹配、故障注入 | Java 为主，配置复杂 | 复杂测试场景 |
| **MSW** | 拦截网络请求，无需服务器 | 主要用于前端 | 前端开发 |
| **Apifox** | 中文友好，集成文档+Mock+测试 | 需要 GUI，商业版收费 | 国内团队 |
| **MockMate (本项目)** | 单文件部署，热加载，配置简单 | 功能较基础 | 内网 Mock，快速开发 |

## 路线图

本项目计划发展为更强大的 Mock 工具，以下是几个可能的方向：

### 🎯 方案 1：极致简单 - 命令行版 json-server

**定位**：5秒启动，零学习成本

- ✅ 单个可执行文件，无需安装依赖
- ✅ 配置即服务，改完自动生效（已实现热加载）
- ✅ 支持 REST 规范自动生成 CRUD
- 🔄 添加 CLI 初始化命令：`mockmate init`
- 🔄 添加自动生成示例数据功能

**适用场景**：个人开发者、快速原型

---

### 🌐 方案 2：团队协作 - 轻量级 Apifox

**定位**：给小团队的内网 Mock 工具

- ✅ 当前基础已就绪
- 🔄 添加 **Web UI 管理界面**（可视化编辑配置）
- 🔄 支持**多人协作**（配置文件放共享目录/数据库）
- 🔄 添加**导入/导出**功能（Postman/Swagger 格式）
- 🔄 添加**请求日志**和**统计面板**
- 🔄 支持**环境管理**（dev/test/prod）

**适用场景**：小团队协作、前后端联调

---

### 🤖 方案 3：智能 Mock - AI 驱动

**定位**：根据 OpenAPI 规范自动生成 Mock 数据

- 🔄 读取 **Swagger/OpenAPI 文档**自动生成接口
- 🔄 **智能数据生成**：
  - 根据字段类型生成符合语义的数据（邮箱、日期、正则等）
  - 支持中文姓名、手机号、地址等本地化数据
  - 支持自定义数据生成规则
- 🔄 支持**场景管理**：
  - 登录成功/失败/超时
  - 数据为空/单条/多条
  - 自定义响应逻辑
- 🔄 支持**请求匹配**：
  - 根据请求参数返回不同数据
  - 根据请求头、Cookie 等条件响应

**适用场景**：规范化团队、自动化测试

---

### 🔀 方案 4：混合模式 - Mock + Proxy

**定位**：部分接口 Mock，部分转发到真实服务

- 🔄 添加 `proxy_pass` 配置项
- 🔄 支持灵活的路由规则：
  ```yaml
  # Mock 接口
  - path: /api/user/*
    mock: true

  # 转发到测试环境
  - path: /api/payment/*
    proxy: https://test-api.example.com
  ```
- 🔄 支持**请求/响应修改**（添加/删除/修改 headers、body）
- 🔄 **抓包调试**：查看所有请求日志
- 🔄 支持**请求重放**

**适用场景**：本地开发连接测试环境、部分接口联调

---

## 你想要哪个方向？

欢迎提 Issue 或 PR 告诉我们你最需要的功能！

**当前优先级**：
1. ⭐ Web UI 管理界面
2. ⭐ 请求日志查看
3. ⭐ OpenAPI 导入功能

---

## License

MIT
