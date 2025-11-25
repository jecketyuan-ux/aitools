# EduFlow 部署手册

## 目录
1. [部署概述](#1-部署概述)
2. [环境要求](#2-环境要求)
3. [快速部署](#3-快速部署)
4. [详细部署步骤](#4-详细部署步骤)
5. [配置说明](#5-配置说明)
6. [验证部署](#6-验证部署)
7. [故障排除](#7-故障排除)

---

## 1. 部署概述

### 1.1 部署架构
EduFlow采用容器化部署方案，包含以下服务：
- **eduflow-api**: 主应用服务
- **mysql**: 数据库服务
- **redis**: 缓存服务
- **minio**: 文件存储服务

### 1.2 部署方式
- **Docker Compose**: 推荐的部署方式，适合单机部署
- **手动部署**: 适合有特殊需求的场景
- **Kubernetes**: 适合生产环境和大规模部署

### 1.3 网络架构
```
Internet
    │
    ▼
┌─────────────┐
│  Load       │
│  Balancer   │
└─────────────┘
    │
    ▼
┌─────────────┐
│  EduFlow    │
│  API:8080   │
└─────────────┘
    │
    ▼
┌─────────────────────────────────────────┐
│  Internal Network (eduflow-network)     │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐   │
│  │ MySQL   │ │ Redis   │ │ MinIO   │   │
│  │ :3306   │ │ :6379   │ │ :9000   │   │
│  └─────────┘ └─────────┘ └─────────┘   │
└─────────────────────────────────────────┘
```

---

## 2. 环境要求

### 2.1 硬件要求

#### 2.1.1 最小配置
- **CPU**: 2核心
- **内存**: 4GB RAM
- **存储**: 20GB 可用空间
- **网络**: 100Mbps

#### 2.1.2 推荐配置
- **CPU**: 4核心
- **内存**: 8GB RAM
- **存储**: 100GB SSD
- **网络**: 1Gbps

#### 2.1.3 生产环境配置
- **CPU**: 8核心
- **内存**: 16GB RAM
- **存储**: 500GB SSD
- **网络**: 10Gbps

### 2.2 软件要求

#### 2.2.1 操作系统
- **Linux**: Ubuntu 20.04+ / CentOS 8+ / RHEL 8+
- **Docker**: 20.10+
- **Docker Compose**: 2.0+

#### 2.2.2 端口要求
- **8080**: EduFlow API服务
- **3306**: MySQL数据库（内部）
- **6379**: Redis缓存（内部）
- **9000**: MinIO API（内部）
- **9001**: MinIO Console（内部）

### 2.3 安全要求
- 防火墙配置
- SSL证书配置
- 数据库访问控制
- 文件存储访问控制

---

## 3. 快速部署

### 3.1 一键部署（推荐）
```bash
# 1. 克隆项目
git clone https://github.com/eduflow/eduflow.git
cd eduflow

# 2. 启动所有服务
docker-compose up -d

# 3. 等待服务启动（约2-3分钟）
docker-compose logs -f

# 4. 访问系统
# 管理后台: http://localhost:8080/backend
# 用户前端: http://localhost:8080
```

### 3.2 默认账户
```
管理员账户:
邮箱: admin@eduflow.com
密码: eduflow123

MinIO管理:
用户名: minioadmin
密码: minioadmin123
```

---

## 4. 详细部署步骤

### 4.1 准备工作

#### 4.1.1 安装Docker
```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# CentOS/RHEL
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install -y docker-ce docker-ce-cli containerd.io
sudo systemctl start docker
sudo systemctl enable docker
```

#### 4.1.2 安装Docker Compose
```bash
# 下载Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker-compose version
```

#### 4.1.3 下载项目代码
```bash
# 克隆项目
git clone https://github.com/eduflow/eduflow.git
cd eduflow

# 检查项目结构
ls -la
```

### 4.2 配置文件准备

#### 4.2.1 环境变量配置
创建 `.env` 文件：
```bash
# 数据库配置
DB_HOST=mysql
DB_PORT=3306
DB_NAME=eduflow
DB_USER=eduflow
DB_PASSWORD=eduflow123

# Redis配置
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=

# MinIO配置
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin123

# 应用配置
APP_PORT=8080
JWT_SECRET=eduflow-secret-key-change-in-production
```

#### 4.2.2 应用配置文件
修改 `configs/config.yaml`：
```yaml
server:
  port: 8080
  mode: release  # 生产环境使用release模式
  read_timeout: 60s
  write_timeout: 60s

database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  dbname: ${DB_NAME}
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600s

redis:
  host: ${REDIS_HOST}
  port: ${REDIS_PORT}
  password: ${REDIS_PASSWORD}
  db: 0

jwt:
  secret: ${JWT_SECRET}
  expire_time: 360h  # 15天

minio:
  endpoint: ${MINIO_ENDPOINT}
  access_key_id: ${MINIO_ACCESS_KEY}
  secret_access_key: ${MINIO_SECRET_KEY}
  use_ssl: false
  bucket_name: eduflow

rate_limit:
  duration: 60
  limit: 360

log:
  level: info
  file_path: logs/eduflow.log
  max_size: 100
  max_backups: 10
  max_age: 30
```

### 4.3 Docker Compose配置

#### 4.3.1 生产环境配置
修改 `docker-compose.yml`：
```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: eduflow-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_ROOT_PASSWORD}
      MYSQL_DATABASE: ${DB_NAME}
      MYSQL_USER: ${DB_USER}
      MYSQL_PASSWORD: ${DB_PASSWORD}
    ports:
      - "127.0.0.1:3306:3306"  # 仅本地访问
    volumes:
      - mysql_data:/var/lib/mysql
      - ./migrations:/docker-entrypoint-initdb.d
    networks:
      - eduflow-network
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      timeout: 20s
      retries: 10
      interval: 30s

  redis:
    image: redis:7-alpine
    container_name: eduflow-redis
    restart: unless-stopped
    ports:
      - "127.0.0.1:6379:6379"  # 仅本地访问
    volumes:
      - redis_data:/data
    networks:
      - eduflow-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 30s
      timeout: 10s
      retries: 3

  minio:
    image: minio/minio:latest
    container_name: eduflow-minio
    restart: unless-stopped
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: ${MINIO_ACCESS_KEY}
      MINIO_ROOT_PASSWORD: ${MINIO_SECRET_KEY}
    ports:
      - "127.0.0.1:9000:9000"  # 仅本地访问
      - "127.0.0.1:9001:9001"  # 仅本地访问
    volumes:
      - minio_data:/data
    networks:
      - eduflow-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: eduflow-api
    restart: unless-stopped
    ports:
      - "0.0.0.0:8080:8080"  # 对外开放
    env_file:
      - .env
    volumes:
      - ./logs:/app/logs
      - ./configs:/app/configs
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
      minio:
        condition: service_healthy
    networks:
      - eduflow-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  mysql_data:
    driver: local
  redis_data:
    driver: local
  minio_data:
    driver: local

networks:
  eduflow-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

### 4.4 构建和启动

#### 4.4.1 构建镜像
```bash
# 构建应用镜像
docker-compose build

# 查看镜像
docker images | grep eduflow
```

#### 4.4.2 启动服务
```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

#### 4.4.3 等待服务就绪
```bash
# 等待数据库就绪
docker-compose exec mysql mysqladmin ping -h localhost

# 等待Redis就绪
docker-compose exec redis redis-cli ping

# 等待MinIO就绪
curl http://localhost:9000/minio/health/live

# 等待API服务就绪
curl http://localhost:8080/health
```

### 4.5 数据库初始化

#### 4.5.1 自动初始化
数据库会通过 `migrations/000001_init_schema.up.sql` 自动初始化。

#### 4.5.2 手动初始化（如需要）
```bash
# 连接数据库
docker-compose exec mysql mysql -u root -p

# 手动执行初始化脚本
mysql -u root -p eduflow < migrations/000001_init_schema.up.sql
```

### 4.6 SSL证书配置

#### 4.6.1 使用Nginx反向代理
创建 `nginx.conf`：
```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/your/cert.pem;
    ssl_certificate_key /path/to/your/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

---

## 5. 配置说明

### 5.1 数据库配置

#### 5.1.1 连接池配置
```yaml
database:
  max_idle_conns: 10      # 最大空闲连接数
  max_open_conns: 100     # 最大打开连接数
  conn_max_lifetime: 3600s # 连接最大生存时间
```

#### 5.1.2 性能优化
- 根据并发量调整连接池大小
- 定期优化数据库表
- 配置适当的索引

### 5.2 Redis配置

#### 5.2.1 内存配置
```yaml
redis:
  maxmemory: 2gb          # 最大内存
  maxmemory-policy: allkeys-lru  # 内存淘汰策略
```

#### 5.2.2 持久化配置
```yaml
redis:
  save: "900 1 300 10 60 10000"  # RDB持久化
  appendonly: yes                 # AOF持久化
```

### 5.3 MinIO配置

#### 5.3.1 存储配置
```yaml
minio:
  use_ssl: true           # 生产环境启用SSL
  region: us-east-1       # 存储区域
```

#### 5.3.2 访问控制
- 设置访问密钥
- 配置存储桶策略
- 启用访问日志

### 5.4 应用配置

#### 5.4.1 安全配置
```yaml
jwt:
  secret: "your-super-secret-key"  # 使用强密钥
  expire_time: 360h                 # 根据需求调整

rate_limit:
  duration: 60      # 限流时间窗口
  limit: 360        # 限流请求数
```

#### 5.4.2 日志配置
```yaml
log:
  level: info        # 日志级别
  file_path: logs/eduflow.log
  max_size: 100      # 日志文件大小(MB)
  max_backups: 10    # 保留日志文件数
  max_age: 30        # 日志保留天数
```

---

## 6. 验证部署

### 6.1 健康检查

#### 6.1.1 API服务检查
```bash
# 检查API服务
curl http://localhost:8080/health

# 预期响应
{"status":"ok"}
```

#### 6.1.2 数据库检查
```bash
# 检查数据库连接
docker-compose exec mysql mysql -u root -p -e "SHOW DATABASES;"

# 检查表结构
docker-compose exec mysql mysql -u root -p eduflow -e "SHOW TABLES;"
```

#### 6.1.3 Redis检查
```bash
# 检查Redis连接
docker-compose exec redis redis-cli ping

# 检查Redis信息
docker-compose exec redis redis-cli info
```

#### 6.1.4 MinIO检查
```bash
# 检查MinIO服务
curl http://localhost:9000/minio/health/live

# 检查存储桶
docker-compose exec minio mc ls local/
```

### 6.2 功能验证

#### 6.2.1 管理员登录
1. 访问: http://localhost:8080/backend
2. 使用默认管理员账户登录
3. 验证管理后台功能

#### 6.2.2 用户注册
1. 访问: http://localhost:8080
2. 点击注册
3. 填写注册信息
4. 验证登录功能

#### 6.2.3 课程管理
1. 登录管理后台
2. 创建测试课程
3. 上传测试资源
4. 验证课程发布

### 6.3 性能测试

#### 6.3.1 压力测试
```bash
# 使用Apache Bench
ab -n 1000 -c 100 http://localhost:8080/health

# 使用wrk
wrk -t12 -c400 -d30s http://localhost:8080/health
```

#### 6.3.2 负载测试
- 模拟多用户并发访问
- 监控系统资源使用
- 测试数据库性能

---

## 7. 故障排除

### 7.1 常见问题

#### 7.1.1 服务启动失败
```bash
# 查看服务状态
docker-compose ps

# 查看服务日志
docker-compose logs [service_name]

# 重启服务
docker-compose restart [service_name]
```

#### 7.1.2 数据库连接失败
```bash
# 检查数据库服务
docker-compose exec mysql mysqladmin ping

# 检查网络连接
docker network ls
docker network inspect eduflow_eduflow-network

# 检查配置文件
cat configs/config.yaml | grep database
```

#### 7.1.3 Redis连接失败
```bash
# 检查Redis服务
docker-compose exec redis redis-cli ping

# 检查Redis配置
docker-compose exec redis redis-cli config get "*"
```

#### 7.1.4 MinIO访问失败
```bash
# 检查MinIO服务
curl http://localhost:9000/minio/health/live

# 检查MinIO配置
docker-compose exec minio mc admin info local
```

### 7.2 性能问题

#### 7.2.1 响应缓慢
- 检查数据库查询性能
- 优化Redis缓存策略
- 调整连接池配置
- 监控系统资源

#### 7.2.2 内存不足
- 增加系统内存
- 优化应用内存使用
- 调整JVM参数（如适用）
- 配置内存监控

### 7.3 日志分析

#### 7.3.1 应用日志
```bash
# 查看应用日志
docker-compose logs api

# 实时查看日志
docker-compose logs -f api

# 查看错误日志
docker-compose logs api | grep ERROR
```

#### 7.3.2 系统日志
```bash
# 查看系统日志
journalctl -u docker

# 查看磁盘使用
df -h

# 查看内存使用
free -h
```

### 7.4 备份和恢复

#### 7.4.1 数据库备份
```bash
# 备份数据库
docker-compose exec mysql mysqldump -u root -p eduflow > backup.sql

# 恢复数据库
docker-compose exec -i mysql mysql -u root -p eduflow < backup.sql
```

#### 7.4.2 文件备份
```bash
# 备份数据卷
docker run --rm -v eduflow_mysql_data:/data -v $(pwd):/backup alpine tar czf /backup/mysql_data.tar.gz -C /data .

# 备份应用文件
tar czf backup_app.tar.gz configs/ logs/
```

---

## 8. 生产环境优化

### 8.1 安全加固
- 更新默认密码
- 配置防火墙规则
- 启用SSL/TLS
- 定期安全扫描

### 8.2 性能优化
- 配置负载均衡
- 启用Gzip压缩
- 优化数据库查询
- 配置CDN加速

### 8.3 监控告警
- 配置系统监控
- 设置告警规则
- 配置日志收集
- 建立运维流程

---

*本手册版本: v1.0*  
*最后更新: 2024年*  
*技术支持: support@eduflow.com*