# Go 后端开发规范

## 项目结构

```
backend/
├── cmd/              # 入口文件（main.go）
├── config/          # 配置加载
├── internal/         # 内部包
│   ├── handler/     # 处理器层（HTTP 请求处理）
│   ├── service/     # 业务逻辑层
│   ├── repository/  # 数据访问层
│   ├── model/       # 数据模型
│   └── middleware/  # 中间件
├── pkg/             # 公共包
│   ├── response/   # 统一响应
│   ├── jwt/         # JWT 工具
│   └── logger/      # 日志工具
└── scripts/         # 脚本（数据库初始化等）
```

## 分层规范

### Handler 层
- 处理请求参数验证
- 调用 Service 层
- 返回统一响应格式
- 禁止写业务逻辑

### Service 层
- 编写业务逻辑
- 调用 Repository 层
- 处理业务错误

### Repository 层
- 数据增删改查
- 禁止业务逻辑

## 代码规范

### 1. 命名规范
- 包名: 简短的小写单词
- 结构体: PascalCase
- 变量/函数: camelCase
- 常量: MixedCase 或 UPPER_SNAKE_CASE

### 2. 错误处理
- 使用自定义错误
- 统一错误码
- 不在 Handler 层返回原始错误
- 使用 errors.Wrap 包装错误并保留上下文

### 3. 数据库规范
- GORM 模型必须内嵌 gorm.Model
- 字段加 comment 标签
- 使用 transactions 处理多表操作
- 使用软删除（gorm.DeletedAt）

### 4. 接口规范
- RESTful 风格
- 统一响应格式: { "code": 200, "data": {}, "msg": "ok", "request_id": "uuid" }
- 使用 JWT 进行认证

### 5. 日志规范
- 使用结构化日志（zap）
- 记录请求 ID（request_id）
- 敏感信息脱敏处理

## 函数长度
- 函数长度不超过 50 行
- 超过则拆分

## 测试规范

### 测试框架
- testing（标准库）
- testify/assert（断言）
- gomock（Mock 测试）

### 测试文件命名
- 单元测试: `*_test.go`
- 放在同目录

### 测试命令
```bash
# 运行所有测试
go test ./...

# 运行测试并显示详情
go test -v ./...

# 运行测试并检查覆盖率
go test -cover ./...

# 运行特定包的测试
go test -v ./internal/repository/...
```

### 测试覆盖率要求
- 核心业务逻辑（service）覆盖率 > 70%
- 数据访问层（repository）覆盖率 > 60%

### Mock 测试
使用 gomock 进行依赖 Mock：
```bash
# 安装 mockgen
go install github.com/golang/mock/mockgen@latest

# 生成 Mock 文件
mockgen -source=./internal/repository/user_repository.go -destination=./internal/repository/mocks/user_repository_mock.go
```
