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
