# JoyShop API (user-web)

## 项目说明

本项目是 JoyShop 电商系统的用户服务 Web 层，提供用户相关的 HTTP API 接口，通过 gRPC 与用户服务进行通信。

### 项目结构

```
user-web/
├── api/            # API 接口层
├── config/         # 配置文件
├── forms/          # 表单验证
├── global/         # 全局变量
├── initialize/     # 初始化相关
├── middlewares/    # 中间件
├── models/         # 数据模型
├── proto/          # protobuf 文件
├── router/         # 路由配置
└── utils/          # 工具函数
```

### 技术栈

- Web 框架：Gin
- 服务发现：Consul
- 数据存储：Redis
- 日志系统：zap
- 配置管理：viper
- 表单验证：validator
- 认证：JWT

## 功能特性

### 1. 服务注册与发现

- 使用 Consul 进行服务注册与发现
- 服务启动时自动注册到 Consul
- 健康检查机制:
  - 每15秒检查一次服务状态
  - 检查失败15秒后自动注销服务
  - 确保服务列表中的服务都是可用的

### 2. gRPC 客户端

- 使用 Consul 进行服务发现
- 自动连接可用的用户服务实例
- 支持服务健康检查
- 使用 `grpc.NewClient` 创建连接
- 全局客户端管理:
  - 在 `global` 包中定义全局客户端
  - 服务启动时初始化连接
  - 服务关闭时自动关闭连接

### 3. API 接口

- 用户注册
- 用户登录
- 获取用户列表
- 用户信息管理

### 4. 中间件

- JWT 认证
- CORS 跨域支持
- 请求日志记录

## 配置说明

### Consul 配置

```yaml
consul:
  host: "localhost"  # Consul 服务器地址
  port: 8500        # Consul 服务器端口
```

### 服务配置

```yaml
name: "user-web"    # 服务名称
host: "localhost"   # 服务地址
port: 8021         # 服务端口
```

## 开发环境要求

- Go 1.16+
- Consul
- Redis

## 快速开始

1. 启动 Consul

```bash
consul agent -dev
```

2. 启动 Web 服务

```bash
go run main.go
```

## 项目特点

1. 微服务架构
   - 使用 gRPC 进行服务间通信
   - Consul 服务注册与发现
   - 服务健康检查

2. 代码组织
   - 清晰的目录结构
   - 模块化的设计
   - 统一的错误处理

3. 开发体验
   - 完整的日志记录
   - 统一的配置管理
   - 便捷的开发工具

## 最近更新

- 优化 gRPC 客户端初始化逻辑
- 改进服务注册机制
- 统一使用 `grpc.NewClient` API
- 增强健康检查配置

# joyshop_api

## 日志配置说明

### 日志模式切换

项目支持通过环境变量 `APP_ENV` 来切换日志模式：

- 开发环境（默认）：

  ```bash
  # 不设置 APP_ENV 或设置为 development
  go run main.go
  # 或
  APP_ENV=development go run main.go
  ```

- 生产环境：

  ```bash
  # 设置 APP_ENV 为 production
  APP_ENV=production go run main.go
  ```

### 构建说明

在构建时设置环境变量会影响编译后的二进制文件的日志模式：

- 开发环境构建：

  ```bash
  go build
  # 或
  APP_ENV=development go build
  ```

- 生产环境构建：

  ```bash
  APP_ENV=production go build
  ```

注意：构建时设置的环境变量会被编译到二进制文件中，这意味着：

1. 构建时设置 `APP_ENV=production`，运行时会始终使用生产环境日志模式
2. 构建时未设置 `APP_ENV`，运行时可以通过环境变量动态切换日志模式

建议：

- 开发时使用默认模式（不设置 APP_ENV）
- 生产环境构建时使用 `APP_ENV=production go build`
- 如果需要在运行时动态切换日志模式，构建时不要设置 APP_ENV

## 项目配置说明

### 环境变量配置

#### 阿里云短信服务配置

在 `~/.zshrc` 或 `~/.bash_profile` 中添加以下环境变量：

```bash
# Aliyun SMS Credentials
export ALIBABA_CLOUD_ACCESS_KEY_ID=your_access_key_id
export ALIBABA_CLOUD_ACCESS_KEY_SECRET=your_access_key_secret
```

### 配置文件说明

项目包含两个配置文件：

- `config-debug.yaml`: 开发环境配置
- `config-prod.yaml`: 生产环境配置

#### 配置文件结构

```yaml
name: 'user-web-debug'  # 服务名称
port: 8022              # 服务端口
lang: 'zh'              # 语言设置

# 用户服务配置
user_srv:
  host: 'localhost'     # 开发环境使用localhost，生产环境使用服务名
  port: 50051

# JWT配置
jwt:
  signing_key: 'your_signing_key'

# 阿里云短信配置
aliyun_sms:
  sign_name: '阿里云短信测试'    # 短信签名
  template_code: 'SMS_xxx'      # 短信模板代码
  phone_numbers: '13800138000'  # 接收短信的手机号

# Redis配置
redis:
  host: 'localhost'     # 开发环境使用localhost，生产环境使用服务名
  port: 6379
```

### 短信服务说明

1. 短信发送功能使用阿里云短信服务
2. 验证码生成使用本地随机数生成器，默认生成6位数字验证码
3. 短信模板参数格式：`{"code":"123456"}`

### 开发环境配置

- 使用 `config-debug.yaml` 配置文件
- 服务地址使用 `localhost`
- Redis 地址使用 `localhost`

### 生产环境配置

- 使用 `config-prod.yaml` 配置文件
- 服务地址使用服务名（如 `user-srv`）
- Redis 地址使用服务名（如 `redis`）

## 配置文件说明

### 配置文件结构

项目支持通过环境变量 `APP_ENV` 来切换配置文件：

- 开发环境（默认）：使用 `config-debug.yaml`
- 生产环境：使用 `config-prod.yaml`

配置文件示例：

```yaml
name: "user-web-debug"  # 服务名称
port: 8022              # 服务端口
user_srv:               # 用户服务配置
  host: "localhost"     # 用户服务主机
  port: 50051          # 用户服务端口
```

### 配置热更新

项目支持配置文件的动态热更新：

1. 当配置文件发生变化时，系统会自动重新加载配置
2. 配置更新会实时生效，无需重启服务
3. 配置更新会通过日志记录变更信息

### 配置使用说明

1. 开发环境配置：

   ```bash
   # 默认使用 config-debug.yaml
   go run main.go
   ```

2. 生产环境配置：

   ```bash
   # 使用 config-prod.yaml
   APP_ENV=production go run main.go
   ```

3. 配置热更新：
   - 直接修改对应的配置文件
   - 系统会自动检测变更并重新加载
   - 变更信息会记录在日志中

### 配置结构说明

配置文件中的字段说明：

- `name`: 服务名称
- `port`: 服务端口
- `user_srv`: 用户服务配置
  - `host`: 用户服务主机地址
  - `port`: 用户服务端口

### 注意事项

1. 配置文件路径：
   - 配置文件应放在 `user-web` 目录下
   - 开发环境：`config-debug.yaml`
   - 生产环境：`config-prod.yaml`

2. 配置热更新：
   - 仅支持 yaml 格式的配置文件
   - 文件变更会触发自动重载
   - 配置错误会导致服务 panic

3. 环境变量：
   - 可以通过环境变量覆盖配置文件中的值
   - 环境变量优先级高于配置文件

## 表单验证说明

### 自定义验证器

项目支持自定义表单验证器，以手机号验证为例：

1. 在 `initialize/validator.go` 中注册自定义验证器：

```go
// 注册自定义验证器
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
    // 注册手机号验证器
    _ = v.RegisterValidation("mobile", validator.Func(func(fl validator.FieldLevel) bool {
        mobile := fl.Field().String()
        // 使用正则表达式验证手机号
        ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
        return ok
    }))
}
```

2. 在表单结构体中使用验证器：

```go
type PassWordLoginForm struct {
    Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`
    Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
}
```

3. 错误处理：

```go
if err := ctx.ShouldBindJSON(&form); err != nil {
    if utils.HandleValidatorError(ctx, err, "FormName") {
        return
    }
}
```

### 验证器使用说明

1. 内置验证器：
   - `required`: 必填字段
   - `min`: 最小长度
   - `max`: 最大长度
   - `email`: 邮箱格式
   - `url`: URL格式
   - 更多验证器请参考 [validator 文档](https://pkg.go.dev/github.com/go-playground/validator/v10)

2. 自定义验证器：
   - 在 `initialize/validator.go` 中注册
   - 使用 `binding:"custom_validator"` 标签
   - 验证器函数返回 `bool` 类型

3. 错误处理：
   - 使用 `utils.HandleValidatorError` 统一处理
   - 自动移除表单名称前缀
   - 统一错误信息格式

### 错误信息格式

验证失败时返回的错误信息格式：

```json
{
    "code": 400,
    "msg": {
        "mobile": "mobile为必填字段",
        "password": "password为必填字段"
    }
}
```

### 注意事项

1. 验证器注册：
   - 必须在应用启动时注册
   - 验证器名称不能重复
   - 验证器函数必须返回 bool 类型

2. 错误处理：
   - 使用统一的错误处理函数
   - 保持错误信息格式一致
   - 记录详细的错误日志

3. 性能考虑：
   - 验证器函数应该尽量简单
   - 避免在验证器中进行复杂操作
   - 必要时可以使用缓存优化
