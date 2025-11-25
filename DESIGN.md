# PlayEdu Go 版本 - 设计文档

## 1. 项目概述

PlayEdu Go 版本是基于原 PlayEdu (Java + Spring Boot) 的完整功能重构实现，使用 Go 语言开发的企业级在线培训平台。

### 1.1 核心功能

1. **部门管理** - 支持多层级部门组织架构
2. **学员管理** - 学员信息管理、部门关联、批量导入
3. **管理员系统** - 管理员账户、角色权限管理
4. **课程管理** - 课程创建、章节管理、课时安排
5. **视频学习** - 在线视频播放、进度跟踪、防快进
6. **学习统计** - 学习时长统计、进度追踪、完成率分析
7. **资源管理** - 视频、图片等资源的上传、分类、存储
8. **课程附件** - 课程相关文档、资料下载
9. **LDAP集成** - 支持 LDAP/AD 域用户和部门同步
10. **系统配置** - 系统参数配置、日志记录
11. **API限流** - 接口访问频率限制

### 1.2 技术栈

#### 后端
- **语言**: Go 1.21+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0+
- **缓存**: Redis
- **认证**: JWT (golang-jwt)
- **文件存储**: MinIO / S3
- **配置管理**: Viper
- **日志**: Zap
- **数据验证**: validator/v10
- **密码加密**: bcrypt
- **定时任务**: cron

#### 前端
- **PC端**: React 18
- **移动端**: React 18 (H5)
- **管理后台**: React 18

## 2. 系统架构

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                         Client Layer                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Admin Portal │  │   PC Client  │  │   H5 Client  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │ HTTP/HTTPS
┌─────────────────────────────────────────────────────────────┐
│                         API Gateway                          │
│                   (Nginx / Load Balancer)                    │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                      Application Layer                       │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                 Gin HTTP Server                      │   │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐    │   │
│  │  │ Middleware │  │  Handlers  │  │   Router   │    │   │
│  │  └────────────┘  └────────────┘  └────────────┘    │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                 Service Layer                        │   │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐    │   │
│  │  │   User     │  │   Course   │  │  Resource  │    │   │
│  │  │  Service   │  │  Service   │  │  Service   │    │   │
│  │  └────────────┘  └────────────┘  └────────────┘    │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │               Repository Layer                       │   │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐    │   │
│  │  │    User    │  │   Course   │  │  Resource  │    │   │
│  │  │Repository  │  │ Repository │  │ Repository │    │   │
│  │  └────────────┘  └────────────┘  └────────────┘    │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                      Infrastructure Layer                    │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐           │
│  │   MySQL    │  │   Redis    │  │   MinIO    │           │
│  └────────────┘  └────────────┘  └────────────┘           │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 项目目录结构

```
playedu-go/
├── cmd/
│   └── api/
│       └── main.go                 # 应用入口
├── internal/
│   ├── config/                     # 配置管理
│   │   └── config.go
│   ├── domain/                     # 领域模型
│   │   ├── user.go
│   │   ├── department.go
│   │   ├── course.go
│   │   ├── resource.go
│   │   └── ...
│   ├── repository/                 # 数据访问层
│   │   ├── user_repository.go
│   │   ├── department_repository.go
│   │   ├── course_repository.go
│   │   └── ...
│   ├── service/                    # 业务逻辑层
│   │   ├── user_service.go
│   │   ├── course_service.go
│   │   ├── resource_service.go
│   │   ├── auth_service.go
│   │   └── ...
│   ├── handler/                    # HTTP处理器
│   │   ├── backend/               # 后台管理API
│   │   │   ├── user_handler.go
│   │   │   ├── course_handler.go
│   │   │   └── ...
│   │   └── frontend/              # 前台学员API
│   │       ├── user_handler.go
│   │       ├── course_handler.go
│   │       └── ...
│   ├── middleware/                 # 中间件
│   │   ├── auth.go                # 认证
│   │   ├── permission.go          # 权限
│   │   ├── logger.go              # 日志
│   │   ├── cors.go                # 跨域
│   │   ├── rate_limit.go          # 限流
│   │   └── recovery.go            # 错误恢复
│   └── pkg/                        # 内部工具包
│       ├── jwt/                   # JWT工具
│       ├── storage/               # 文件存储
│       ├── crypto/                # 加密工具
│       ├── ldap/                  # LDAP客户端
│       └── response/              # 响应格式
├── pkg/                            # 公共工具包
│   ├── utils/
│   └── constants/
├── migrations/                     # 数据库迁移
│   ├── 000001_init_schema.up.sql
│   └── 000001_init_schema.down.sql
├── configs/                        # 配置文件
│   ├── config.yaml
│   └── config.example.yaml
├── scripts/                        # 脚本文件
│   ├── init_db.sh
│   └── build.sh
├── docker/                         # Docker配置
│   ├── Dockerfile
│   └── docker-compose.yml
├── docs/                           # 文档
│   ├── api.md
│   └── deployment.md
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 3. 数据库设计

### 3.1 核心表结构

#### 用户相关

**users - 学员表**
```sql
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL COMMENT '邮箱',
  `name` varchar(255) NOT NULL COMMENT '姓名',
  `avatar` int DEFAULT NULL COMMENT '头像资源ID',
  `password` varchar(255) DEFAULT NULL COMMENT '密码',
  `salt` varchar(255) DEFAULT NULL COMMENT '盐值',
  `id_card` varchar(255) DEFAULT NULL COMMENT '身份证号',
  `credit` int DEFAULT 0 COMMENT '学分',
  `create_ip` varchar(255) DEFAULT NULL COMMENT '注册IP',
  `create_city` varchar(255) DEFAULT NULL COMMENT '注册城市',
  `is_active` tinyint DEFAULT 1 COMMENT '是否激活',
  `is_lock` tinyint DEFAULT 0 COMMENT '是否锁定',
  `is_verify` tinyint DEFAULT 0 COMMENT '是否实名认证',
  `verify_at` datetime DEFAULT NULL COMMENT '实名认证时间',
  `is_set_password` tinyint DEFAULT 0 COMMENT '是否设置密码',
  `login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_email` (`email`),
  KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学员表';
```

**departments - 部门表**
```sql
CREATE TABLE `departments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '部门名称',
  `parent_id` int NOT NULL DEFAULT 0 COMMENT '父部门ID',
  `parent_chain` varchar(1000) DEFAULT NULL COMMENT '父级链',
  `sort` int DEFAULT 0 COMMENT '排序',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='部门表';
```

**user_department - 用户部门关联表**
```sql
CREATE TABLE `user_department` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL COMMENT '用户ID',
  `department_id` int NOT NULL COMMENT '部门ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_dep` (`user_id`, `department_id`),
  KEY `idx_department_id` (`department_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户部门关联表';
```

#### 管理员相关

**admin_users - 管理员表**
```sql
CREATE TABLE `admin_users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '姓名',
  `email` varchar(255) NOT NULL COMMENT '邮箱',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `salt` varchar(255) NOT NULL COMMENT '盐值',
  `login_ip` varchar(255) DEFAULT NULL COMMENT '登录IP',
  `login_at` datetime DEFAULT NULL COMMENT '登录时间',
  `is_ban_login` tinyint DEFAULT 0 COMMENT '是否禁止登录',
  `login_times` int DEFAULT 0 COMMENT '登录次数',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员表';
```

**admin_roles - 管理员角色表**
```sql
CREATE TABLE `admin_roles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '角色名称',
  `slug` varchar(255) NOT NULL COMMENT '角色标识',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员角色表';
```

**admin_permissions - 管理员权限表**
```sql
CREATE TABLE `admin_permissions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `type` varchar(50) NOT NULL COMMENT '权限类型',
  `group_name` varchar(255) NOT NULL COMMENT '分组名称',
  `sort` int DEFAULT 0 COMMENT '排序',
  `name` varchar(255) NOT NULL COMMENT '权限名称',
  `slug` varchar(255) NOT NULL COMMENT '权限标识',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员权限表';
```

**admin_role_permission - 角色权限关联表**
```sql
CREATE TABLE `admin_role_permission` (
  `id` int NOT NULL AUTO_INCREMENT,
  `role_id` int NOT NULL COMMENT '角色ID',
  `permission_id` int NOT NULL COMMENT '权限ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_perm` (`role_id`, `permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';
```

**admin_user_role - 管理员角色关联表**
```sql
CREATE TABLE `admin_user_role` (
  `id` int NOT NULL AUTO_INCREMENT,
  `admin_id` int NOT NULL COMMENT '管理员ID',
  `role_id` int NOT NULL COMMENT '角色ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_admin_role` (`admin_id`, `role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员角色关联表';
```

#### 课程相关

**courses - 课程表**
```sql
CREATE TABLE `courses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(500) NOT NULL COMMENT '课程标题',
  `thumb` int DEFAULT NULL COMMENT '课程封面',
  `charge` int DEFAULT 0 COMMENT '课程价格(分)',
  `short_desc` text COMMENT '课程简介',
  `is_required` tinyint DEFAULT 0 COMMENT '是否必修',
  `class_hour` int DEFAULT 0 COMMENT '课时数',
  `is_show` tinyint DEFAULT 1 COMMENT '是否显示',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `sort_at` datetime DEFAULT NULL COMMENT '排序时间',
  `published_at` datetime DEFAULT NULL COMMENT '发布时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_title` (`title`(255)),
  KEY `idx_is_show` (`is_show`),
  KEY `idx_sort_at` (`sort_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程表';
```

**course_chapters - 课程章节表**
```sql
CREATE TABLE `course_chapters` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL COMMENT '课程ID',
  `name` varchar(500) NOT NULL COMMENT '章节名称',
  `sort` int DEFAULT 0 COMMENT '排序',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_course_id` (`course_id`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程章节表';
```

**course_hours - 课时表**
```sql
CREATE TABLE `course_hours` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL COMMENT '课程ID',
  `chapter_id` int DEFAULT 0 COMMENT '章节ID',
  `sort` int DEFAULT 0 COMMENT '排序',
  `title` varchar(500) NOT NULL COMMENT '课时标题',
  `type` varchar(50) NOT NULL COMMENT '课时类型(video)',
  `rid` int NOT NULL COMMENT '资源ID',
  `duration` int DEFAULT 0 COMMENT '时长(秒)',
  `published_at` datetime DEFAULT NULL COMMENT '发布时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_course_id` (`course_id`),
  KEY `idx_chapter_id` (`chapter_id`),
  KEY `idx_rid` (`rid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课时表';
```

**course_categories - 课程分类关联表**
```sql
CREATE TABLE `course_categories` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL COMMENT '课程ID',
  `category_id` int NOT NULL COMMENT '分类ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_course_cat` (`course_id`, `category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程分类关联表';
```

**course_attachments - 课程附件表**
```sql
CREATE TABLE `course_attachments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL COMMENT '课程ID',
  `sort` int DEFAULT 0 COMMENT '排序',
  `title` varchar(500) NOT NULL COMMENT '附件标题',
  `type` varchar(50) NOT NULL COMMENT '附件类型',
  `rid` int NOT NULL COMMENT '资源ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_course_id` (`course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程附件表';
```

**course_attachment_download_log - 附件下载记录表**
```sql
CREATE TABLE `course_attachment_download_log` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL COMMENT '用户ID',
  `course_id` int NOT NULL COMMENT '课程ID',
  `title` varchar(500) NOT NULL COMMENT '附件标题',
  `attachment_id` int NOT NULL COMMENT '附件ID',
  `rid` int NOT NULL COMMENT '资源ID',
  `ip` varchar(255) DEFAULT NULL COMMENT 'IP地址',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_course_id` (`course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='附件下载记录表';
```

**course_department_user - 课程可见部门用户表**
```sql
CREATE TABLE `course_department_user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL COMMENT '课程ID',
  `dep_id` int NOT NULL COMMENT '部门ID',
  `user_id` int NOT NULL COMMENT '用户ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_course_id` (`course_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程可见部门用户表';
```

#### 学习记录相关

**user_course_records - 用户课程学习记录表**
```sql
CREATE TABLE `user_course_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL COMMENT '用户ID',
  `course_id` int NOT NULL COMMENT '课程ID',
  `hour_count` int DEFAULT 0 COMMENT '课时数量',
  `finished_count` int DEFAULT 0 COMMENT '已完成课时数',
  `progress` int DEFAULT 0 COMMENT '学习进度(0-100)',
  `is_finished` tinyint DEFAULT 0 COMMENT '是否完成',
  `finished_at` datetime DEFAULT NULL COMMENT '完成时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_course` (`user_id`, `course_id`),
  KEY `idx_course_id` (`course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户课程学习记录表';
```

**user_course_hour_records - 用户课时学习记录表**
```sql
CREATE TABLE `user_course_hour_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL COMMENT '用户ID',
  `course_id` int NOT NULL COMMENT '课程ID',
  `hour_id` int NOT NULL COMMENT '课时ID',
  `total_duration` int DEFAULT 0 COMMENT '总时长(秒)',
  `finished_duration` int DEFAULT 0 COMMENT '已学习时长(秒)',
  `real_duration` int DEFAULT 0 COMMENT '实际学习时长(秒)',
  `is_finished` tinyint DEFAULT 0 COMMENT '是否完成',
  `finished_at` datetime DEFAULT NULL COMMENT '完成时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_hour` (`user_id`, `hour_id`),
  KEY `idx_course_id` (`course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户课时学习记录表';
```

**user_learn_duration_records - 用户学习时长记录表**
```sql
CREATE TABLE `user_learn_duration_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL COMMENT '用户ID',
  `course_id` int NOT NULL COMMENT '课程ID',
  `hour_id` int NOT NULL COMMENT '课时ID',
  `duration` int NOT NULL COMMENT '学习时长(秒)',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户学习时长记录表';
```

**user_learn_duration_stats - 用户学习时长统计表**
```sql
CREATE TABLE `user_learn_duration_stats` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL COMMENT '用户ID',
  `duration` int NOT NULL COMMENT '学习时长(秒)',
  `created_date` date NOT NULL COMMENT '统计日期',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_date` (`user_id`, `created_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户学习时长统计表';
```

**user_latest_learn - 用户最近学习记录表**
```sql
CREATE TABLE `user_latest_learn` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL COMMENT '用户ID',
  `course_id` int NOT NULL COMMENT '课程ID',
  `hour_id` int NOT NULL COMMENT '课时ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_course` (`user_id`, `course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户最近学习记录表';
```

#### 资源相关

**resources - 资源表**
```sql
CREATE TABLE `resources` (
  `id` int NOT NULL AUTO_INCREMENT,
  `admin_id` int NOT NULL COMMENT '上传管理员ID',
  `type` varchar(50) NOT NULL COMMENT '资源类型(video/image)',
  `category_id` int DEFAULT 0 COMMENT '分类ID',
  `url` varchar(2000) NOT NULL COMMENT '资源URL',
  `name` varchar(500) NOT NULL COMMENT '资源名称',
  `extension` varchar(50) NOT NULL COMMENT '扩展名',
  `size` bigint DEFAULT 0 COMMENT '文件大小(字节)',
  `disk` varchar(50) NOT NULL COMMENT '存储磁盘',
  `file_id` varchar(500) DEFAULT NULL COMMENT '文件ID',
  `path` varchar(2000) DEFAULT NULL COMMENT '存储路径',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_name` (`name`(255))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资源表';
```

**resource_categories - 资源分类表**
```sql
CREATE TABLE `resource_categories` (
  `id` int NOT NULL AUTO_INCREMENT,
  `parent_id` int NOT NULL DEFAULT 0 COMMENT '父分类ID',
  `parent_chain` varchar(1000) DEFAULT NULL COMMENT '父级链',
  `name` varchar(255) NOT NULL COMMENT '分类名称',
  `sort` int DEFAULT 0 COMMENT '排序',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资源分类表';
```

**resource_videos - 视频资源扩展表**
```sql
CREATE TABLE `resource_videos` (
  `id` int NOT NULL AUTO_INCREMENT,
  `rid` int NOT NULL COMMENT '资源ID',
  `duration` int DEFAULT 0 COMMENT '时长(秒)',
  `poster` varchar(2000) DEFAULT NULL COMMENT '封面图',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_rid` (`rid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频资源扩展表';
```

#### 分类相关

**categories - 课程分类表**
```sql
CREATE TABLE `categories` (
  `id` int NOT NULL AUTO_INCREMENT,
  `parent_id` int NOT NULL DEFAULT 0 COMMENT '父分类ID',
  `parent_chain` varchar(1000) DEFAULT NULL COMMENT '父级链',
  `name` varchar(255) NOT NULL COMMENT '分类名称',
  `sort` int DEFAULT 0 COMMENT '排序',
  `is_show` tinyint DEFAULT 1 COMMENT '是否显示',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程分类表';
```

#### 系统相关

**app_config - 系统配置表**
```sql
CREATE TABLE `app_config` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `key_name` varchar(255) NOT NULL COMMENT '配置键',
  `key_value` text COMMENT '配置值',
  `is_private` tinyint DEFAULT 0 COMMENT '是否私有',
  `is_hidden` tinyint DEFAULT 0 COMMENT '是否隐藏',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_key` (`key_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';
```

**admin_logs - 管理员操作日志表**
```sql
CREATE TABLE `admin_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `admin_id` int NOT NULL COMMENT '管理员ID',
  `admin_name` varchar(255) NOT NULL COMMENT '管理员名称',
  `module` varchar(255) NOT NULL COMMENT '模块',
  `title` varchar(500) NOT NULL COMMENT '操作标题',
  `opt` varchar(50) NOT NULL COMMENT '操作类型',
  `method` varchar(50) NOT NULL COMMENT '请求方法',
  `url` varchar(2000) NOT NULL COMMENT '请求URL',
  `ip` varchar(255) NOT NULL COMMENT 'IP地址',
  `ip_area` varchar(500) DEFAULT NULL COMMENT 'IP地区',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_admin_id` (`admin_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员操作日志表';
```

**user_login_records - 用户登录记录表**
```sql
CREATE TABLE `user_login_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL COMMENT '用户ID',
  `jti` varchar(255) NOT NULL COMMENT 'JWT ID',
  `ip` varchar(255) DEFAULT NULL COMMENT 'IP地址',
  `ip_area` varchar(500) DEFAULT NULL COMMENT 'IP地区',
  `browser` varchar(500) DEFAULT NULL COMMENT '浏览器',
  `browser_version` varchar(255) DEFAULT NULL COMMENT '浏览器版本',
  `os` varchar(255) DEFAULT NULL COMMENT '操作系统',
  `is_logout` tinyint DEFAULT 0 COMMENT '是否登出',
  `expired_at` datetime NOT NULL COMMENT '过期时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_jti` (`jti`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户登录记录表';
```

**user_upload_image_logs - 用户上传图片日志表**
```sql
CREATE TABLE `user_upload_image_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL COMMENT '用户ID',
  `scene` varchar(255) NOT NULL COMMENT '场景',
  `rid` int NOT NULL COMMENT '资源ID',
  `ip` varchar(255) DEFAULT NULL COMMENT 'IP地址',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户上传图片日志表';
```

#### LDAP相关

**ldap_users - LDAP用户表**
```sql
CREATE TABLE `ldap_users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(255) NOT NULL COMMENT 'LDAP UUID',
  `ou` varchar(500) DEFAULT NULL COMMENT '组织单元',
  `cn` varchar(500) DEFAULT NULL COMMENT '通用名',
  `display_name` varchar(500) DEFAULT NULL COMMENT '显示名',
  `user_id` int DEFAULT NULL COMMENT '关联用户ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='LDAP用户表';
```

**ldap_departments - LDAP部门表**
```sql
CREATE TABLE `ldap_departments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(255) NOT NULL COMMENT 'LDAP UUID',
  `ou` varchar(500) DEFAULT NULL COMMENT '组织单元',
  `name` varchar(500) DEFAULT NULL COMMENT '部门名称',
  `dep_id` int DEFAULT NULL COMMENT '关联部门ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='LDAP部门表';
```

**ldap_sync_records - LDAP同步记录表**
```sql
CREATE TABLE `ldap_sync_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `action` varchar(50) NOT NULL COMMENT '操作类型',
  `start_at` datetime NOT NULL COMMENT '开始时间',
  `end_at` datetime DEFAULT NULL COMMENT '结束时间',
  `status` varchar(50) NOT NULL COMMENT '状态',
  `error_msg` text COMMENT '错误信息',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='LDAP同步记录表';
```

## 4. API 接口设计

### 4.1 后台管理 API (Backend)

#### 认证模块

**POST /backend/v1/auth/login**
- 描述：管理员登录
- 请求：`{ "email": "string", "password": "string" }`
- 响应：`{ "token": "string", "admin": {...} }`

**POST /backend/v1/auth/logout**
- 描述：管理员登出
- 认证：需要

**GET /backend/v1/auth/detail**
- 描述：获取当前管理员信息
- 认证：需要

#### 学员管理

**GET /backend/v1/user**
- 描述：学员列表
- 参数：page, size, name, email, dep_ids
- 认证：需要
- 权限：user

**POST /backend/v1/user**
- 描述：创建学员
- 认证：需要
- 权限：user.store

**PUT /backend/v1/user/:id**
- 描述：更新学员
- 认证：需要
- 权限：user.update

**DELETE /backend/v1/user/:id**
- 描述：删除学员
- 认证：需要
- 权限：user.delete

**GET /backend/v1/user/:id**
- 描述：学员详情
- 认证：需要
- 权限：user

**POST /backend/v1/user/import**
- 描述：批量导入学员
- 认证：需要
- 权限：user.store

**GET /backend/v1/user/:id/learn**
- 描述：学员学习统计
- 认证：需要
- 权限：user

#### 部门管理

**GET /backend/v1/department**
- 描述：部门列表（树形）
- 认证：需要
- 权限：department

**POST /backend/v1/department**
- 描述：创建部门
- 认证：需要
- 权限：department.store

**PUT /backend/v1/department/:id**
- 描述：更新部门
- 认证：需要
- 权限：department.update

**DELETE /backend/v1/department/:id**
- 描述：删除部门
- 认证：需要
- 权限：department.delete

#### 课程管理

**GET /backend/v1/course**
- 描述：课程列表
- 参数：page, size, title, dep_ids, category_ids, is_required
- 认证：需要
- 权限：course

**POST /backend/v1/course**
- 描述：创建课程
- 认证：需要
- 权限：course.store

**PUT /backend/v1/course/:id**
- 描述：更新课程
- 认证：需要
- 权限：course.update

**DELETE /backend/v1/course/:id**
- 描述：删除课程
- 认证：需要
- 权限：course.delete

**GET /backend/v1/course/:id**
- 描述：课程详情
- 认证：需要
- 权限：course

**GET /backend/v1/course/:id/users**
- 描述：课程学员列表及学习进度
- 认证：需要
- 权限：course

#### 章节管理

**GET /backend/v1/course-chapter/:course_id**
- 描述：课程章节列表
- 认证：需要
- 权限：course

**POST /backend/v1/course-chapter**
- 描述：创建章节
- 认证：需要
- 权限：course.store

**PUT /backend/v1/course-chapter/:id**
- 描述：更新章节
- 认证：需要
- 权限：course.update

**DELETE /backend/v1/course-chapter/:id**
- 描述：删除章节
- 认证：需要
- 权限：course.delete

#### 课时管理

**GET /backend/v1/course-hour/:course_id**
- 描述：课程课时列表
- 认证：需要
- 权限：course

**POST /backend/v1/course-hour**
- 描述：创建课时
- 认证：需要
- 权限：course.store

**PUT /backend/v1/course-hour/:id**
- 描述：更新课时
- 认证：需要
- 权限：course.update

**DELETE /backend/v1/course-hour/:id**
- 描述：删除课时
- 认证：需要
- 权限：course.delete

#### 课程附件管理

**GET /backend/v1/course-attachment/:course_id**
- 描述：课程附件列表
- 认证：需要
- 权限：course

**POST /backend/v1/course-attachment**
- 描述：创建附件
- 认证：需要
- 权限：course.store

**PUT /backend/v1/course-attachment/:id**
- 描述：更新附件
- 认证：需要
- 权限：course.update

**DELETE /backend/v1/course-attachment/:id**
- 描述：删除附件
- 认证：需要
- 权限：course.delete

**GET /backend/v1/course-attachment-download-log**
- 描述：附件下载日志
- 认证：需要
- 权限：course

#### 资源管理

**GET /backend/v1/resource**
- 描述：资源列表
- 参数：page, size, type, category_id, name
- 认证：需要
- 权限：resource

**POST /backend/v1/resource/video/upload**
- 描述：上传视频
- 认证：需要
- 权限：resource.upload

**POST /backend/v1/resource/image/upload**
- 描述：上传图片
- 认证：需要
- 权限：resource.upload

**DELETE /backend/v1/resource/:id**
- 描述：删除资源
- 认证：需要
- 权限：resource.delete

#### 资源分类管理

**GET /backend/v1/resource-category**
- 描述：资源分类列表（树形）
- 认证：需要
- 权限：resource-category

**POST /backend/v1/resource-category**
- 描述：创建资源分类
- 认证：需要
- 权限：resource-category.store

**PUT /backend/v1/resource-category/:id**
- 描述：更新资源分类
- 认证：需要
- 权限：resource-category.update

**DELETE /backend/v1/resource-category/:id**
- 描述：删除资源分类
- 认证：需要
- 权限：resource-category.delete

#### 课程分类管理

**GET /backend/v1/category**
- 描述：课程分类列表（树形）
- 认证：需要
- 权限：category

**POST /backend/v1/category**
- 描述：创建课程分类
- 认证：需要
- 权限：category.store

**PUT /backend/v1/category/:id**
- 描述：更新课程分类
- 认证：需要
- 权限：category.update

**DELETE /backend/v1/category/:id**
- 描述：删除课程分类
- 认证：需要
- 权限：category.delete

#### 管理员管理

**GET /backend/v1/admin-user**
- 描述：管理员列表
- 认证：需要
- 权限：admin-user

**POST /backend/v1/admin-user**
- 描述：创建管理员
- 认证：需要
- 权限：admin-user.store

**PUT /backend/v1/admin-user/:id**
- 描述：更新管理员
- 认证：需要
- 权限：admin-user.update

**DELETE /backend/v1/admin-user/:id**
- 描述：删除管理员
- 认证：需要
- 权限：admin-user.delete

#### 角色管理

**GET /backend/v1/admin-role**
- 描述：角色列表
- 认证：需要
- 权限：admin-role

**POST /backend/v1/admin-role**
- 描述：创建角色
- 认证：需要
- 权限：admin-role.store

**PUT /backend/v1/admin-role/:id**
- 描述：更新角色
- 认证：需要
- 权限：admin-role.update

**DELETE /backend/v1/admin-role/:id**
- 描述：删除角色
- 认证：需要
- 权限：admin-role.delete

#### 权限管理

**GET /backend/v1/admin-permission**
- 描述：权限列表
- 认证：需要

#### 系统配置

**GET /backend/v1/system/config**
- 描述：获取系统配置
- 认证：需要
- 权限：system

**PUT /backend/v1/system/config**
- 描述：更新系统配置
- 认证：需要
- 权限：system.config

#### 管理员日志

**GET /backend/v1/admin-log**
- 描述：管理员操作日志
- 认证：需要
- 权限：admin-log

#### Dashboard

**GET /backend/v1/dashboard**
- 描述：仪表板统计数据
- 认证：需要

#### LDAP管理

**GET /backend/v1/ldap/config**
- 描述：获取LDAP配置
- 认证：需要
- 权限：system

**PUT /backend/v1/ldap/config**
- 描述：更新LDAP配置
- 认证：需要
- 权限：system.config

**POST /backend/v1/ldap/sync**
- 描述：手动触发LDAP同步
- 认证：需要
- 权限：ldap.sync

**GET /backend/v1/ldap/sync-records**
- 描述：LDAP同步记录
- 认证：需要
- 权限：ldap

### 4.2 前台学员 API (Frontend)

#### 认证模块

**POST /api/v1/auth/login**
- 描述：学员登录
- 请求：`{ "email": "string", "password": "string" }`
- 响应：`{ "token": "string", "user": {...} }`

**POST /api/v1/auth/register**
- 描述：学员注册
- 请求：`{ "email": "string", "password": "string", "name": "string" }`

**POST /api/v1/auth/logout**
- 描述：学员登出
- 认证：需要

**GET /api/v1/auth/detail**
- 描述：获取当前学员信息
- 认证：需要

**PUT /api/v1/auth/password**
- 描述：修改密码
- 认证：需要

**PUT /api/v1/auth/profile**
- 描述：修改个人信息
- 认证：需要

#### 课程模块

**GET /api/v1/courses**
- 描述：课程列表
- 参数：page, size, category_id, title
- 认证：需要

**GET /api/v1/course/:id**
- 描述：课程详情
- 认证：需要

**GET /api/v1/course/:id/chapters**
- 描述：课程章节列表
- 认证：需要

**GET /api/v1/course/:id/attachments**
- 描述：课程附件列表
- 认证：需要

**GET /api/v1/course/attachment/:id/download**
- 描述：下载课程附件
- 认证：需要

#### 课时模块

**GET /api/v1/hour/:id**
- 描述：课时详情
- 认证：需要

**POST /api/v1/hour/:id/record**
- 描述：提交学习进度
- 认证：需要
- 请求：`{ "duration": 10 }`

#### 学习记录

**GET /api/v1/user/courses**
- 描述：我的课程列表
- 认证：需要

**GET /api/v1/user/course/:id/progress**
- 描述：课程学习进度
- 认证：需要

**GET /api/v1/user/learn/stats**
- 描述：学习统计
- 认证：需要
- 参数：start_date, end_date

**GET /api/v1/user/latest-learn**
- 描述：最近学习记录
- 认证：需要

#### 分类模块

**GET /api/v1/categories**
- 描述：课程分类列表
- 认证：需要

#### 部门模块

**GET /api/v1/departments**
- 描述：部门列表
- 认证：需要

#### 系统信息

**GET /api/v1/system/info**
- 描述：系统信息

### 4.3 响应格式

#### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": {...}
}
```

#### 分页响应
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "size": 10
  }
}
```

#### 错误响应
```json
{
  "code": 1001,
  "message": "error message",
  "data": null
}
```

#### 错误码定义
- 0: 成功
- 1000: 通用错误
- 1001: 参数错误
- 1002: 未授权
- 1003: 无权限
- 1004: 资源不存在
- 1005: 资源已存在
- 1006: 操作失败
- 1007: 限流错误
- 2001: 用户不存在
- 2002: 密码错误
- 2003: 用户已锁定
- 3001: 课程不存在
- 3002: 课程已关闭
- 4001: Token无效
- 4002: Token过期

## 5. 核心功能实现

### 5.1 认证授权

#### JWT Token
- 使用 golang-jwt/jwt 实现
- Token 有效期：15天
- 包含信息：user_id, email, role, jti
- 支持Token刷新机制

#### 权限控制
- RBAC 模型（Role-Based Access Control）
- 管理员角色-权限关联
- 权限检查中间件
- 前台学员基于部门的访问控制

### 5.2 文件存储

#### MinIO 集成
- 支持视频、图片上传
- 自动生成缩略图
- 视频转码（可选）
- CDN 加速支持
- 分片上传大文件

#### 存储结构
```
bucket/
├── videos/
│   ├── 2024/01/
│   │   └── xxx.mp4
├── images/
│   ├── 2024/01/
│   │   └── xxx.jpg
└── attachments/
    └── 2024/01/
        └── xxx.pdf
```

### 5.3 学习进度追踪

#### 进度计算
- 课时完成条件：观看时长 >= 总时长 * 80%
- 课程完成条件：所有课时完成
- 实时更新进度百分比

#### 防快进机制
- 记录实际学习时长
- 检测异常快进行为
- 进度提交频率限制

#### 时长统计
- 按天统计学习时长
- 按课程统计学习时长
- 学习时长排行榜

### 5.4 LDAP 集成

#### 同步功能
- 定时同步（每小时执行）
- 手动触发同步
- 支持增量同步
- 同步部门组织架构
- 同步用户信息

#### 同步流程
1. 连接LDAP服务器
2. 读取部门信息并同步
3. 读取用户信息并同步
4. 建立部门-用户关联
5. 记录同步日志

### 5.5 缓存策略

#### Redis 缓存
- 用户信息缓存（5分钟）
- 课程信息缓存（10分钟）
- 分类树缓存（30分钟）
- Token黑名单（15天）

#### 缓存更新
- 数据变更时主动清除缓存
- 缓存预热机制
- 缓存穿透防护

### 5.6 限流策略

#### API 限流
- IP级别限流：60秒内最多360次请求
- 用户级别限流：60秒内最多180次请求
- 接口级别限流：特定接口自定义限制

#### 实现方式
- 使用 Redis + Lua 实现滑动窗口限流
- 令牌桶算法
- 支持动态调整限流参数

## 6. 部署方案

### 6.1 Docker 部署

#### docker-compose.yml
```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: playedu
    volumes:
      - mysql_data:/var/lib/mysql
      - ./migrations:/docker-entrypoint-initdb.d
    ports:
      - "3306:3306"
  
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
  
  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio_data:/data
  
  api:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - mysql
      - redis
      - minio
    environment:
      DB_HOST: mysql
      DB_PORT: 3306
      DB_NAME: playedu
      DB_USER: root
      DB_PASSWORD: root_password
      REDIS_HOST: redis
      REDIS_PORT: 6379
      MINIO_ENDPOINT: minio:9000
      MINIO_ACCESS_KEY: minioadmin
      MINIO_SECRET_KEY: minioadmin
    volumes:
      - ./configs:/app/configs

volumes:
  mysql_data:
  redis_data:
  minio_data:
```

### 6.2 生产环境部署

#### 服务器要求
- CPU: 4核+
- 内存: 8GB+
- 硬盘: 100GB+ (根据视频存储需求调整)
- 操作系统: Ubuntu 20.04+ / CentOS 7+

#### 部署步骤
1. 安装 Docker 和 Docker Compose
2. 克隆项目代码
3. 配置环境变量
4. 执行数据库迁移
5. 启动服务
6. 配置 Nginx 反向代理
7. 配置 SSL 证书

#### Nginx 配置示例
```nginx
upstream playedu_api {
    server 127.0.0.1:8080;
}

server {
    listen 80;
    server_name api.playedu.com;
    
    location / {
        proxy_pass http://playedu_api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## 7. 开发规范

### 7.1 代码规范

#### Go 代码规范
- 遵循 Go 官方代码规范
- 使用 gofmt 格式化代码
- 使用 golangci-lint 进行代码检查
- 错误处理不能忽略
- 避免全局变量

#### 命名规范
- 包名：小写，简短，有意义
- 文件名：小写，下划线分隔
- 变量名：驼峰命名
- 常量名：大写，下划线分隔
- 接口名：以 er 结尾（如 Reader, Writer）

### 7.2 Git 规范

#### 分支管理
- main: 主分支，生产环境代码
- develop: 开发分支
- feature/*: 功能分支
- bugfix/*: 修复分支
- release/*: 发布分支

#### Commit 规范
```
<type>(<scope>): <subject>

<body>

<footer>
```

**type 类型：**
- feat: 新功能
- fix: 修复
- docs: 文档
- style: 格式
- refactor: 重构
- test: 测试
- chore: 构建/工具

### 7.3 测试规范

#### 单元测试
- 覆盖率要求 > 70%
- 使用 testify 进行断言
- Mock 外部依赖

#### 集成测试
- 测试完整业务流程
- 使用测试数据库

## 8. 监控与运维

### 8.1 日志管理

#### 日志级别
- DEBUG: 调试信息
- INFO: 一般信息
- WARN: 警告信息
- ERROR: 错误信息
- FATAL: 致命错误

#### 日志格式
```json
{
  "time": "2024-01-01T12:00:00Z",
  "level": "info",
  "msg": "user login",
  "user_id": 123,
  "ip": "192.168.1.1"
}
```

### 8.2 性能监控

#### 监控指标
- API 响应时间
- 数据库查询时间
- 缓存命中率
- 错误率
- 并发用户数

#### 监控工具
- Prometheus: 指标采集
- Grafana: 可视化展示
- Alertmanager: 告警管理

## 9. 安全性

### 9.1 数据安全

- 密码使用 bcrypt 加密
- 敏感数据传输使用 HTTPS
- SQL 注入防护（使用参数化查询）
- XSS 防护（输入验证和输出转义）
- CSRF 防护（Token验证）

### 9.2 访问控制

- 基于角色的权限控制
- API 访问频率限制
- 登录失败次数限制
- 异常登录检测

### 9.3 数据备份

- 数据库每日自动备份
- 保留最近 30 天备份
- 异地备份存储
- 定期恢复测试

## 10. 后续优化方向

### 10.1 功能扩展

- 在线考试系统
- 学习任务系统
- 证书系统
- 直播课程
- 社交互动（评论、讨论）
- 移动APP

### 10.2 性能优化

- 数据库读写分离
- 分布式缓存集群
- CDN 加速
- 视频分片上传
- 数据库分表分库

### 10.3 运维优化

- 容器编排（Kubernetes）
- 自动化部署（CI/CD）
- 服务监控告警
- 日志分析系统
- 自动扩缩容

## 11. 总结

PlayEdu Go 版本是一个功能完整、架构清晰、易于扩展的企业级在线培训平台。采用 Go 语言开发，具有高性能、高并发的特点。通过模块化设计，各功能模块职责清晰，便于维护和扩展。

本设计文档涵盖了系统架构、数据库设计、API 接口、核心功能实现、部署方案等各个方面，为后续开发提供了完整的指导。
