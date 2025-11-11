[根目录](../../CLAUDE.md) > [internal](../) > **docker**

# docker - Docker API 封装模块

> 更新时间：2025-11-11 14:11:43

## 模块职责

docker 模块是对 Docker Engine API 的高层封装，提供：
- **统一的客户端管理**：创建和管理 Docker 客户端连接
- **容器服务**：容器查询、状态检查、启动、停止等操作
- **镜像服务**：镜像检查、拉取、版本比较等操作
- **资源管理**：统一的资源清理和连接管理

## 入口与启动

### ClientManager 创建
```go
clientManager, err := docker.NewClientManager()
if err != nil {
    return nil, fmt.Errorf("创建 Docker 客户端管理器失败: %w", err)
}
```

### 服务实例化
```go
containerSvc := docker.NewContainerService(clientManager)
imageSvc := docker.NewImageService(clientManager)
```

## 对外接口

### ClientManager
- 管理底层 Docker 客户端连接
- 提供连接检查和关闭方法

### ContainerService 主要方法
- `ListContainersByLabel(labelKey, labelValue)` - 按标签查询容器
- `ListContainersByName(names)` - 按名称查询容器
- `GetAll()` - 获取所有容器信息
- `GetContainerInfo(containerID)` - 获取容器详细信息
- `StartContainer(containerID)` - 启动容器
- `StopContainer(containerID)` - 停止容器
- `RemoveContainer(containerID)` - 删除容器
- `CreateContainer(config)` - 创建新容器

### ImageService 主要方法
- `CheckImageUpdate(imageName)` - 检查镜像是否有更新
- `PullImage(imageName)` - 拉取新镜像
- `GetImageHash(imageName)` - 获取镜像哈希值

## 关键依赖与配置

### 第三方依赖
- `github.com/docker/docker/client`: 官方的 Docker Go 客户端
- `github.com/docker/docker/api/types`: Docker API 类型定义
- `github.com/docker/docker/api/types/container`: 容器相关类型
- `github.com/docker/docker/api/types/filters`: 查询过滤器
- `github.com/docker/docker/api/types/network`: 网络配置类型

### 配置方式
- 使用环境变量自动配置 Docker 连接
- 支持 Docker Socket 连接（默认方式）
- 支持 TCP 连接（通过环境变量配置）

## 数据模型

### 主要数据结构
- 使用官方的 Docker API 类型定义
- 自定义的过滤器和查询条件
- 容器配置和网络配置结构体

## 测试与质量

> ⚠️ **测试状态**: 无单元测试

### 测试难点
- 需要真实的 Docker 环境进行集成测试
- API 调用可能涉及权限和网络问题
- 容器和镜像操作的副作用较大

### 测试建议
- 使用 Docker-in-Docker 进行集成测试
- Mock 测试用于单元测试场景
- 增加连接失败和权限错误的测试用例

## 常见问题 (FAQ)

### Q: 连接 Docker 失败怎么办？
A: 检查 Docker 守护进程是否运行，确保有权限访问 Docker Socket

### Q: 镜像检查的网络超时如何设置？
A: 当前使用默认超时，可能需要增加连接和请求超时配置

### Q: 如何支持私有镜像仓库认证？
A: 需要扩展认证配置，使用 Docker 的认证配置机制

## 相关文件清单

- `client.go` - Docker 客户端管理
- `container.go` - 容器服务实现
- `image.go` - 镜像服务实现

## 变更记录 (Changelog)

### 2025-11-11 14:11:43
- 初始化模块文档
- 记录 Docker API 封装结构
- 识别测试挑战和策略
- 梳理服务接口和依赖