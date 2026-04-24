# 可配置模板仓库 CI/CD 设计指南

> **文档编号**：TECH-DEV-003  
> **归档日期**：2026-04-23  
> **关联项目**：NeuraMind 智脑笔记系统（Vue3 + Gin 全栈模板）  
> **文档类型**：工程化实践与模板设计规范

---

## 目录

1. [设计目标与核心原则](#一设计目标与核心原则)
2. [模板仓库目录结构](#二模板仓库目录结构)
3. [可配置 CI 架构设计](#三可配置-ci-架构设计)
4. [配置文件详解](#四配置文件详解)
5. [初始化与配置流程](#五初始化与配置流程)
6. [多环境支持方案](#六多环境支持方案)
7. [安全与权限设计](#七安全与权限设计)
8. [使用手册（给模板使用者）](#八使用手册给模板使用者)
9. [完整文件参考](#九完整文件参考)

---

## 一、设计目标与核心原则

### 1.1 模板仓库的特殊性

模板仓库（Template Repository）被用于 **"Use this template"** 生成新仓库时，以下资源**不会**自动继承：

- **Secrets**：所有加密凭据（SSH 密钥、API Token）全部清空
- **Variables**：仓库变量（非敏感配置）全部清空
- **Actions 权限**：首次需要手动授权 Workflow 写入权限
- **Runner 配置**：自托管 Runner 不会跟随模板

### 1.2 设计目标

让从模板创建的新仓库，在 **5 分钟内** 完成 CI 配置并首次部署成功，无需修改 YAML 文件本身。

### 1.3 核心原则

| 原则           | 说明                                    | 实现方式                                               |
| -------------- | --------------------------------------- | ------------------------------------------------------ |
| **零硬编码**   | YAML 中不出现 IP、路径、域名等具体值    | 全部使用 `${{ vars.XXX }}` 和 `${{ secrets.XXX }}`     |
| **防御性编程** | 未配置完成时，CI 优雅跳过而非报错失败   | `if: vars.XXX != ''` 条件判断                          |
| **自描述配置** | 使用者一看 Variables 列表就知道要填什么 | 命名规范：`DEPLOY_SERVER_HOST`、`DEPLOY_FRONTEND_PATH` |
| **多环境就绪** | 支持开发/生产多环境，通过分支或变量切换 | `vars.DEPLOY_ENV` 控制部署目标                         |
| **技术栈无关** | 虽然基于 Vue3+Gin，但目录名、端口可配置 | `vars.PROJECT_FRONTEND_DIR`                            |

---

## 二、模板仓库目录结构

```
neuramind-template/                 ← 模板仓库根目录
├── .github/
│   └── workflows/
│       ├── ci-check.yml            ← 代码质量检查（无部署，开箱即用）
│       └── deploy.yml              ← 可配置部署流水线（未配置时自动跳过）
├── scripts/
│   └── setup-ci.sh               ← 本地初始化脚本，引导配置 Secrets/Vars
├── frontend/                       ← Vue3 前端（目录名可配置）
├── backend/                        ← Gin 后端（目录名可配置）
├── .env.example                    ← 环境变量模板
├── README.md                       ← 必须包含"从模板创建后必做"章节
└── LICENSE
```

---

## 三、可配置 CI 架构设计

### 3.1 配置分层模型

```
┌─────────────────────────────────────────────┐
│  层 1：代码层（YAML 文件）                     │
│  └── 只定义流程逻辑，不出现任何具体值             │
│       │                                       │
│  层 2：仓库变量层（Settings → Variables）      │
│  └── 非敏感配置：IP、路径、端口、环境标识         │
│       │                                       │
│  层 3：仓库密钥层（Settings → Secrets）        │
│  └── 敏感信息：SSH 私钥、数据库密码、API Token   │
└─────────────────────────────────────────────┘
```

### 3.2 命名规范

| 类型           | 前缀           | 示例                   | 存放位置  |
| -------------- | -------------- | ---------------------- | --------- |
| 部署服务器相关 | `DEPLOY_`      | `DEPLOY_SERVER_HOST`   | Variables |
| 项目路径相关   | `PROJECT_`     | `PROJECT_FRONTEND_DIR` | Variables |
| 运行时配置     | `APP_`         | `APP_API_PORT`         | Variables |
| 敏感凭据       | `SSH_` / `DB_` | `SSH_PRIVATE_KEY`      | Secrets   |

---

## 四、配置文件详解

### 4.1 代码质量检查（ci-check.yml）

**特点**：无需任何配置，从模板创建后立即生效。

```yaml
name: CI Check

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  frontend-check:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: ${{ vars.PROJECT_FRONTEND_DIR || './frontend' }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: ${{ vars.NODE_VERSION || '20' }}
          cache: 'npm'
          cache-dependency-path: '${{ vars.PROJECT_FRONTEND_DIR || ''./frontend'' }}/package-lock.json'
      - run: npm ci
      - run: npm run lint
      - run: npm run type-check || true   # 如果有 tsc 检查命令
      - run: npm run build

  backend-check:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: ${{ vars.PROJECT_BACKEND_DIR || './backend' }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ vars.GO_VERSION || '1.22' }}
      - run: go mod download
      - run: go test -v ./...
      - run: |
          CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
          go build -ldflags="-s -w" -o app ${{ vars.PROJECT_BACKEND_ENTRY || './cmd/server' }}
```

**设计要点**：

- 所有路径使用 `${{ vars.XXX || '默认值' }}`，未配置时使用模板默认结构
- 不依赖任何 Secrets，新仓库创建后立即能跑通

---

### 4.2 可配置部署流水线（deploy.yml）

**特点**：未配置完成时，Job 自动跳过，不会报错。

```yaml
name: Deploy

on:
  push:
    branches: [main]
  workflow_dispatch:
    inputs:
      environment:
        description: '部署环境'
        required: true
        default: 'production'
        type: choice
        options:
          - production
          - staging

env:
  # 从 Variables 读取，提供默认值
  FRONTEND_DIR: ${{ vars.PROJECT_FRONTEND_DIR || './frontend' }}
  BACKEND_DIR: ${{ vars.PROJECT_BACKEND_DIR || './backend' }}
  BACKEND_ENTRY: ${{ vars.PROJECT_BACKEND_ENTRY || './cmd/server' }}
  NODE_VER: ${{ vars.NODE_VERSION || '20' }}
  GO_VER: ${{ vars.GO_VERSION || '1.22' }}

jobs:
  # ==========================================
  # Job 1: 构建前端
  # ==========================================
  build-frontend:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: ${{ env.FRONTEND_DIR }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VER }}
          cache: 'npm'
          cache-dependency-path: '${{ env.FRONTEND_DIR }}/package-lock.json'
      - run: npm ci
      - run: npm run build
      - uses: actions/upload-artifact@v4
        with:
          name: frontend-dist
          path: ${{ env.FRONTEND_DIR }}/dist
          retention-days: 3

  # ==========================================
  # Job 2: 构建后端
  # ==========================================
  build-backend:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: ${{ env.BACKEND_DIR }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VER }}
      - run: go mod download
      - run: go test -v ./...
      - run: |
          CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
          go build -ldflags="-s -w" -o app ${{ env.BACKEND_ENTRY }}
      - uses: actions/upload-artifact@v4
        with:
          name: backend-binary
          path: ${{ env.BACKEND_DIR }}/app
          retention-days: 3

  # ==========================================
  # Job 3: 部署（带前置条件检查）
  # ==========================================
  deploy:
    needs: [build-frontend, build-backend]
    runs-on: ubuntu-latest
    
    # 防御性条件：关键变量未配置时直接跳过
    if: |
      vars.DEPLOY_SERVER_HOST != '' &&
      vars.DEPLOY_FRONTEND_PATH != '' &&
      vars.DEPLOY_BACKEND_PATH != '' &&
      secrets.SSH_PRIVATE_KEY != ''
    
    steps:
      - name: Download frontend artifact
        uses: actions/download-artifact@v4
        with:
          name: frontend-dist
          path: ./dist
      
      - name: Download backend artifact
        uses: actions/download-artifact@v4
        with:
          name: backend-binary
          path: ./api

      # 部署前端静态文件
      - name: Deploy Frontend
        uses: appleboy/scp-action@v0.1.7
        with:
          host: ${{ vars.DEPLOY_SERVER_HOST }}
          username: ${{ vars.DEPLOY_SERVER_USER || 'root' }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          port: ${{ vars.DEPLOY_SERVER_PORT || '22' }}
          source: "dist/*"
          target: ${{ vars.DEPLOY_FRONTEND_PATH }}
          strip_components: 1

      # 部署后端 + 重启服务
      - name: Deploy Backend
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ vars.DEPLOY_SERVER_HOST }}
          username: ${{ vars.DEPLOY_SERVER_USER || 'root' }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          port: ${{ vars.DEPLOY_SERVER_PORT || '22' }}
          script: |
            set -e
            
            BACKEND_PATH="${{ vars.DEPLOY_BACKEND_PATH }}"
            SERVICE_NAME="${{ vars.DEPLOY_SERVICE_NAME || 'neuramind' }}"
            
            # 备份旧版本
            if [ -f "$BACKEND_PATH/app" ]; then
              cp "$BACKEND_PATH/app" "$BACKEND_PATH/app.bak.$(date +%s)"
            fi
            
            # 确保目录存在
            mkdir -p "$BACKEND_PATH"
            
            # 移动新二进制（SCP 默认传到用户 home，再移到目标位置）
            mv /home/${{ vars.DEPLOY_SERVER_USER || 'root' }}/app "$BACKEND_PATH/app" 2>/dev/null || \
            mv ~/app "$BACKEND_PATH/app" 2>/dev/null || \
            mv /tmp/app "$BACKEND_PATH/app" 2>/dev/null || true
            
            chmod +x "$BACKEND_PATH/app"
            
            # 重启服务（systemd 优先，fallback supervisor）
            if command -v systemctl &> /dev/null && systemctl list-unit-files | grep -q "$SERVICE_NAME"; then
              sudo systemctl restart "$SERVICE_NAME"
              sleep 3
              systemctl is-active --quiet "$SERVICE_NAME" || {
                echo "服务启动失败，尝试回滚..."
                cp "$BACKEND_PATH/app.bak."* "$BACKEND_PATH/app"
                sudo systemctl restart "$SERVICE_NAME"
                exit 1
              }
            elif command -v supervisorctl &> /dev/null; then
              sudo supervisorctl restart "$SERVICE_NAME"
            else
              echo "未找到服务管理工具，请手动重启"
              exit 1
            fi
            
            echo "部署成功：$(date)"

      # 可选：部署通知
      - name: Notify
        if: vars.DEPLOY_NOTIFY_WEBHOOK != ''
        run: |
          curl -X POST ${{ vars.DEPLOY_NOTIFY_WEBHOOK }} \
            -H 'Content-Type: application/json' \
            -d '{"msg":"部署成功: ${{ github.repository }}@${{ github.sha }}"}' || true
```

---

### 4.3 配置清单（README 中必须提供）

从模板创建新仓库后，在 **Settings → Secrets and variables → Actions** 中配置：

#### Secrets（加密）

| 名称              | 必填 | 获取方式                                 |
| ----------------- | ---- | ---------------------------------------- |
| `SSH_PRIVATE_KEY` | ✅    | 服务器执行 `cat ~/.ssh/id_rsa`，粘贴全文 |

#### Variables（明文）

| 名称                    | 必填 | 示例                            | 说明                         |
| ----------------------- | ---- | ------------------------------- | ---------------------------- |
| `DEPLOY_SERVER_HOST`    | ✅    | `123.456.78.90`                 | 服务器公网 IP                |
| `DEPLOY_SERVER_USER`    | ❌    | `root`                          | SSH 用户名，默认 root        |
| `DEPLOY_SERVER_PORT`    | ❌    | `22`                            | SSH 端口，默认 22            |
| `DEPLOY_FRONTEND_PATH`  | ✅    | `/www/wwwroot/neuramind`        | 前端部署目录（宝塔站点路径） |
| `DEPLOY_BACKEND_PATH`   | ✅    | `/www/neuramind`                | 后端二进制存放目录           |
| `DEPLOY_SERVICE_NAME`   | ❌    | `neuramind`                     | systemd/supervisor 服务名    |
| `PROJECT_FRONTEND_DIR`  | ❌    | `./frontend`                    | 前端代码相对路径             |
| `PROJECT_BACKEND_DIR`   | ❌    | `./backend`                     | 后端代码相对路径             |
| `PROJECT_BACKEND_ENTRY` | ❌    | `./cmd/server`                  | Go 主程序入口                |
| `NODE_VERSION`          | ❌    | `20`                            | Node.js 版本                 |
| `GO_VERSION`            | ❌    | `1.22`                          | Go 版本                      |
| `DEPLOY_NOTIFY_WEBHOOK` | ❌    | `https://oapi.dingtalk.com/...` | 可选，部署成功通知地址       |

---

## 五、初始化与配置流程

### 5.1 模板使用者操作步骤

```bash
# 1. 在 GitHub 点击 "Use this template" 创建新仓库

# 2. 克隆新仓库到本地
git clone https://github.com/你的用户名/新项目名.git
cd 新项目名

# 3. 按 README 配置 Secrets 和 Variables（网页操作）

# 4. 服务器上准备 SSH 公钥授权
ssh-keygen -t rsa -b 4096 -C "github-actions" -f ~/.ssh/github_actions
cat ~/.ssh/github_actions.pub >> ~/.ssh/authorized_keys
# 把私钥填到 GitHub Secrets：SSH_PRIVATE_KEY

# 5. 服务器上准备服务管理（systemd）
sudo tee /etc/systemd/system/neuramind.service > /dev/null <<'EOF'
[Unit]
Description=NeuraMind API
After=network.target
[Service]
Type=simple
User=root
WorkingDirectory=/www/neuramind
ExecStart=/www/neuramind/app
Restart=always
[Install]
WantedBy=multi-user.target
EOF
sudo systemctl daemon-reload && sudo systemctl enable neuramind

# 6. 本地首次提交，触发 CI
git add .
git commit -m "init: 从模板初始化项目"
git push origin main

# 7. 去 GitHub → Actions 查看流水线状态
```

---

## 六、多环境支持方案

通过 **分支策略** 或 **手动触发参数** 实现：

### 方案 A：分支对应环境（推荐）

```yaml
on:
  push:
    branches:
      - main      # 自动部署到生产
      - develop   # 自动部署到测试
```

在 deploy Job 中根据分支选择路径：

```yaml
- name: Set environment variables
  run: |
    if [ "${{ github.ref_name }}" == "main" ]; then
      echo "TARGET_PATH=${{ vars.DEPLOY_FRONTEND_PATH }}" >> $GITHUB_ENV
    else
      echo "TARGET_PATH=${{ vars.DEPLOY_FRONTEND_PATH_STAGING || vars.DEPLOY_FRONTEND_PATH }}" >> $GITHUB_ENV
    fi
```

### 方案 B：手动触发选择

已在 `workflow_dispatch` 中配置 `environment` 参数，在 Actions 页面点击 **Run workflow** 时选择环境。

---

## 七、安全与权限设计

### 7.1 模板仓库本身的安全

- **不要在模板里放任何真实凭据**：包括 `.env`、测试密钥、个人服务器 IP
- **提供 `.env.example`**：只留空模板，标注必填项
- **Actions 权限设置**：模板仓库建议开启 **Settings → Actions → General**：
  - **Workflow permissions**：Read and write permissions（用于上传 artifact）
  - **Fork pull request workflows**：选择 **Require approval for first-time contributors**

### 7.2 新仓库创建后的安全加固

| 操作                  | 路径                         | 说明                                 |
| --------------------- | ---------------------------- | ------------------------------------ |
| 限制 Actions 读写权限 | Settings → Actions → General | 防止恶意 Workflow 篡改仓库           |
| 开启分支保护          | Settings → Branches          | `main` 分支要求 PR + CI 通过才能合并 |
| 开启 Secret scanning  | Settings → Security          | GitHub 自动检测是否误提交密钥        |

---

## 八、使用手册（给模板使用者）

将此段放入模板仓库的 **README.md**：

```markdown
## 🚀 从模板创建后必做（5 分钟完成 CI 配置）

### 1. 配置仓库变量
进入仓库 **Settings → Secrets and variables → Actions → Variables**（New repository variable），添加：

| 变量名 | 示例值 |
|--------|--------|
| `DEPLOY_SERVER_HOST` | `123.456.78.90` |
| `DEPLOY_FRONTEND_PATH` | `/www/wwwroot/your-project` |
| `DEPLOY_BACKEND_PATH` | `/www/your-project` |

### 2. 配置密钥
进入 **Secrets**，添加：

- `SSH_PRIVATE_KEY`：你的服务器私钥（`cat ~/.ssh/id_rsa`）

### 3. 服务器准备
确保服务器已添加公钥授权，且 systemd 服务已创建（见上方文档）。

### 4. 首次推送
```bash
git push origin main
```

然后去 **Actions** 标签页查看部署状态。

### 5. 自定义项目结构

如果前端目录不叫 `frontend`，在 Variables 中添加：

- `PROJECT_FRONTEND_DIR` = `./your-dir-name`

```
---

## 九、完整文件参考

### `scripts/setup-ci.sh`

```bash
#!/bin/bash
# 本地初始化辅助脚本，可选使用

echo "=========================================="
echo "  NeuraMind 模板 CI 配置向导"
echo "=========================================="
echo ""

read -p "请输入服务器 IP: " HOST
read -p "请输入前端部署路径 (默认: /www/wwwroot/neuramind): " FRONTEND_PATH
FRONTEND_PATH=${FRONTEND_PATH:-/www/wwwroot/neuramind}

read -p "请输入后端部署路径 (默认: /www/neuramind): " BACKEND_PATH
BACKEND_PATH=${BACKEND_PATH:-/www/neuramind}

echo ""
echo "请在 GitHub 仓库页面手动配置以下 Secrets 和 Variables："
echo ""
echo "【Secrets】"
echo "  SSH_PRIVATE_KEY"
echo ""
echo "【Variables】"
echo "  DEPLOY_SERVER_HOST = $HOST"
echo "  DEPLOY_FRONTEND_PATH = $FRONTEND_PATH"
echo "  DEPLOY_BACKEND_PATH = $BACKEND_PATH"
echo ""
echo "配置完成后，执行 git push 即可触发首次部署。"
echo "=========================================="
```

### `.env.example`

```bash
# 本地开发环境变量模板
# 复制为 .env 后填入实际值

# 后端配置
APP_PORT=8080
APP_MODE=release
DB_DSN=root:password@tcp(127.0.0.1:3306)/neuramind?charset=utf8mb4

# 前端配置
VITE_API_BASE_URL=/api
```

---

## 总结

| 你（模板作者）需要做的                                  | 使用者需要做的                    |
| ------------------------------------------------------- | --------------------------------- |
| 写好 `.github/workflows/*.yml`，全部用 `vars`/`secrets` | 在 GitHub 网页填入服务器 IP、路径 |
| 提供 `README.md` 配置清单                               | 在服务器生成 SSH 密钥对，公钥授权 |
| 提供 systemd 服务文件模板                               | 把私钥粘贴到 Secrets              |
| 确保 `if:` 条件防止未配置时失败                         | `git push` 触发部署               |

> **核心思想**：模板仓库的 CI 文件是"万能插头"，具体插到哪里（哪台服务器、哪个路径），由 Variables 和 Secrets 决定，绝不写死在代码里。

---

*归档人：AI 助手*  
*关联文档：TECH-DEV-002《GitHub Actions 工作原理与最佳实践》*  
*技术标签：`模板仓库` `GitHub Actions` `可配置CI` `Vue3` `Gin` `工程化`*