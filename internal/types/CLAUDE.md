[根目录](../../CLAUDE.md) > [internal](../) > **types**

# types - 类型定义模块

> 更新时间：2025-11-11 14:11:43

## 模块职责

types 模块定义了 WatchDucker 项目中的核心数据结构和类型，包括：
- 容器信息的数据结构
- 镜像检查结果的定义
- 批量操作的结果汇总
- 检查模式和回调函数类型

## 入口与启动

本模块不包含可执行代码，仅提供类型定义和常量声明，所有类型都是导出类型，可被其他模块直接引用。

## 对外接口

### 主要类型定义

#### ContainerInfo - 容器信息
```go
type ContainerInfo struct {
    ID     string            `json:"id"`
    Name   string            `json:"name"`
    Image  string            `json:"image"`
    Labels map[string]string `json:"labels"`
    State  string            `json:"state"`
}
```

#### ImageCheckResult - 镜像检查结果
```go
type ImageCheckResult struct {
    Name       string    `json:"name"`
    LocalHash  string    `json:"local_hash"`
    RemoteHash string    `json:"remote_hash"`
    IsUpdated  bool      `json:"is_updated"`
    CheckedAt  time.Time `json:"checked_at"`
    Error      string    `json:"error,omitempty"`
}
```

#### BatchCheckResult - 批量检查结果
```go
type BatchCheckResult struct {
    Containers []ContainerInfo     `json:"containers"`
    Images     []*ImageCheckResult `json:"images"`
    Summary    struct {
        TotalContainers int           `json:"total_containers"`
        TotalImages     int           `json:"total_images"`
        Updated         int           `json:"updated"`
        Failed          int           `json:"failed"`
        UpToDate        int           `json:"up_to_date"`
        Duration        time.Duration `json:"duration"`
    } `json:"summary"`
}
```

### 回调和枚举类型

#### CheckCallback - 检查回调函数
```go
type CheckCallback func(*ImageCheckResult)
```

#### CheckMode - 检查模式枚举
```go
type CheckMode int

const (
    CheckByName  CheckMode = iota // 按名称检查
    CheckByLabel                  // 按标签检查
    CheckAll                      // 检查所有容器
)
```

## 关键依赖与配置

### 标准库依赖
- `time`: 时间相关类型
- `encoding/json`: JSON 序列化标签

### 与其他模块的关系
- 本模块是数据定义层，被所有业务模块引用
- 不依赖任何其他内部模块
- 提供跨模块数据交换的标准格式

## 数据模型

### 数据结构关系
```
BatchCheckResult        ContainerInfo        ImageCheckResult
├── Containers [] ──────→ (1:1 对应)          ───────→ (1:1 对应)
├── Images []* ──────────────────────────────→ (引用)
└── Summary
```

### JSON 序列化支持
所有主要结构体都支持 JSON 序列化，便于：
- 日志输出和调试
- 外部 API 集成
- 配置导出和导入

## 测试与质量

> ⚠️ **测试状态**: 无测试代码

### 验证必要性
虽然类型定义相对简单，但仍建议：
- 类型零值和默认值的验证
- JSON 序列化/反序列化测试
- 边界条件的验证

### 测试建议
- 简单的单元测试验证结构体初始化
- JSON 标签的正确性测试
- 枚举类型的完整性测试

## 常见问题 (FAQ)

### Q: 为什么使用指针类型的 ImageCheckResult？
A: 允许为空值表示检查失败或未执行的情况

### Q: Duration 字段的精度是什么？
A: 使用标准库的 time.Duration 类型，精度为纳秒

### Q: 如何扩展新的检查模式？
A: 在 CheckMode 枚举中添加新的常量即可

## 相关文件清单

- `types.go` - 所有类型定义

## 变更记录 (Changelog)

### 2025-11-11 14:11:43
- 初始化模块文档
- 记录所有导出类型定义
- 梳理数据结构关系
- 提供 JSON 序列化说明