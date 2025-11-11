[根目录](../../CLAUDE.md) > [pkg](../) > **logger**

# logger - 日志系统模块

> 更新时间：2025-11-11 14:11:43

## 模块职责

logger 模块提供了统一的日志记录功能，主要包括：
- 多级别日志支持（DEBUG、INFO、WARN、ERROR）
- 格式化日志输出
- 日志级别动态配置
- 错误和致命错误的特殊处理

## 入口与启动

### 日志级别设置
```go
// 通过环境变量或配置文件设置
// export WATCHDUCKER_LOG_LEVEL=DEBUG

// 自动从配置中读取日志级别
level := config.Get().LogLevel()
logger.SetLevel(logger.DEBUG)
```

### 基础使用
```go
logger.Debug("调试信息: %s", debugData)
logger.Info("常规信息")
logger.Warn("警告信息")
logger.Error("错误信息: %v", err)
logger.Fatal("致命错误: %v", err) // 会退出程序
```

## 对外接口

### 日志级别枚举
- `DEBUG` - 调试信息，最详细
- `INFO` - 常规信息，默认级别
- `WARN` - 警告信息
- `ERROR` - 错误信息

### 主要导出函数
- `SetLevel(level Level)` - 设置全局日志级别
- `Debug(format string, args ...interface{})` - 调试级别日志
- `Info(format string, args ...interface{})` - 信息级别日志
- `Warn(format string, args ...interface{})` - 警告级别日志
- `Error(format string, args ...interface{})` - 错误级别日志
- `Fatal(format string, args ...interface{})` - 致命错误日志（退出程序）

## 关键依赖与配置

### 标准库依赖
- `fmt`: 格式化输出
- `io`: 输出流控制
- `os`: 标准输出和程序退出
- `time`: 时间戳格式化

### 输出配置
- 默认输出到 `os.Stdout`
- 支持自定义输出流
- 自动添加时间戳和日志级别前缀
- 格式：`[级别] 时间 消息`

### 性能考虑
- 避免在循环中频繁创建日志消息
- 使用格式化字符串减少字符串拼接
- 支持 level-based filtering 减少不必要输出

## 数据模型

### Level 枚举
```go
type Level int

const (
    DEBUG Level = iota
    INFO
    WARN
    ERROR
)
```

### 全局状态
- `currentLevel`: 当前日志级别，默认为 INFO
- `writer`: 输出流，默认为 os.Stdout
- 使用标准库的并发安全机制

## 测试与质量

> ⚠️ **测试状态**: 无单元测试

### 测试难点
- 测试输出内容需要捕获标准输出
- 测试 Fatal 函数需要特殊的测试框架
- 并发写入的安全性测试

### 测试建议
- 使用 bytes.Buffer 捕获输出进行验证
- 测试不同日志级别的过滤效果
- 测试格式化字符串的正确性
- 测试时间戳格式和时区

## 常见问题 (FAQ)

### Q: 如何禁用调试日志？
A: 通过环境变量设置 `WATCHDUCKER_LOG_LEVEL=INFO` 或更高

### Q: Fatal 和 Error 有什么区别？
A: Error 记录错误但程序继续运行，Fatal 记录错误并退出程序

### Q: 可以自定义日志格式吗？
A: 当前使用固定格式，可以扩展支持自定义格式化器

### Q: 日志输出到文件如何配置？
A: 当前输出到标准输出，可以通过重定向或扩展文件输出功能

## 相关文件清单

- `logger.go` - 日志系统的主要实现

## 变更记录 (Changelog)

### 2025-11-11 14:11:43
- 初始化模块文档
- 记录日志级别和使用方式
- 梳理配置和输出机制
- 识别测试需求