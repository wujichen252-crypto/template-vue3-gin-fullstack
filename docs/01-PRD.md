# 需求文档

## 1. 项目概述

### 项目名称
NeuraMind 智脑笔记系统（Vue3 + Gin 全栈模板）

### 项目背景
本项目是一个前后端分离的全栈开发模板，旨在为开发者提供一套开箱即用的 Vue3 + Go Gin 项目架构。模板包含完整的用户认证系统、标准化的 API 规范、可配置的 CI/CD 流水线，适合作为各类 Web 应用的基础框架。

### 项目目标
- 提供一套完整的前后端分离架构模板
- 实现标准化的 RESTful API 规范
- 包含可配置的 GitHub Actions CI/CD 流水线
- 覆盖单元测试、集成测试的完整测试体系
- 支持 Docker 容器化部署

## 2. 功能需求

### 2.1 模块列表

| 模块 | 功能描述 | 优先级 |
|------|---------|--------|
| 用户模块 | 用户注册、登录、登出、信息管理、Token刷新 | P0 |
| 笔记模块 | 笔记CRUD、AI摘要生成（预留接口） | P1 |
| 标签模块 | 标签管理（预留） | P2 |
| 分享模块 | 笔记分享链接（预留） | P2 |

### 2.2 详细需求

#### 用户模块

##### 用户注册
- 功能描述: 新用户注册
- 输入: 用户名、密码、邮箱
- 处理流程:
  1. 验证输入参数（用户名4-64字符、密码6-128字符、邮箱格式）
  2. 检查用户名是否已存在
  3. 检查邮箱是否已存在
  4. 密码使用 bcrypt 加密存储
  5. 创建用户记录
  6. 返回用户信息（不自动登录）
- 输出: 用户信息
- 优先级: P0

##### 用户登录
- 功能描述: 用户登录系统
- 输入: 用户名、密码
- 处理流程:
  1. 验证输入参数
  2. 查询用户
  3. 验证密码（bcrypt）
  4. 检查用户状态（1=正常，0=禁用）
  5. 生成 JWT Access Token 和 Refresh Token
  6. 返回用户信息和 Token
- 输出: 用户信息、Access Token、Refresh Token
- 优先级: P0

##### 获取用户信息
- 功能描述: 获取当前登录用户的详细信息
- 输入: JWT Token（在 Authorization Header）
- 处理流程:
  1. 验证 Token 有效性
  2. 检查 Token 是否在黑名单（Redis）
  3. 从 Token 中提取用户ID
  4. 查询用户信息
  5. 返回用户信息
- 输出: 用户信息（包含 created_at、updated_at）
- 优先级: P0

##### 刷新Token
- 功能描述: 使用 Refresh Token 获取新的 Access Token
- 输入: JWT Token（在 Authorization Header）
- 处理流程:
  1. 验证 Token 有效性
  2. 检查是否为 Refresh Token
  3. 从 Token 中提取用户ID
  4. 检查用户状态
  5. 生成新的 Access Token
- 输出: 新的 Access Token
- 优先级: P0

##### 用户登出
- 功能描述: 使当前 Token 失效
- 输入: JWT Token（在 Authorization Header）
- 处理流程:
  1. 验证 Token 有效性
  2. 将 Token 加入 Redis 黑名单
  3. 设置黑名单过期时间为 Token 剩余有效期
  4. 返回成功
- 输出: 操作结果
- 优先级: P0

## 3. 非功能需求

### 3.1 性能需求
- 响应时间: < 200ms（P95）
- 并发用户: >= 100
- 数据库连接池: 最大100连接

### 3.2 安全需求
- 密码加密存储（bcrypt cost=10）
- JWT Token 认证（HS256 签名）
- Token 黑名单机制（Redis）
- 防止 SQL 注入（GORM 参数化查询）
- 防止 XSS 攻击（前端输出转义）
- CORS 跨域限制（可配置允许来源）
- 请求限流（Redis + Token Bucket）

### 3.3 兼容性需求
- 浏览器: Chrome, Firefox, Safari, Edge 最新版本
- 移动端: iOS 12+, Android 8+
- 后端: Go 1.22+, PostgreSQL 14+, Redis 6+

### 3.4 可靠性需求
- Graceful Shutdown（优雅关闭）
- 数据库连接超时处理
- Redis 连接失败降级

## 4. 数据字典

### 4.1 用户表 (users)

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | uint | 主键 | 自增, PK |
| created_at | datetime | 创建时间 | 自动填充 |
| updated_at | datetime | 更新时间 | 自动填充 |
| deleted_at | gorm.DeletedAt | 软删除时间 | 可为空,索引 |
| username | string(64) | 用户名 | 唯一索引, 必填 |
| password_hash | string(255) | 密码哈希 | 必填 |
| email | string(128) | 邮箱 | 唯一索引, 必填 |
| avatar_url | string(255) | 头像URL | 可为空 |
| status | int8 | 状态 | 默认1（1=正常,0=禁用） |

### 4.2 索引设计

| 索引名 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| idx_users_deleted_at | deleted_at | BTree | 软删除查询优化 |
| idx_users_username | username | UNIQUE BTree | 用户名唯一约束 |
| idx_users_email | email | UNIQUE BTree | 邮箱唯一约束 |

## 5. 风险评估

| 风险 | 影响 | 概率 | 应对措施 |
|------|------|------|----------|
| 用户密码泄露 | 高 | 低 | bcrypt加密、密码强度验证、日志脱敏 |
| JWT Token 被盗用 | 高 | 中 | Token黑名单、短期AccessToken、RefreshToken轮换 |
| 数据库密码泄露 | 高 | 低 | 环境变量配置、不提交.env文件到版本库 |
| Redis 服务中断 | 中 | 低 | JWT验证降级、Token黑名单功能暂时失效 |
| 并发登录导致密码验证失败 | 低 | 中 | 数据库连接池配置、限流保护 |
| 软删除数据泄露 | 低 | 低 | deleted_at 索引优化、查询自动过滤 |

## 6. 技术架构

### 6.1 前端技术栈
- Vue 3.4 + TypeScript 5
- Vite 5（构建工具）
- Pinia 2（状态管理）
- Vue Router 4（路由）
- Axios（HTTP 客户端）
- Tailwind CSS 3（样式）
- Element Plus（UI 组件库）
- Vitest（单元测试）

### 6.2 后端技术栈
- Go 1.22（Golang）
- Gin 1.9（Web 框架）
- GORM 1.30（ORM）
- PostgreSQL 14+（关系数据库）
- Redis 6+（缓存/Token黑名单）
- golang-jwt/jwt/v5（JWT 认证）
- Zap（结构化日志）
- Swaggo/swag（Swagger 文档）

### 6.3 DevOps
- GitHub Actions（CI/CD）
- Docker（容器化）
- systemd（服务管理）
