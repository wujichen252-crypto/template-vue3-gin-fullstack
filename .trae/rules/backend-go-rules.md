# Go 后端开发规范

## 项目结构

```
backend/
├── cmd/              # 入口文件
├── config/          # 配置加载
├── internal/         # 内部包
│   ├── handler/     # 处理器层
│   ├── service/     # 业务逻辑层
│   ├── repository/  # 数据访问层
│   ├── model/       # 数据模型
│   └── middleware/  # 中间件
└── pkg/             # 公共包
    ├── response/   # 统一响应
    ├── jwt/         # JWT 工具
    └── logger/      # 日志工具
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

### 3. 数据库规范
- GORM 模型必须内嵌 gorm.Model
- 字段加 comment 标签
- 使用 transactions 处理多表操作

### 4. 接口规范
- RESTful 风格
- 统一响应格式: { "code": 200, "data": {}, "msg": "ok", "request_id": "uuid" }
- 使用 JWT 进行认证

## 函数长度
- 函数长度不超过 50 行
- 超过则拆分
