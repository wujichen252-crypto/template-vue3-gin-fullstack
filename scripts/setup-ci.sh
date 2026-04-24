#!/bin/bash
# 本地初始化辅助脚本，引导配置 Secrets 和 Variables

echo "=========================================="
echo "  Vue3 + Gin 模板 CI 配置向导"
echo "=========================================="
echo ""

read -p "请输入服务器 IP: " HOST
read -p "请输入 SSH 用户名 (默认: root): " USER
USER=${USER:-root}
read -p "请输入 SSH 端口 (默认: 22): " PORT
PORT=${PORT:-22}

echo ""
read -p "请输入前端部署路径 (默认: /www/wwwroot/template-app): " FRONTEND_PATH
FRONTEND_PATH=${FRONTEND_PATH:-/www/wwwroot/template-app}
read -p "请输入后端部署路径 (默认: /www/template-app): " BACKEND_PATH
BACKEND_PATH=${BACKEND_PATH:-/www/template-app}
read -p "请输入服务名称 (默认: template-app): " SERVICE_NAME
SERVICE_NAME=${SERVICE_NAME:-template-app}

echo ""
echo "=========================================="
echo "  请在 GitHub 仓库页面配置以下内容："
echo "=========================================="
echo ""
echo "【Settings → Secrets and variables → Actions → Secrets】"
echo ""
echo "  SSH_PRIVATE_KEY"
echo "    获取方式: 在服务器执行 cat ~/.ssh/id_rsa"
echo "    粘贴完整的私钥内容（包含 BEGIN/END 行）"
echo ""
echo "【Settings → Secrets and variables → Actions → Variables】"
echo ""
echo "  必填项:"
echo "    DEPLOY_SERVER_HOST    = $HOST"
echo "    DEPLOY_FRONTEND_PATH  = $FRONTEND_PATH"
echo "    DEPLOY_BACKEND_PATH   = $BACKEND_PATH"
echo ""
echo "  可选项:"
echo "    DEPLOY_SERVER_USER    = $USER"
echo "    DEPLOY_SERVER_PORT    = $PORT"
echo "    DEPLOY_SERVICE_NAME   = $SERVICE_NAME"
echo "    PROJECT_FRONTEND_DIR  = ./frontend"
echo "    PROJECT_BACKEND_DIR   = ./backend"
echo "    PROJECT_BACKEND_ENTRY = ./cmd/main.go"
echo "    NODE_VERSION          = 20"
echo "    GO_VERSION            = 1.22"
echo "    DEPLOY_NOTIFY_WEBHOOK = (部署通知 webhook 地址)"
echo ""
echo "=========================================="
echo "  服务器准备步骤："
echo "=========================================="
echo ""
echo "1. 生成 SSH 密钥对（如尚未生成）:"
echo "   ssh-keygen -t rsa -b 4096 -C \"github-actions\" -f ~/.ssh/github_actions"
echo ""
echo "2. 添加公钥到 authorized_keys:"
echo "   cat ~/.ssh/github_actions.pub >> ~/.ssh/authorized_keys"
echo ""
echo "3. 创建 systemd 服务文件 /etc/systemd/system/$SERVICE_NAME.service:"
echo ""
cat <<EOF
[Unit]
Description=$SERVICE_NAME API
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$BACKEND_PATH
ExecStart=$BACKEND_PATH/app
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
echo ""
echo "4. 启用并启动服务:"
echo "   sudo systemctl daemon-reload"
echo "   sudo systemctl enable $SERVICE_NAME"
echo "   sudo systemctl start $SERVICE_NAME"
echo ""
echo "=========================================="
echo "  配置完成后，执行 git push 即可触发部署"
echo "=========================================="
