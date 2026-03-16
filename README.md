<div align="center">

# GRA - Go React Admin

基于 **Gin + GORM** 的后台管理系统后端

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-1.12-blue?style=flat-square)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/GORM-1.31-red?style=flat-square)](https://gorm.io/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](./LICENSE)

</div>

---

## 特性

- **Clean Architecture** — Handler / Service / Repository 三层分离，职责清晰
- **域隔离架构** — `system`（系统管理）与 `business`（业务逻辑）完全分离
- **Provider 模式** — 依赖注入集中管理，main.go 永远只需 3 行
- **接口隔离** — 跨域调用通过接口解耦，可 mock 测试
- **路由自注册** — 每个模块自带 `RegisterRoutes`，路由文件永不膨胀
- **JWT 认证** — Bearer Token 鉴权，开箱即用
- **统一响应** — 标准化 JSON 响应格式，内置分页支持
- **软删除** — GORM 软删除，数据安全可追溯

---

## 技术栈

| 组件 | 技术 | 说明 |
|:-----|:-----|:-----|
| Web 框架 | Gin | 高性能 HTTP 框架 |
| ORM | GORM | MySQL 驱动 |
| 配置管理 | Viper | YAML 配置加载 |
| 日志 | Zap | 结构化高性能日志 |
| 认证 | golang-jwt | JWT Token 签发与验证 |
| 加密 | bcrypt | 密码哈希 |

---

## 项目结构

```
gra/
├── cmd/
│   └── server/
│       └── main.go                 # 程序入口
│
├── internal/                       # 内部业务代码（不可被外部引用）
│   ├── system/                     # 🔧 系统管理域
│   │   ├── provider.go             #    域级依赖注入入口
│   │   └── user/                   #    用户管理模块
│   │       ├── model.go            #      模型 + DTO
│   │       ├── repository.go       #      数据访问层
│   │       ├── service.go          #      业务逻辑层
│   │       ├── handler.go          #      HTTP 处理 + 路由注册
│   │       └── provider.go         #      模块级依赖注入
│   │
│   ├── business/                   # 💼 业务逻辑域
│   │   └── provider.go             #    域级依赖注入入口
│   │
│   ├── middleware/                  # 中间件
│   │   ├── cors.go                 #    跨域处理
│   │   └── jwt.go                  #    JWT 认证
│   │
│   └── router/                     # 路由注册
│       ├── router.go               #    总入口
│       ├── system.go               #    系统域路由
│       └── business.go             #    业务域路由
│
├── pkg/                            # 公共基础库（可被外部引用）
│   ├── config/config.go            #    配置加载
│   ├── database/database.go        #    数据库初始化
│   ├── logger/logger.go            #    日志初始化
│   └── response/response.go        #    统一响应封装
│
├── config/
│   └── config.yaml                 # 配置文件
│
├── go.mod
└── go.sum
```

---

## 快速开始

### 环境要求

- Go 1.25+
- MySQL 5.7+

### 1. 克隆项目

```bash
git clone <your-repo-url>
cd gra
```

### 2. 修改配置

编辑 `config/config.yaml`，配置数据库连接信息：

```yaml
server:
  port: 8888
  mode: debug

database:
  driver: mysql
  host: 127.0.0.1
  port: 3306
  username: root
  password: "your-password"
  dbname: gra
  charset: utf8mb4
```

### 3. 启动服务

```bash
go run cmd/server/main.go
```

服务启动后访问 `http://localhost:8888`

---

## API 接口

### 公开接口

| 方法 | 路径 | 说明 |
|:-----|:-----|:-----|
| POST | `/api/login` | 用户登录，获取 Token |

### 用户管理（需认证）

> 请求头添加：`Authorization: Bearer <token>`

| 方法 | 路径 | 说明 |
|:-----|:-----|:-----|
| POST | `/api/users` | 创建用户 |
| GET | `/api/users` | 用户列表（分页） |
| GET | `/api/users/:id` | 用户详情 |
| PUT | `/api/users/:id` | 更新用户 |
| DELETE | `/api/users/:id` | 删除用户 |

### 请求示例

**登录**

```bash
curl -X POST http://localhost:8888/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "123456"}'
```

**分页查询**

```bash
curl "http://localhost:8888/api/users?page=1&size=10" \
  -H "Authorization: Bearer <token>"
```

### 响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

分页响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "size": 10
  }
}
```

---

## 架构设计

### 依赖注入链路

```
main.go
  │
  │  sysHandlers, sysSvc := system.Init(db)     // 系统域初始化
  │  bizHandlers := business.Init(db, sysSvc)    // 业务域初始化（接收系统域服务）
  │  router.Setup(r, sysHandlers, bizHandlers)   // 路由注册
  │
  ├── system/provider.go       ← 聚合系统域所有模块
  │   └── user/provider.go     ← 模块内部自接线
  │
  └── business/provider.go     ← 聚合业务域，接收跨域依赖
      └── order/provider.go    ← 模块内部自接线
```

### 跨域调用（接口隔离）

```
  system/user             business/order
  ┌──────────┐            ┌──────────────┐
  │ Service  │───注入────▶│ UserQuerier  │  （消费方定义接口）
  │ (实现方) │            │ (消费方)      │
  └──────────┘            └──────────────┘
```

- 消费方定义接口，只声明需要的方法
- 提供方无需感知（Go 隐式接口）
- Provider 层负责接线

### 新增模块步骤

**新增系统模块**（如角色管理）：

1. 创建 `internal/system/role/` 目录，编写 `model.go`、`repository.go`、`service.go`、`handler.go`、`provider.go`
2. `handler.go` 中实现 `RegisterRoutes(r *gin.RouterGroup)`
3. `internal/system/provider.go` 的 `Handlers` 加字段，`Init()` 加一行
4. `internal/router/system.go` 加一行 `h.Role.RegisterRoutes(auth)`

**新增业务模块**（如订单管理）：

1. 创建 `internal/business/order/` 目录，编写四件套 + `provider.go`
2. `handler.go` 中实现 `RegisterRoutes(r *gin.RouterGroup)`
3. `internal/business/provider.go` 的 `Handlers` 加字段，`Init()` 加一行
4. `internal/router/business.go` 加一行 `h.Order.RegisterRoutes(auth)`

> main.go 无需修改。系统域与业务域互不干扰。

---

## 配置说明

| 配置项 | 说明 | 默认值 |
|:-------|:-----|:-------|
| `server.port` | 服务端口 | 8888 |
| `server.mode` | 运行模式（debug/release/test） | debug |
| `database.driver` | 数据库驱动 | mysql |
| `database.max_idle_conns` | 最大空闲连接数 | 10 |
| `database.max_open_conns` | 最大打开连接数 | 100 |
| `jwt.secret` | JWT 签名密钥 | - |
| `jwt.expire` | Token 过期时间（秒） | 7200 |
| `log.level` | 日志级别（debug/info/warn/error） | info |
| `log.format` | 日志格式（console/json） | console |

---

## License

MIT
