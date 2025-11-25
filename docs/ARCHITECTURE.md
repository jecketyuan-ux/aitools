# EduFlow 系统架构设计文档

## 1. 项目概述

### 1.1 项目简介
EduFlow是一个基于Go语言开发的企业级在线培训平台，提供完整的课程管理、用户管理、资源管理和学习进度跟踪功能。

### 1.2 技术栈
- **后端语言**: Go 1.24+
- **Web框架**: Gin
- **ORM框架**: GORM
- **数据库**: MySQL 8.0+
- **缓存**: Redis 7+
- **文件存储**: MinIO (S3兼容)
- **认证**: JWT (golang-jwt)
- **配置管理**: Viper
- **容器化**: Docker & Docker Compose

## 2. 系统架构

### 2.1 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                        前端应用层                              │
├─────────────────────┬───────────────────────────────────────┤
│   管理后台           │         用户学习端                      │
│   (Admin Portal)    │      (Learning Portal)                │
└─────────────────────┴───────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      API网关层                               │
│                   (Gin Router)                             │
│  ┌─────────────────┬─────────────────┬─────────────────────┐ │
│  │   认证中间件      │    限流中间件      │     CORS中间件      │ │
│  │  Auth Middleware │ Rate Limiting   │  CORS Middleware   │ │
│  └─────────────────┴─────────────────┴─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      业务处理层                               │
│                   (Handler Layer)                           │
│  ┌─────────────────┬─────────────────┬─────────────────────┐ │
│  │   管理端处理器    │    用户端处理器    │     公共处理器       │ │
│  │ Backend Handler  │ Frontend Handler │  Common Handler    │ │
│  └─────────────────┴─────────────────┴─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      业务逻辑层                               │
│                   (Service Layer)                           │
│  ┌─────────────────┬─────────────────┬─────────────────────┐ │
│  │   认证服务       │    用户服务       │     课程服务         │ │
│  │  Auth Service   │  User Service   │  Course Service    │ │
│  └─────────────────┴─────────────────┴─────────────────────┘ │
│  ┌─────────────────┬─────────────────┬─────────────────────┐ │
│  │   资源服务       │    部门服务       │     权限服务         │ │
│  │Resource Service │Department Svc   │ Permission Service  │ │
│  └─────────────────┴─────────────────┴─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      数据访问层                               │
│                 (Repository Layer)                          │
│  ┌─────────────────┬─────────────────┬─────────────────────┐ │
│  │   用户仓储       │    课程仓储       │     资源仓储         │ │
│  │ User Repository │Course Repository │Resource Repository │ │
│  └─────────────────┴─────────────────┴─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       数据存储层                              │
│  ┌─────────────────┬─────────────────┬─────────────────────┐ │
│  │     MySQL       │      Redis      │       MinIO         │ │
│  │   (主数据库)      │    (缓存/会话)    │     (文件存储)       │ │
│  └─────────────────┴─────────────────┴─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 数据库设计

#### 2.2.1 核心实体关系图

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     User        │    │   Department    │    │   AdminUser     │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ id (PK)         │    │ id (PK)         │    │ id (PK)         │
│ email           │    │ name            │    │ email           │
│ name            │    │ parent_id       │    │ name            │
│ password        │    │ parent_chain    │    │ password        │
│ salt            │    │ sort            │    │ role_id         │
│ is_active       │    │ created_at      │    │ is_active       │
│ created_at      │    │ updated_at      │    │ created_at      │
│ updated_at      │    └─────────────────┘    │ updated_at      │
└─────────────────┘                            └─────────────────┘
         │                                               │
         │                                               │
    ┌────┴─────┐                                  ┌─────┴─────┐
    │UserDept  │                                  │AdminRole  │
    ├──────────┤                                  ├───────────┤
    │ user_id  │                                  │ id (PK)   │
    │ dep_id   │                                  │ name      │
    └──────────┘                                  │ permission│
                                                   └───────────┘
         │                                               │
         └───────────────────────┬───────────────────────┘
                                 │
                                 ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Course      │    │    Resource     │    │    Category     │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ id (PK)         │    │ id (PK)         │    │ id (PK)         │
│ title           │    │ name            │    │ name            │
│ thumb           │    │ type            │    │ sort            │
│ short_desc      │    │ rid             │    │ created_at      │
│ class_hour      │    │ size            │    │ updated_at      │
│ is_required     │    │ created_at      │    └─────────────────┘
│ is_show         │    │ updated_at      │            │
│ created_at      │    └─────────────────┘            │
│ updated_at      │            │                       │
└─────────────────┘            │                       │
         │                     │                       │
    ┌────┴─────┐               │                       │
    │CourseCh  │               │                       │
    ├──────────┤               │                       │
    │ id (PK)  │               │                       │
    │ course_id│               │                       │
    │ name     │               │                       │
    │ sort     │               │                       │
    └──────────┘               │                       │
         │                     │                       │
    ┌────┴─────┐               │                       │
    │CourseHr  │               │                       │
    ├──────────┤               │                       │
    │ id (PK)  │               │                       │
    │ course_id│               │                       │
    │ chapter_id│              │                       │
    │ title    │               │                       │
    │ type     │               │                       │
    │ rid      │               │                       │
    │ duration │               │                       │
    └──────────┘               │                       │
                                │                       │
                                └───────────────────────┘
```

## 3. 模块设计

### 3.1 用户管理模块

#### 3.1.1 功能特性
- 用户注册/登录/登出
- 用户信息管理
- 部门管理
- 用户权限控制
- LDAP集成支持

#### 3.1.2 核心组件
- **UserRepository**: 用户数据访问
- **UserService**: 用户业务逻辑
- **AuthHandler**: 认证处理
- **UserHandler**: 用户管理处理

#### 3.1.3 数据模型
```go
type User struct {
    ID            int        `json:"id"`
    Email         string     `json:"email"`
    Name          string     `json:"name"`
    Password      string     `json:"-"`
    Salt          string     `json:"-"`
    IsActive      int        `json:"is_active"`
    CreatedAt     time.Time  `json:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at"`
}
```

### 3.2 课程管理模块

#### 3.2.1 功能特性
- 课程创建/编辑/删除
- 章节管理
- 课时管理
- 课程分类
- 学习进度跟踪

#### 3.2.2 核心组件
- **CourseRepository**: 课程数据访问
- **CourseService**: 课程业务逻辑
- **CourseHandler**: 课程管理处理

#### 3.2.3 数据模型
```go
type Course struct {
    ID          int        `json:"id"`
    Title       string     `json:"title"`
    Thumb       *int       `json:"thumb"`
    Charge      int        `json:"charge"`
    ShortDesc   string     `json:"short_desc"`
    IsRequired  int        `json:"is_required"`
    ClassHour   int        `json:"class_hour"`
    IsShow      int        `json:"is_show"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}
```

### 3.3 资源管理模块

#### 3.3.1 功能特性
- 视频上传/管理
- 图片上传/管理
- 文件存储
- 资源分类
- 附件管理

#### 3.3.2 核心组件
- **ResourceRepository**: 资源数据访问
- **ResourceService**: 资源业务逻辑
- **MinIOStorage**: 文件存储处理
- **ResourceHandler**: 资源管理处理

## 4. API设计

### 4.1 API架构
- RESTful API设计
- 统一响应格式
- JWT认证
- 请求限流
- CORS支持

### 4.2 API路由结构

#### 4.2.1 管理端API (/backend/v1)
```
POST   /auth/login          # 管理员登录
POST   /auth/logout         # 管理员登出
GET    /auth/detail         # 获取管理员信息

GET    /user                # 获取用户列表
POST   /user                # 创建用户
GET    /user/:id            # 获取用户详情
PUT    /user/:id            # 更新用户
DELETE /user/:id            # 删除用户

GET    /course              # 获取课程列表
POST   /course              # 创建课程
GET    /course/:id          # 获取课程详情
PUT    /course/:id          # 更新课程
DELETE /course/:id          # 删除课程

GET    /resource            # 获取资源列表
POST   /resource/video/upload # 上传视频
POST   /resource/image/upload # 上传图片
DELETE /resource/:id        # 删除资源
```

#### 4.2.2 用户端API (/api/v1)
```
POST   /auth/login          # 用户登录
POST   /auth/register       # 用户注册
POST   /auth/logout         # 用户登出
GET    /auth/detail         # 获取用户信息

GET    /courses             # 获取课程列表
GET    /course/:id          # 获取课程详情
```

### 4.3 响应格式
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    // 具体数据
  }
}
```

## 5. 安全设计

### 5.1 认证机制
- JWT Token认证
- Token过期时间控制
- Redis存储Token黑名单
- 多端登录管理

### 5.2 权限控制
- 基于角色的权限控制(RBAC)
- 管理员权限分级
- API接口权限验证
- 数据访问权限控制

### 5.3 安全防护
- 密码加密存储(bcrypt)
- SQL注入防护(GORM)
- XSS防护
- CSRF防护
- 请求限流
- CORS配置

## 6. 性能设计

### 6.1 数据库优化
- 索引优化
- 查询优化
- 连接池配置
- 读写分离支持

### 6.2 缓存策略
- Redis缓存热点数据
- 用户会话缓存
- 课程信息缓存
- 资源信息缓存

### 6.3 文件存储
- MinIO分布式存储
- CDN加速支持
- 文件压缩优化
- 缩略图生成

## 7. 部署架构

### 7.1 容器化部署
- Docker容器化
- Docker Compose编排
- 多环境配置
- 健康检查配置

### 7.2 服务依赖
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   EduFlow   │────│   MySQL     │    │    Redis    │
│    API      │    │  Database   │    │    Cache    │
└─────────────┘    └─────────────┘    └─────────────┘
         │                                      │
         │                                      │
         └──────────────┬───────────────────────┘
                        │
                        ▼
                ┌─────────────┐
                │    MinIO    │
                │  Storage    │
                └─────────────┘
```

## 8. 监控与运维

### 8.1 日志管理
- 结构化日志
- 日志级别控制
- 日志轮转
- 错误日志追踪

### 8.2 监控指标
- 应用性能监控
- 数据库性能监控
- Redis监控
- 系统资源监控

### 8.3 健康检查
- HTTP健康检查接口
- 数据库连接检查
- Redis连接检查
- 依赖服务检查

## 9. 扩展性设计

### 9.1 水平扩展
- 无状态应用设计
- 负载均衡支持
- 数据库分片支持
- 缓存集群支持

### 9.2 功能扩展
- 插件化架构
- 微服务迁移准备
- API版本管理
- 配置中心集成

## 10. 技术选型说明

### 10.1 Go语言优势
- 高性能并发
- 简洁语法
- 丰富标准库
- 快速编译
- 跨平台支持

### 10.2 Gin框架优势
- 高性能路由
- 中间件支持
- JSON处理优秀
- 文档完善
- 社区活跃

### 10.3 GORM优势
- ORM功能完整
- 数据库迁移
- 关联查询
- 事务支持
- 多数据库支持

这份架构设计文档为EduFlow项目提供了全面的技术架构指导，涵盖了系统的各个层面，为开发、部署和维护提供了详细的参考。