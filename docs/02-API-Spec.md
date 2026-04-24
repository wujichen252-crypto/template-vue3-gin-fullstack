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

## 通用错误码

| code | 说明 |
|------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（Token无效、过期、或在黑名单） |
| 403 | 禁止访问（用户已被禁用） |
| 404 | 资源不存在 |
| 409 | 资源冲突（用户名或邮箱已存在） |
| 429 | 请求过于频繁（限流触发） |
| 500 | 服务器内部错误 |

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

**请求参数说明**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 4-64字符，唯一 |
| password | string | 是 | 6-128字符 |
| email | string | 是 | 合法邮箱格式，唯一 |

**响应（成功）**
```json
{
  "code": 200,
  "data": {
    "user_info": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "avatar_url": "",
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z"
    }
  },
  "msg": "ok",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**错误码**
| code | 说明 |
|------|------|
| 400 | 参数错误（用户名/密码/邮箱格式不正确） |
| 409 | 用户名或邮箱已被注册 |
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

**请求参数说明**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**响应（成功）**
```json
{
  "code": 200,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
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

**响应字段说明**
| 字段 | 类型 | 说明 |
|------|------|------|
| token | string | Access Token，有效期2小时 |
| refresh_token | string | Refresh Token，有效期7天 |
| user_info | object | 用户基本信息 |

**错误码**
| code | 说明 |
|------|------|
| 400 | 参数错误 |
| 401 | 用户名或密码错误 |
| 403 | 用户已被禁用 |
| 500 | 服务器错误 |

---

### 3. 获取用户信息

**请求**
```http
GET /api/v1/auth/userinfo
Authorization: Bearer <token>
```

**请求头说明**
| 头部 | 必填 | 说明 |
|------|------|------|
| Authorization | 是 | Bearer {token} |

**响应（成功）**
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
| 401 | 未授权（Token无效、过期、或在黑名单） |
| 500 | 服务器错误 |

---

### 4. 刷新Token

**请求**
```http
POST /api/v1/auth/refresh
Authorization: Bearer <token>
```

**请求头说明**
| 头部 | 必填 | 说明 |
|------|------|------|
| Authorization | 是 | Bearer {refresh_token} |

**响应（成功）**
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

**说明**
- 该接口需要使用 Refresh Token
- 返回新的 Access Token

**错误码**
| code | 说明 |
|------|------|
| 401 | 未授权（不是Refresh Token、已过期、或在黑名单） |
| 403 | 用户已被禁用 |
| 500 | 服务器错误 |

---

### 5. 用户登出

**请求**
```http
POST /api/v1/auth/logout
Authorization: Bearer <token>
```

**请求头说明**
| 头部 | 必填 | 说明 |
|------|------|------|
| Authorization | 是 | Bearer {token} |

**响应（成功）**
```json
{
  "code": 200,
  "data": null,
  "msg": "ok",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**说明**
- 将当前 Token 加入黑名单
- 黑名单过期时间 = Token 剩余有效期
- 黑名单内的 Token 将无法使用

**错误码**
| code | 说明 |
|------|------|
| 401 | 未授权（Token无效或已过期） |
| 500 | 服务器错误 |

---

## 健康检查接口

### GET /health

**请求**
```http
GET /health
```

**响应（成功）**
```json
{
  "status": "ok",
  "time": "2024-01-01 12:00:00"
}
```

---

## 状态码说明

### 用户状态 (status)
| 值 | 说明 |
|----|------|
| 1 | 正常 |
| 0 | 禁用 |

### Token 类型
| 类型 | 有效期 | 用途 |
|------|--------|------|
| Access Token | 2小时 | API访问认证 |
| Refresh Token | 7天 | 刷新Access Token |
