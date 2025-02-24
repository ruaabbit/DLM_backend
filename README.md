# 挡粮门记录后端服务

## 项目简介

这是一个基于 Go 语言和 Gin 框架开发的后端服务，用于支持“挡粮门记录”微信小程序。该服务提供用户登录、挡粮门点检信息记录、查询等功能，旨在通过高效、稳定的后端接口支持小程序的日常运行。

## 技术栈

- **后端框架**：Go Gin
- **数据库**：MySQL（可根据需求替换为其他数据库）
- **开发工具**：GoLand / VS Code
- **测试工具**：Postman / cURL
- **代码规范**：Go Code Review Comments

## 项目结构

```
DLM_backend/
├── cmd/
│ └── main.go # 项目入口文件
├── config/
│ └── config.go # 配置文件加载与解析
├── controllers/
│ ├── auth.go # 用户认证控制器
│ └── inspection.go # 挡粮门点检控制器
├── models/
│ ├── user.go # 用户模型
│ └── inspection_record.go # 点检记录模型
├── routers/
│ └── router.go # 路由配置
├── services/
│ ├── auth_service.go # 用户认证服务
│ └── inspection_service.go # 点检记录服务
├── utils/
│ ├── response.go # 统一响应格式
│ └── jwt.go # JWT 工具
├── README.md # 项目说明文档
└── go.mod # Go 模块依赖管理
```

## 功能模块

### 1. 用户认证

- **登录接口**：支持保管员和管理员登录，返回 JWT 令牌。
- **鉴权中间件**：基于 JWT 的权限校验，确保接口安全。

### 2. 挡粮门点检记录

- **新增记录**：允许保管员提交挡粮门点检信息，包括检查日期、检查项目、检查情况及图片。
- **查询记录**：支持按条件查询点检记录，如按日期、仓号、责任人等。
- **图片上传**：支持图片上传功能，存储至本地或云存储服务。

### 3. 管理员功能

- **用户管理**：管理员可以查看、添加、删除保管员账号。
- **记录管理**：管理员可以查看所有点检记录，支持导出功能。

## 环境配置

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 数据库配置

在 config/config.go 中配置数据库连接信息：

```go
type Config struct {
    DBHost     string `env:"DB_HOST"`
    DBPort     string `env:"DB_PORT"`
    DBUser     string `env:"DB_USER"`
    DBPassword string `env:"DB_PASSWORD"`
    DBName     string `env:"DB_NAME"`
    JWTSecret  string `env:"JWT_SECRET"`
}
```

### 3. 启动服务

```bash
go run cmd/main.go
```

默认服务运行在 localhost:8080。
