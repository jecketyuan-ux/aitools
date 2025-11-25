# EduFlow 流程设计图

## 1. 系统整体流程图

```mermaid
graph TB
    A[用户访问] --> B{用户类型}
    B -->|管理员| C[管理后台]
    B -->|学员| D[学习平台]
    
    C --> E[用户管理]
    C --> F[课程管理]
    C --> G[资源管理]
    C --> H[部门管理]
    
    D --> I[课程学习]
    D --> J[个人中心]
    D --> K[学习统计]
    
    E --> L[(MySQL数据库)]
    F --> L
    G --> M[(MinIO存储)]
    H --> L
    I --> L
    J --> L
    K --> L
    
    L --> N[(Redis缓存)]
    M --> O[文件服务]
    N --> P[会话管理]
```

## 2. 用户认证流程图

```mermaid
sequenceDiagram
    participant User as 用户
    participant Frontend as 前端
    participant API as API服务
    participant Auth as 认证服务
    participant DB as 数据库
    participant Redis as Redis缓存
    
    User->>Frontend: 输入登录信息
    Frontend->>API: POST /auth/login
    API->>Auth: 验证用户凭据
    Auth->>DB: 查询用户信息
    DB-->>Auth: 返回用户数据
    Auth->>Auth: 验证密码
    Auth->>Auth: 生成JWT Token
    Auth->>Redis: 存储Token信息
    Auth-->>API: 返回Token和用户信息
    API-->>Frontend: 返回登录结果
    Frontend-->>User: 显示登录成功
```

## 3. 课程学习流程图

```mermaid
stateDiagram-v2
    [*] --> 选择课程
    选择课程 --> 查看课程详情
    查看课程详情 --> 开始学习
    开始学习 --> 学习章节
    学习章节 --> 学习课时
    学习课时 --> 完成课时
    完成课时 --> 下一课时?
    下一课时? -->|是| 学习课时
    下一课时? -->|否| 完成章节
    完成章节 --> 下一章节?
    下一章节? -->|是| 学习章节
    下一章节? -->|否| 完成课程
    完成课程 --> 获得证书
    获得证书 --> [*]
    
    学习课时 --> 暂停学习
    暂停学习 --> 继续学习
    继续学习 --> 学习课时
```

## 4. 课程管理流程图

```mermaid
flowchart TD
    A[管理员登录] --> B[进入课程管理]
    B --> C{操作类型}
    
    C -->|创建课程| D[填写课程基本信息]
    C -->|编辑课程| E[选择已有课程]
    C -->|删除课程| F[选择要删除的课程]
    
    D --> G[添加章节]
    E --> G
    G --> H[添加课时]
    H --> I[上传资源]
    I --> J[设置发布范围]
    J --> K[发布课程]
    K --> L[课程上线]
    
    F --> M[确认删除]
    M -->|确认| N[删除课程数据]
    M -->|取消| B
    N --> O[更新课程列表]
    O --> B
```

## 5. 资源上传流程图

```mermaid
graph LR
    A[选择文件] --> B[文件验证]
    B --> C{验证通过?}
    C -->|否| A
    C -->|是| D[上传到MinIO]
    D --> E[生成资源记录]
    E --> F[保存到数据库]
    F --> G[返回资源ID]
    G --> H[关联到课程]
```

## 6. 用户注册流程图

```mermaid
sequenceDiagram
    participant User as 用户
    participant Frontend as 前端
    participant API as API服务
    participant UserService as 用户服务
    participant DB as 数据库
    participant Email as 邮件服务
    
    User->>Frontend: 填写注册信息
    Frontend->>API: POST /auth/register
    API->>UserService: 创建用户
    UserService->>DB: 检查邮箱是否存在
    DB-->>UserService: 返回检查结果
    alt 邮箱已存在
        UserService-->>API: 返回错误
        API-->>Frontend: 返回错误信息
        Frontend-->>User: 显示错误
    else 邮箱可用
        UserService->>UserService: 密码加密
        UserService->>DB: 保存用户信息
        DB-->>UserService: 返回用户ID
        UserService-->>API: 返回成功
        API-->>Frontend: 返回注册成功
        Frontend-->>User: 显示注册成功
    end
```

## 7. 权限控制流程图

```mermaid
flowchart TD
    A[用户请求] --> B[JWT验证]
    B --> C{Token有效?}
    C -->|否| D[返回401错误]
    C -->|是| E[解析用户信息]
    E --> F[查询用户权限]
    F --> G{有权限?}
    G -->|否| H[返回403错误]
    G -->|是| I[执行业务逻辑]
    I --> J[返回结果]
```

## 8. 数据同步流程图

```mermaid
graph TB
    A[LDAP服务器] --> B[同步触发]
    B --> C[获取LDAP用户数据]
    C --> D[数据转换]
    D --> E[与本地用户对比]
    E --> F{用户存在?}
    F -->|否| G[创建新用户]
    F -->|是| H[更新用户信息]
    G --> I[保存到数据库]
    H --> I
    I --> J[记录同步日志]
    J --> K[同步完成]
```

## 9. 学习进度跟踪流程图

```mermaid
stateDiagram-v2
    [*] --> 开始学习课时
    开始学习课时 --> 记录开始时间
    记录开始时间 --> 学习进行中
    学习进行中 --> 定期保存进度
    定期保存进度 --> 学习完成?
    学习完成? -->|否| 学习进行中
    学习完成? -->|是| 记录完成时间
    记录完成时间 --> 更新课程进度
    更新课程进度 --> 课程完成?
    课程完成? -->|否| [*]
    课程完成? -->|是| 发放证书
    发放证书 --> [*]
```

## 10. 系统监控流程图

```mermaid
graph LR
    A[应用服务] --> B[日志收集]
    C[数据库] --> B
    D[Redis] --> B
    E[MinIO] --> B
    
    B --> F[日志分析]
    F --> G[指标提取]
    G --> H[监控告警]
    H --> I{异常?}
    I -->|是| J[发送告警]
    I -->|否| K[正常记录]
    J --> L[运维处理]
    K --> M[数据存储]
    L --> M
```

## 11. 部署流程图

```mermaid
flowchart TD
    A[准备部署环境] --> B[安装Docker]
    B --> C[安装Docker Compose]
    C --> D[下载项目代码]
    D --> E[配置环境变量]
    E --> F[构建应用镜像]
    F --> G[启动数据库服务]
    G --> H[等待数据库就绪]
    H --> I[启动应用服务]
    I --> J[健康检查]
    J --> K{服务正常?}
    K -->|否| L[查看日志排查]
    L --> M[重启服务]
    M --> J
    K -->|是| N[部署完成]
    N --> O[访问验证]
```

## 12. 备份恢复流程图

```mermaid
graph TB
    A[定时备份任务] --> B[数据库备份]
    A --> C[文件备份]
    A --> D[配置备份]
    
    B --> E[压缩备份文件]
    C --> E
    D --> E
    E --> F[上传到备份存储]
    F --> G[清理过期备份]
    
    H[恢复请求] --> I[选择备份版本]
    I --> J[下载备份文件]
    J --> K[停止应用服务]
    K --> L[恢复数据库]
    L --> M[恢复文件]
    M --> N[恢复配置]
    N --> O[启动应用服务]
    O --> P[验证恢复结果]
```

## 13. API请求处理流程图

```mermaid
sequenceDiagram
    participant Client as 客户端
    participant LB as 负载均衡器
    participant API as API服务
    participant Middleware as 中间件
    participant Handler as 处理器
    participant Service as 服务层
    participant Repository as 数据层
    participant DB as 数据库
    participant Cache as 缓存
    
    Client->>LB: HTTP请求
    LB->>API: 转发请求
    API->>Middleware: 请求预处理
    Middleware->>Middleware: 认证验证
    Middleware->>Middleware: 限流检查
    Middleware->>Handler: 传递请求
    Handler->>Service: 调用业务逻辑
    Service->>Cache: 检查缓存
    alt 缓存命中
        Cache-->>Service: 返回缓存数据
    else 缓存未命中
        Service->>Repository: 查询数据
        Repository->>DB: SQL查询
        DB-->>Repository: 返回数据
        Repository-->>Service: 返回数据
        Service->>Cache: 更新缓存
    end
    Service-->>Handler: 返回结果
    Handler-->>Middleware: 返回响应
    Middleware-->>API: 返回响应
    API-->>LB: 返回响应
    LB-->>Client: 返回响应
```

## 14. 文件处理流程图

```mermaid
flowchart TD
    A[文件上传请求] --> B[文件类型检查]
    B --> C{类型支持?}
    C -->|否| D[返回错误]
    C -->|是| E[文件大小检查]
    E --> F{大小合规?}
    F -->|否| D
    F -->|是| G[生成文件ID]
    G --> H[上传到MinIO]
    H --> I{上传成功?}
    I -->|否| J[删除临时文件]
    I -->|是| K[创建资源记录]
    K --> L[保存到数据库]
    L --> M[返回资源信息]
    J --> D
```

## 15. 课程发布流程图

```mermaid
stateDiagram-v2
    [*] --> 创建课程
    创建课程 --> 编辑基本信息
    编辑基本信息 --> 添加章节
    添加章节 --> 添加课时
    添加课时 --> 上传资源
    上传资源 --> 预览课程
    预览课程 --> {内容完整?}
    {内容完整?} -->|否| 添加课时
    {内容完整?} -->|是| 设置发布范围
    设置发布范围 --> 确认发布
    确认发布 --> 课程上线
    课程上线 --> 通知学员
    通知学员 --> [*]
```

---

*本文档包含EduFlow系统的主要业务流程和技术流程设计图*  
*版本: v1.0*  
*更新时间: 2024年*