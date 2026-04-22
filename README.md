# template-vue3-gin-fullstack

前后端分离的全栈项目模板，基于 Vue3 + Go Gin 实现。

## 技术栈

### 前端
- Vue 3.4 + TypeScript 5
- Vite 5
- Pinia 2 (状态管理)
- Vue Router 4
- Axios (HTTP 客户端)
- Tailwind CSS 3
- Element Plus

### 后端
- Go 1.22
- Gin 1.9
- GORM 1.25
- PostgreSQL
- Redis
- JWT (golang-jwt/jwt/v5)
- Zap (日志)

## 目录结构

```
template-vue3-gin-fullstack/
├── frontend/                 # Vue3 前端
│   ├── src/
│   │   ├── api/             # API 接口
│   │   ├── components/      # 组件
│   │   ├── composables/     # 组合式函数
│   │   ├── router/          # 路由
│   │   ├── stores/          # Pinia 状态
│   │   ├── types/           # 类型定义
│   │   ├── utils/           # 工具函数
│   │   └── views/           # 页面
│   └── ...
├── backend/                  # Go 后端
│   ├── cmd/                 # 入口
│   ├── config/              # 配置
│   ├── internal/             # 内部包
│   │   ├── handler/         # 处理器
│   │   ├── service/         # 业务逻辑
│   │   ├── repository/      # 数据访问
│   │   ├── model/           # 数据模型
│   │   └── middleware/      # 中间件
│   ├── pkg/                 # 公共包
│   └── ...
└── ...
```

## 快速开始

### 前置条件

- Node.js >= 18
- Go >= 1.22
- PostgreSQL >= 14
- Redis >= 6

### 1. 克隆项目

```bash
git clone <repository-url>
cd template-vue3-gin-fullstack
```

### 2. 启动后端

```bash
cd backend

cp .env.example .env

# 编辑 .env 文件，配置数据库和 Redis 连接信息

go mod tidy

go run cmd/main.go
```

后端启动后运行于 `http://localhost:8080`

### 3. 启动前端

```bash
cd frontend

npm install

npm run dev
```

前端启动后访问 `http://localhost:5173`

### 4. 初始化数据库

```bash
cd backend/scripts

psql -U postgres -f init_db.sql
```

## 环境变量

### 后端 (.env)

| 变量 | 说明 | 默认值 |
|------|------|--------|
| SERVER_PORT | 服务端口 | 8080 |
| SERVER_MODE | 运行模式 | debug |
| ALLOW_ORIGINS | 允许的跨域来源（逗号分隔） | http://localhost:5173,http://localhost:3000 |
| DB_HOST | 数据库地址 | localhost |
| DB_PORT | 数据库端口 | 5432 |
| DB_USER | 数据库用户 | postgres |
| DB_PASSWORD | 数据库密码 | - |
| DB_NAME | 数据库名 | template_db |
| DB_SSLMODE | SSL模式 | disable |
| REDIS_HOST | Redis地址 | localhost |
| REDIS_PORT | Redis端口 | 6379 |
| REDIS_PASSWORD | Redis密码 | - |
| REDIS_DB | Redis数据库 | 0 |
| JWT_SECRET | JWT密钥 | - |
| JWT_ACCESS_EXPIRE | Access Token有效期(小时) | 2 |
| JWT_REFRESH_EXPIRE | Refresh Token有效期(小时) | 168 |

### 前端

| 变量 | 说明 |
|------|------|
| VITE_API_BASE_URL | API 基础路径 |

## API 接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/v1/auth/register | 注册 | 否 |
| POST | /api/v1/auth/login | 登录 | 否 |
| GET | /api/v1/auth/userinfo | 获取用户信息 | 是 |
| POST | /api/v1/auth/refresh | 刷新Token | 是 |

## 前端构建

```bash
cd frontend

npm run build
```

构建产物在 `frontend/dist/` 目录。

## 后端构建

```bash
cd backend

# Linux 交叉编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-server ./cmd/main.go

# Windows
go build -o app-server.exe ./cmd/main.go
```

## 部署

### Docker

```bash
cd backend

docker build -t template-backend .
docker run -p 8080:8080 --env-file .env template-backend
```

### Nginx 配置

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        root /path/to/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Swagger 文档

开发模式下访问：`http://localhost:8080/swagger/index.html`

## License

MIT
