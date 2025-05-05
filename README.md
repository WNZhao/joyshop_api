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
