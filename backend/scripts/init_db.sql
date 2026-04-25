-- 创建数据库
CREATE DATABASE template_db;

-- 连接数据库
\c template_db;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    username VARCHAR(64) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(128) NOT NULL,
    avatar_url VARCHAR(255) DEFAULT '',
    status SMALLINT DEFAULT 1
);

-- 创建唯一约束和索引
ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username);
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- 注释
COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.username IS '用户名';
COMMENT ON COLUMN users.password_hash IS '密码哈希';
COMMENT ON COLUMN users.email IS '邮箱';
COMMENT ON COLUMN users.avatar_url IS '头像URL';
COMMENT ON COLUMN users.status IS '状态：1正常，0禁用';
