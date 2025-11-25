# 项目更名说明

## 项目名称变更

项目已从 **PlayEdu** 更名为 **EduFlow**，这是一个完全独立的企业培训平台项目。

## 变更内容

### 1. 项目标识

- **旧名称**: PlayEdu / playedu-go
- **新名称**: EduFlow / eduflow
- **Go 模块**: `github.com/eduflow/eduflow`

### 2. 数据库

- **数据库名**: `eduflow`
- **MinIO Bucket**: `eduflow`
- **配置文件**: 已更新所有 `playedu` 为 `eduflow`

### 3. 默认管理员账号

- **邮箱**: `admin@eduflow.com`
- **密码**: `eduflow123`
- **说明**: 首次登录后请立即修改密码

### 4. Docker 容器命名

所有 Docker 容器已重命名：
- `eduflow-mysql` (原 playedu-mysql)
- `eduflow-redis` (原 playedu-redis)
- `eduflow-minio` (原 playedu-minio)
- `eduflow-api` (原 playedu-api)

### 5. 网络命名

- Docker 网络: `eduflow-network` (原 playedu-network)

### 6. 文件和目录

- 二进制文件: `eduflow` (原 playedu)
- 日志文件: `logs/eduflow.log` (原 logs/playedu.log)
- 工作目录: `/opt/eduflow` (推荐)

## 代码变更

### Go 导入路径

所有 Go 文件中的导入路径已从：
```go
import "github.com/playedu/playedu-go/..."
```

更新为：
```go
import "github.com/eduflow/eduflow/..."
```

### 配置文件

`configs/config.yaml` 中的关键配置已更新：
- `database.dbname`: `eduflow`
- `minio.bucket_name`: `eduflow`
- `jwt.secret`: `eduflow-secret-key-change-in-production`
- `log.file_path`: `logs/eduflow.log`

## 升级指南

### 从 PlayEdu 迁移（如果适用）

1. **备份数据库**
   ```bash
   mysqldump -u root -p playedu > playedu_backup.sql
   ```

2. **创建新数据库**
   ```bash
   mysql -u root -p -e "CREATE DATABASE eduflow CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
   ```

3. **导入数据**
   ```bash
   mysql -u root -p eduflow < playedu_backup.sql
   ```

4. **更新管理员账号**
   ```sql
   UPDATE admin_users SET email = 'admin@eduflow.com' WHERE email = 'admin@playedu.xyz';
   ```

5. **更新配置文件**
   - 修改 `configs/config.yaml` 中的数据库名称
   - 更新其他相关配置

6. **重启服务**
   ```bash
   docker-compose down
   docker-compose up -d
   ```

### 全新部署

对于全新部署，直接按照 README.md 中的说明操作即可。

## 构建和部署

### 本地开发

```bash
# 克隆项目
git clone https://github.com/eduflow/eduflow.git
cd eduflow

# 安装依赖
go mod download

# 运行
go run cmd/api/main.go
```

### Docker 部署

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f api

# 停止服务
docker-compose down
```

### 生产环境编译

```bash
# Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o eduflow cmd/api/main.go

# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o eduflow.exe cmd/api/main.go

# macOS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o eduflow cmd/api/main.go
```

## 注意事项

1. **数据库兼容性**: 数据库结构完全兼容，只需更新数据库名称和管理员邮箱
2. **API 兼容性**: API 接口保持不变，客户端无需修改
3. **配置文件**: 确保所有配置文件中的引用已更新
4. **环境变量**: 如果使用环境变量，请更新相关名称

## 验证

部署完成后，可以通过以下方式验证：

1. **健康检查**
   ```bash
   curl http://localhost:8080/health
   ```

2. **管理员登录**
   - URL: http://localhost:8080/backend/v1/auth/login
   - Body: 
     ```json
     {
       "email": "admin@eduflow.com",
       "password": "eduflow123"
     }
     ```

3. **查看容器状态**
   ```bash
   docker-compose ps
   ```

## 支持

如有问题，请访问：
- GitHub Issues: https://github.com/eduflow/eduflow/issues
- 文档: 查看项目 `docs/` 目录

## 版本历史

- **v1.0.0** (2024-01-01): 项目初始版本（更名后）
- 基于 Go 1.21+ 完整实现
- 包含所有核心培训平台功能
