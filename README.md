# EduFlow - 企业培训平台

EduFlow 是一个功能完整的企业级在线培训平台，使用 Go 语言开发，提供完整的企业培训解决方案。

## 特性

- 🏢 **部门管理** - 支持多层级部门组织架构
- 👥 **学员管理** - 学员信息管理、部门关联、批量导入
- 👔 **管理员系统** - 管理员账户、角色权限管理
- 📚 **课程管理** - 课程创建、章节管理、课时安排
- 🎥 **视频学习** - 在线视频播放、进度跟踪、防快进
- 📊 **学习统计** - 学习时长统计、进度追踪、完成率分析
- 📦 **资源管理** - 视频、图片等资源的上传、分类、存储
- 📎 **课程附件** - 课程相关文档、资料下载
- 🔐 **LDAP集成** - 支持 LDAP/AD 域用户和部门同步
- ⚙️ **系统配置** - 系统参数配置、日志记录
- 🚦 **API限流** - 接口访问频率限制

## 技术栈

- **语言**: Go 1.21+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0+
- **缓存**: Redis
- **认证**: JWT
- **文件存储**: MinIO / S3
- **配置管理**: Viper

## 快速开始

### 前置要求

- Go 1.21 或更高版本
- MySQL 8.0+
- Redis
- MinIO (可选)

### 安装

1. 克隆项目

```bash
git clone https://github.com/eduflow/eduflow.git
cd eduflow
```

2. 安装依赖

```bash
go mod download
```

3. 配置环境

复制配置文件并修改：

```bash
cp configs/config.yaml configs/config.local.yaml
```

编辑 `configs/config.local.yaml` 修改数据库、Redis、MinIO 等配置。

4. 初始化数据库

```bash
mysql -u root -p < migrations/000001_init_schema.up.sql
```

5. 运行应用

```bash
go run cmd/api/main.go
```

或使用编译后的二进制文件：

```bash
go build -o eduflow cmd/api/main.go
./eduflow
```

应用将在 http://localhost:8080 启动。

### Docker 部署

使用 Docker Compose 快速部署：

```bash
docker-compose up -d
```

这将启动以下服务：
- API 服务 (端口 8080)
- MySQL 8.0 (端口 3306)
- Redis (端口 6379)
- MinIO (端口 9000, 9001)

## API 文档

### 默认管理员账号

- 邮箱: `admin@eduflow.com`
- 密码: `eduflow123`

### 后台管理 API (Backend)

**基础路径**: `/backend/v1`

#### 认证

- `POST /auth/login` - 管理员登录
- `POST /auth/logout` - 管理员登出
- `GET /auth/detail` - 获取当前管理员信息

#### 用户管理

- `GET /user` - 用户列表
- `POST /user` - 创建用户
- `GET /user/:id` - 获取用户详情
- `PUT /user/:id` - 更新用户
- `DELETE /user/:id` - 删除用户

#### 课程管理

- `GET /course` - 课程列表
- `POST /course` - 创建课程
- `GET /course/:id` - 获取课程详情
- `PUT /course/:id` - 更新课程
- `DELETE /course/:id` - 删除课程

#### 资源管理

- `GET /resource` - 资源列表
- `POST /resource/video/upload` - 上传视频
- `POST /resource/image/upload` - 上传图片
- `DELETE /resource/:id` - 删除资源

### 前台学员 API (Frontend)

**基础路径**: `/api/v1`

#### 认证

- `POST /auth/login` - 学员登录
- `POST /auth/register` - 学员注册
- `POST /auth/logout` - 学员登出
- `GET /auth/detail` - 获取当前学员信息

#### 课程

- `GET /courses` - 课程列表
- `GET /course/:id` - 获取课程详情

### 响应格式

成功响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {...}
}
```

错误响应：

```json
{
  "code": 1001,
  "message": "error message",
  "data": null
}
```

## 项目结构

```
eduflow/
├── cmd/
│   └── api/                    # 应用入口
├── internal/
│   ├── config/                 # 配置管理
│   ├── domain/                 # 领域模型
│   ├── repository/             # 数据访问层
│   ├── service/                # 业务逻辑层
│   ├── handler/                # HTTP处理器
│   │   ├── backend/           # 后台管理API
│   │   └── frontend/          # 前台学员API
│   ├── middleware/             # 中间件
│   └── pkg/                    # 内部工具包
│       ├── jwt/               # JWT工具
│       ├── storage/           # 文件存储
│       ├── crypto/            # 加密工具
│       └── response/          # 响应格式
├── pkg/                        # 公共工具包
│   ├── utils/
│   └── constants/
├── migrations/                 # 数据库迁移
├── configs/                    # 配置文件
├── docs/                       # 文档
└── README.md
```

## 开发指南

### 代码规范

- 遵循 Go 官方代码规范
- 使用 gofmt 格式化代码
- 使用 golangci-lint 进行代码检查

### 数据库迁移

向上迁移：

```bash
mysql -u root -p < migrations/000001_init_schema.up.sql
```

向下迁移：

```bash
mysql -u root -p < migrations/000001_init_schema.down.sql
```

### 测试

运行测试：

```bash
go test ./...
```

运行带覆盖率的测试：

```bash
go test -cover ./...
```

## 配置说明

主要配置项：

```yaml
server:
  port: 8080              # 服务端口
  mode: debug             # 运行模式: debug/release

database:
  host: localhost         # 数据库主机
  port: 3306             # 数据库端口
  user: root             # 数据库用户
  password: password     # 数据库密码
  dbname: eduflow        # 数据库名称

redis:
  host: localhost        # Redis主机
  port: 6379            # Redis端口
  password: ""          # Redis密码
  db: 0                 # Redis数据库

jwt:
  secret: your-secret-key    # JWT密钥
  expire_time: 360h          # Token有效期

minio:
  endpoint: localhost:9000    # MinIO地址
  access_key_id: minioadmin  # 访问密钥
  secret_access_key: minioadmin  # 密钥
  bucket_name: eduflow       # 存储桶名称

rate_limit:
  duration: 60          # 限流时间窗口(秒)
  limit: 360           # 限流次数
```

## 部署

### 生产环境部署

1. 编译生产版本

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o eduflow cmd/api/main.go
```

2. 配置 systemd 服务

创建 `/etc/systemd/system/eduflow.service`:

```ini
[Unit]
Description=EduFlow API Service
After=network.target

[Service]
Type=simple
User=eduflow
WorkingDirectory=/opt/eduflow
ExecStart=/opt/eduflow/eduflow
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

3. 启动服务

```bash
sudo systemctl daemon-reload
sudo systemctl enable eduflow
sudo systemctl start eduflow
```

### Nginx 反向代理

```nginx
upstream eduflow_api {
    server 127.0.0.1:8080;
}

server {
    listen 80;
    server_name api.eduflow.com;
    
    location / {
        proxy_pass http://eduflow_api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目基于 Apache License 2.0 许可证开源。详见 [LICENSE](LICENSE) 文件。

## 联系方式

- 项目主页: https://github.com/eduflow/eduflow
- 问题反馈: https://github.com/eduflow/eduflow/issues

## 更新日志

### v1.0.0 (2024-01-01)

- ✨ 初始版本发布
- ✅ 完整的用户和部门管理
- ✅ 课程管理功能
- ✅ 资源管理和文件上传
- ✅ JWT 认证和权限控制
- ✅ API 限流
- ✅ Docker 支持
