# 项目开发流程

## 1. 开发准备

### 环境配置
- 安装 Node.js >= 18
- 安装 Go >= 1.22
- 安装 PostgreSQL >= 14
- 安装 Redis >= 6
- 配置 Git

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

## 2. 开发流程

### 分支管理
- main: 主分支，保护分支
- develop: 开发分支
- feature/*: 功能分支
- bugfix/*: 修复分支

### 开发步骤
1. 从 develop 创建功能分支
2. 编写代码
3. 本地测试
4. 提交代码 (遵守提交规范)
5. 创建 Pull Request
6. 代码审查
7. 合并到 develop

### 提交规范
```
<type>(<scope>): <subject>

<body>

<footer>
```

类型:
- feat: 新功能
- fix: 修复
- docs: 文档
- style: 格式
- refactor: 重构
- test: 测试
- chore: 构建

## 3. 测试

### 前端测试
```bash
cd frontend
npm run typecheck
npm run lint
```

### 后端测试
```bash
cd backend
go build ./...
go test ./...
```

## 4. 部署

### 准备部署
1. 合并到 main 分支
2. 创建 Tag
3. 构建生产版本

### 部署检查
- [ ] 前端构建成功
- [ ] 后端编译成功
- [ ] 数据库迁移完成
- [ ] 环境变量配置正确
