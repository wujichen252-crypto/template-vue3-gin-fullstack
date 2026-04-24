# 宝塔部署手册

## 服务器环境要求

- 系统: CentOS 7+ / Ubuntu 20.04+
- 内存: >= 2GB
- 磁盘: >= 20GB
- Nginx >= 1.18
- PostgreSQL >= 14
- Redis >= 6

## 1. 环境安装

### 1.1 安装宝塔面板
```bash
# CentOS
yum install -y wget && wget -O install.sh http://download.bt.cn/install/install_6.0.sh && sh install.sh

# Ubuntu
wget -O install.sh http://download.bt.cn/install/install-ubuntu_6.0.sh && sudo bash install.sh
```

### 1.2 安装软件
在宝塔面板中安装以下软件:
- Nginx 1.18+
- PostgreSQL 14+
- Redis 6+
- Node.js 18+

## 2. 数据库配置

### 2.1 创建数据库
```sql
CREATE DATABASE template_db;
```

### 2.2 初始化表结构
```bash
cd /www/wwwroot/template-vue3-gin-fullstack/backend/scripts
psql -U postgres -d template_db -f init_db.sql
```

## 3. 后端部署

### 3.1 上传代码
```bash
cd /www/wwwroot
git clone <repository-url> template-vue3-gin-fullstack
cd template-vue3-gin-fullstack/backend
```

### 3.2 配置环境变量
```bash
cp .env.example .env
nano .env
```

修改以下配置:
```
SERVER_PORT=8080
SERVER_MODE=release
ALLOW_ORIGINS=https://your-domain.com
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=template_db
DB_SSLMODE=disable
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
JWT_SECRET=your-256-bit-secret
JWT_ACCESS_EXPIRE=2
JWT_REFRESH_EXPIRE=168
```

### 3.3 安装依赖和构建
```bash
go mod tidy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go
```

### 3.4 创建systemd服务
```bash
nano /etc/systemd/system/template-vue3-gin-fullstack.service
```

```ini
[Unit]
Description=NeuraMind Backend API
After=network.target
Wants=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/www/wwwroot/template-vue3-gin-fullstack/backend
ExecStart=/www/wwwroot/template-vue3-gin-fullstack/backend/app
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

```bash
systemctl daemon-reload
systemctl enable template-vue3-gin-fullstack
systemctl start template-vue3-gin-fullstack
```

### 3.5 检查后端状态
```bash
systemctl status template-vue3-gin-fullstack
curl http://localhost:8080/health
```

## 4. 前端部署

### 4.1 构建前端
```bash
cd /www/wwwroot/template-vue3-gin-fullstack/frontend
npm install
npm run build
```

### 4.2 构建产物位置
`/www/wwwroot/template-vue3-gin-fullstack/frontend/dist`

## 5. Nginx 配置

### 5.1 创建站点
在宝塔面板中创建站点，域名指向你的域名。

### 5.2 配置Nginx
```nginx
server {
    listen 80;
    server_name your-domain.com;

    root /www/wwwroot/template-vue3-gin-fullstack/frontend/dist;
    index index.html;

    # 前端路由
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Swagger 文档
    location /swagger {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
    }

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2)$ {
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
```

### 5.3 SSL配置
在宝塔面板中为站点申请Let's Encrypt证书。

## 6. CI/CD 自动化部署（可选）

### 6.1 配置 GitHub Actions
本项目支持 GitHub Actions 自动部署，详见 [CI配置文件.md](../CI配置文件.md)。

需要在 GitHub 仓库中配置以下 Secrets 和 Variables:
- Secrets: `SSH_PRIVATE_KEY`
- Variables: `DEPLOY_SERVER_HOST`, `DEPLOY_FRONTEND_PATH`, `DEPLOY_BACKEND_PATH`

### 6.2 手动触发部署
```bash
git push origin main
```
然后在 GitHub Actions 页面查看部署状态。

## 7. 验证部署

### 7.1 检查后端
```bash
curl http://localhost:8080/health
```

### 7.2 检查前端
访问 https://your-domain.com

### 7.3 检查API
```bash
curl http://localhost:8080/api/v1/auth/login -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test"}'
```

## 8. 常见问题

### 8.1 后端无法启动
检查:
- 端口是否被占用: `lsof -i :8080`
- 数据库连接是否正常: `psql -U postgres -d template_db -c "SELECT 1"`
- .env 配置是否正确
- 服务日志: `journalctl -u template-vue3-gin-fullstack -f`

### 8.2 前端页面空白
检查:
- Nginx 配置是否正确
- 静态文件路径 `/www/wwwroot/template-vue3-gin-fullstack/frontend/dist` 是否存在
- 路由是否正确配置（try_files）
- 浏览器控制台是否有跨域错误

### 8.3 API 请求失败
检查:
- Nginx 代理配置是否正确
- 后端服务是否运行: `systemctl status template-vue3-gin-fullstack`
- CORS 配置是否允许前端域名
- API 路径是否正确 `/api/v1/...`

### 8.4 数据库迁移问题
如果需要重新初始化数据库:
```bash
cd /www/wwwroot/template-vue3-gin-fullstack/backend/scripts
psql -U postgres -d template_db -f init_db.sql
```
