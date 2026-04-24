# 项目开发流程

## 1. 开发准备

### 环境配置
- 安装 Node.js >= 18
- 安装 Go >= 1.22
- 安装 PostgreSQL >= 14
- 安装 Redis >= 6
- 安装 Git

### 项目初始化
```bash
# 克隆项目
git clone <repository-url>
cd template-vue3-gin-fullstack

# 前端依赖
cd frontend
npm install

# 后端依赖
cd ../backend
go mod tidy
```

## 2. 分支管理

### 分支策略（Git Flow 简化版）
```
main        # 主分支，生产环境代码，保护分支
├── develop # 开发分支（可选，视团队规模而定）
├── feature/*  # 功能分支，从 main 或 develop 创建
└── bugfix/*   # 修复分支，从 main 或 develop 创建
```

### 简化分支策略（单人/小团队）
对于简单项目或单人开发，可使用更简化的策略：
```
main        # 主分支，所有代码直接提交到 main
├── feature/*  # 功能分支（可选）
└── bugfix/*   # 修复分支（可选）
```

### 分支命名规范
- 功能分支: `feature/<功能描述>` (如: `feature/user-auth`)
- 修复分支: `bugfix/<问题描述>` (如: `bugfix/login-error`)
- 版本标签: `v<版本号>` (如: `v1.0.0`)

## 3. 开发流程

### 标准开发步骤
1. 从 main 创建功能分支: `git checkout -b feature/user-auth`
2. 编写代码（遵循各模块规范）
3. 本地测试
4. 提交代码（遵守提交规范）
5. 推送分支: `git push origin feature/user-auth`
6. 创建 Pull Request 或 Merge Request
7. 代码审查
8. 合并到 main

### 提交规范
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type 类型**:
| 类型 | 说明 |
|------|------|
| feat | 新功能 |
| fix | 修复 bug |
| docs | 文档更新 |
| style | 代码格式（不影响功能） |
| refactor | 重构（不影响功能） |
| test | 测试相关 |
| chore | 构建/工具相关 |

**示例**:
```
feat(auth): 添加用户注册功能

- 添加用户名唯一性检查
- 添加邮箱格式验证
- 密码使用 bcrypt 加密

Closes #123
```

## 4. 测试

### 前端测试
```bash
cd frontend

# 运行所有测试
npm run test

# 监听模式（开发时）
npm run test:watch

# 类型检查
npm run typecheck

# ESLint 检查
npm run lint

# 完整检查（推荐在 PR 前运行）
npm run typecheck && npm run lint
```

### 后端测试
```bash
cd backend

# 运行所有测试
go test ./...

# 运行测试并显示详情
go test -v ./...

# 运行测试并检查覆盖率
go test -cover ./...

# 构建检查
go build ./...
```

### 测试覆盖率
- 前端: Vitest 覆盖率报告
- 后端: `go test -coverprofile=coverage.out && go tool cover`

## 5. CI/CD

### GitHub Actions 工作流

本项目包含两个 CI/CD 工作流：

#### ci-check.yml - 代码检查
- 触发: PR 和 push 到 main/develop
- 前端: typecheck, lint, test
- 后端: go vet, go build, go test

#### deploy.yml - 部署
- 触发: push 到 main 并打 tag
- 部署前后端到生产服务器

### 环境变量配置
在 GitHub 仓库 Settings 中配置：
- Secrets: `SSH_PRIVATE_KEY`, `DB_PASSWORD`, `JWT_SECRET`
- Variables: `DEPLOY_SERVER_HOST`, `DEPLOY_FRONTEND_PATH`, `DEPLOY_BACKEND_PATH`

## 6. 部署

### 准备部署
1. 确保 main 分支代码通过所有检查
2. 创建版本标签: `git tag v1.0.0 && git push origin v1.0.0`
3. GitHub Actions 自动部署，或手动部署

### 手动部署检查清单
- [ ] 前端构建成功: `npm run build`
- [ ] 后端编译成功: `go build -o app ./cmd/main.go`
- [ ] 数据库迁移完成（如有）
- [ ] 环境变量配置正确（.env）
- [ ] systemd 服务已启动并运行正常

### 生产环境注意事项
- 使用 `NODE_ENV=production` 构建前端
- 后端使用 `GIN_MODE=release`
- 启用 HTTPS
- 配置正确的 CORS 源（AllowOrigins）
