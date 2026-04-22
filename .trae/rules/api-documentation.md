# API 文档规范

## RESTful 接口设计

### 基本原则
- 使用标准的 HTTP 方法
- 使用名词而非动词
- 使用复数形式
- 层级结构使用 URL 路径

### HTTP 方法映射

| 方法 | 用途 | 示例 |
|------|------|------|
| GET | 查询资源 | GET /users |
| POST | 创建资源 | POST /users |
| PUT | 更新资源 | PUT /users/1 |
| DELETE | 删除资源 | DELETE /users/1 |

### 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器错误 |

## 响应格式

```json
{
  "code": 200,
  "data": {},
  "msg": "ok",
  "request_id": "uuid"
}
```

## 认证

- 使用 JWT Bearer Token
- Token 放在 Authorization Header
- 格式: `Authorization: Bearer <token>`

## 请求示例

### Headers
```
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Body
```json
{
  "username": "test",
  "password": "123456"
}
```
