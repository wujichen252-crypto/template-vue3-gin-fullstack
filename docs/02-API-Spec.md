# 接口规范文档

## 概述
- 基础路径: `/api/v1`
- 认证方式: JWT Bearer Token
- 数据格式: JSON

## 通用响应格式

### 成功响应
```json
{
  "code": 200,
  "data": {},
  "msg": "ok",
  "request_id": "uuid"
}
```

### 错误响应
```json
{
  "code": 400,
  "data": null,
  "msg": "错误信息",
  "request_id": "uuid"
}
```

## 认证接口

### 1. 用户注册

**请求**
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com"
}
```

**响应**
```json
{
  "code": 200,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_info": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "avatar_url": "",
      "status": 1
    }
  },
  "msg": "ok",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**错误码**
| code | 说明 |
|------|------|
| 400 | 参数错误 |
| 500 | 服务器错误 |

---

### 2. 用户登录

**请求**
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

**响应**
```json
{
  "code": 200,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_info": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "avatar_url": "",
      "status": 1
    }
  },
  "msg": "ok",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**错误码**
| code | 说明 |
|------|------|
| 400 | 参数错误 |
| 401 | 用户名或密码错误 |

---

### 3. 获取用户信息

**请求**
```http
GET /api/v1/auth/userinfo
Authorization: Bearer <token>
```

**响应**
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "avatar_url": "",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "msg": "ok",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**错误码**
| code | 说明 |
|------|------|
| 401 | 未授权 |

---

### 4. 刷新Token

**请求**
```http
POST /api/v1/auth/refresh
Authorization: Bearer <token>
```

**响应**
```json
{
  "code": 200,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "msg": "ok",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**错误码**
| code | 说明 |
|------|------|
| 401 | 未授权 |
