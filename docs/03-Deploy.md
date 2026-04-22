# 宝塔部署手册

## 服务器环境要求

- 系统: CentOS 7+ / Ubuntu 20.04+
- 内存: >= 2GB
- 磁盘: >= 20GB
- Nginx >= 1.18
- MySQL >= 5.7 / PostgreSQL >= 13
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
- Nginx 1.18
- PostgreSQL 13
- Redis 6
- Node.js 18

## 2. 数据库配置

### 2.1 创建数据库
```sql
CREATE DATABASE template_db;
```

### 2.2 初始化表结构
```bash
cd /www/wwwroot/template/backend/scripts
psql -U postgres -d template_db -f init_db.sql
```

## 3. 后端部署

### 3.1 上传代码
```bash
cd /www/wwwroot/template
git clone <repository-url> .
cd backend
```

### 3.2 配置环境变量
```bash
cp .env.example .env
nano .env
```

修改以下配置:
```
SERVER_PORT=8080
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=template_db
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
JWT_SECRET=your-256-bit-secret
```

### 3.3 安装依赖和构建
```bash
go mod tidy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-server ./cmd/main.go
```

### 3.4 创建systemd服务
```bash
nano /etc/systemd/system/template.service
```

```ini
[Unit]
Description=Template Backend
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/www/wwwroot/template/backend
ExecStart=/www/wwwroot/template/backend/app-server
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
systemctl daemon-reload
systemctl enable template
systemctl start template
```

## 4. 前端部署

### 4.1 构建前端
```bash
cd /www/wwwroot/template/frontend
npm install
npm run build
```

### 4.2 构建产物位置
`/www/wwwroot/template/frontend/dist`

## 5. Nginx 配置

### 5.1 创建站点
在宝塔面板中创建站点，域名指向你的域名。

### 5.2 配置Nginx
```nginx
server {
    listen 80;
    server_name your-domain.com;

    root /www/wwwroot/template/frontend/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /swagger {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
    }
}
```

### 5.3 SSL配置
在宝塔面板中为站点申请Let's Encrypt证书。

## 6. 验证部署

### 6.1 检查后端
```bash
curl http://localhost:8080/
```

### 6.2 检查前端
访问 http://your-domain.com

### 6.3 检查API
```bash
curl http://localhost:8080/api/v1/auth/login -X POST -H "Content-Type: application/json" -d '{"username":"test","password":"test"}'
```

## 7. 常见问题

### 7.1 后端无法启动
检查:
- 端口是否被占用
- 数据库连接是否正常
- .env 配置是否正确

### 7.2 前端页面空白
检查:
- Nginx 配置是否正确
- 静态文件路径是否正确
- 路由是否正确配置

### 7.3 API 请求失败
检查:
- Nginx 代理是否配置正确
- 后端服务是否运行
- CORS 配置是否正确
