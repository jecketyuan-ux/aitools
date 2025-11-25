# EduFlow 快速启动指南

## 5分钟快速体验

### 使用 Docker Compose（推荐）

```bash
# 1. 克隆项目
git clone https://github.com/eduflow/eduflow.git
cd eduflow

# 2. 启动所有服务
docker-compose up -d

# 3. 等待服务启动（约30秒）
docker-compose logs -f api

# 4. 访问应用
# API: http://localhost:8080
# MinIO Console: http://localhost:9001
```

### 默认登录信息

**管理员账号**:
- 邮箱: `admin@eduflow.com`
- 密码: `eduflow123`

**MinIO**:
- 用户名: `minioadmin`
- 密码: `minioadmin123`

**MySQL**:
- 用户名: `eduflow`
- 密码: `eduflow123`
- 数据库: `eduflow`

## 测试 API

### 1. 健康检查

```bash
curl http://localhost:8080/health
```

预期响应：
```json
{
  "status": "ok"
}
```

### 2. 管理员登录

```bash
curl -X POST http://localhost:8080/backend/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@eduflow.com",
    "password": "eduflow123"
  }'
```

预期响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "admin": {
      "id": 1,
      "name": "Administrator",
      "email": "admin@eduflow.com",
      ...
    }
  }
}
```

### 3. 创建用户（需要先登录获取 token）

```bash
TOKEN="your_token_here"

curl -X POST http://localhost:8080/backend/v1/user \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "email": "user@example.com",
    "name": "Test User",
    "password": "password123"
  }'
```

### 4. 获取用户列表

```bash
curl -X GET "http://localhost:8080/backend/v1/user?page=1&size=10" \
  -H "Authorization: Bearer $TOKEN"
```

### 5. 创建课程

```bash
curl -X POST http://localhost:8080/backend/v1/course \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Go 语言入门",
    "short_desc": "从零开始学习 Go 语言",
    "is_required": 1,
    "is_show": 1
  }'
```

## 前端学员 API

### 1. 学员注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@example.com",
    "name": "Student Name",
    "password": "password123"
  }'
```

### 2. 学员登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@example.com",
    "password": "password123"
  }'
```

### 3. 浏览课程

```bash
STUDENT_TOKEN="student_token_here"

curl -X GET "http://localhost:8080/api/v1/courses?page=1&size=10" \
  -H "Authorization: Bearer $STUDENT_TOKEN"
```

### 4. 查看课程详情

```bash
curl -X GET http://localhost:8080/api/v1/course/1 \
  -H "Authorization: Bearer $STUDENT_TOKEN"
```

## 上传资源

### 上传图片

```bash
curl -X POST http://localhost:8080/backend/v1/resource/image/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/image.jpg"
```

### 上传视频

```bash
curl -X POST http://localhost:8080/backend/v1/resource/video/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/video.mp4"
```

## 常用命令

### Docker Compose

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 查看 API 日志
docker-compose logs -f api

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 清理所有数据（包括数据库）
docker-compose down -v
```

### 查看服务状态

```bash
# 查看所有容器
docker-compose ps

# 查看 API 容器日志
docker logs eduflow-api

# 进入 MySQL 容器
docker exec -it eduflow-mysql mysql -u eduflow -p

# 进入 Redis 容器
docker exec -it eduflow-redis redis-cli
```

## 数据库操作

### 连接数据库

```bash
# 使用 Docker
docker exec -it eduflow-mysql mysql -u eduflow -p eduflow

# 或直接连接
mysql -h 127.0.0.1 -P 3306 -u eduflow -p eduflow
```

### 常用查询

```sql
-- 查看所有管理员
SELECT * FROM admin_users;

-- 查看所有用户
SELECT * FROM users;

-- 查看所有课程
SELECT * FROM courses WHERE deleted_at IS NULL;

-- 查看用户学习记录
SELECT * FROM user_course_records;
```

## 故障排查

### 服务无法启动

1. 检查端口占用
```bash
netstat -tuln | grep -E '3306|6379|8080|9000|9001'
```

2. 查看容器日志
```bash
docker-compose logs
```

3. 重启服务
```bash
docker-compose restart
```

### 数据库连接失败

1. 检查 MySQL 容器状态
```bash
docker-compose ps mysql
```

2. 检查 MySQL 日志
```bash
docker-compose logs mysql
```

3. 等待 MySQL 完全启动
```bash
docker-compose logs -f mysql | grep "ready for connections"
```

### API 返回 500 错误

1. 查看 API 日志
```bash
docker-compose logs api
```

2. 检查配置文件
```bash
cat configs/config.yaml
```

3. 验证数据库连接
```bash
docker exec -it eduflow-mysql mysql -u eduflow -p -e "SHOW DATABASES;"
```

## 停止和清理

### 停止服务但保留数据

```bash
docker-compose down
```

### 完全清理（包括数据）

```bash
docker-compose down -v
```

### 清理 Docker 镜像

```bash
docker-compose down --rmi all
```

## 下一步

现在你已经成功启动了 EduFlow，可以：

1. **阅读完整文档**: 查看 `README.md` 和 `docs/` 目录
2. **了解架构设计**: 阅读 `DESIGN.md`
3. **查看实现细节**: 阅读 `docs/IMPLEMENTATION_SUMMARY.md`
4. **开始开发**: 参考 `Makefile` 中的开发命令
5. **部署到生产**: 查看 README 中的部署指南

## 获取帮助

- 📖 文档: 项目 `docs/` 目录
- 🐛 问题反馈: [GitHub Issues](https://github.com/eduflow/eduflow/issues)
- 💬 讨论: [GitHub Discussions](https://github.com/eduflow/eduflow/discussions)

## 安全建议

⚠️ **生产环境部署前请务必：**

1. 修改所有默认密码
2. 更新 JWT 密钥
3. 配置 HTTPS
4. 启用防火墙
5. 定期备份数据库
6. 更新依赖包
7. 配置日志轮转
8. 设置监控告警

---

祝你使用愉快！ 🚀
