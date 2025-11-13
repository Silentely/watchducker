# pkg/notify 模块文档

> 更新时间：2025-11-13 14:26:00

## 模块概述

`notify` 模块是 WatchDucker 的**多平台通知推送服务层**，提供统一的消息通知接口，支持 **15+ 种通知渠道**。该模块负责将容器更新、状态变化等重要事件实时推送到用户指定的通知平台。

## 核心职责

| 职责 | 说明 |
|-----|------|
| **配置管理** | 从 `push.yaml` 配置文件加载各通知平台的认证信息 |
| **多渠道推送** | 支持并行调用多个通知服务 |
| **错误隔离** | 单个推送渠道失败不影响其他渠道 |
| **灵活扩展** | 通过新增函数可轻松支持新的通知平台 |
| **日志记录** | 记录每次推送的成功/失败状态 |

## 代码结构

```
pkg/notify/
├── notify.go                # 主模块文件（520+ 行）
└── CLAUDE.md                # 本文档
```

### 文件说明

| 文件 | 行数 | 职责 |
|------|------|------|
| `notify.go` | 520+ | 完整的配置加载、HTTP/SMTP/加密工具、15+ 推送函数、主逻辑 |

## 支持的通知平台（15 种）

### 即时通讯类

| 平台 | 配置字段 | 认证方式 | 适用场景 |
|------|---------|--------|--------|
| **Telegram** | `telegram.api_url/bot_token/chat_id` | Token + Chat ID | 国际用户首选 |
| **企业微信** | `wecom.*` | Corp ID + Secret + Agent ID | 企业内部通知 |
| **企业微信机器人** | `wecomrobot.url/mobile` | Webhook URL | 群组通知 |
| **QQ (CQHTTP)** | `cqhttp.url/qq` | CQHTTP 服务地址 | QQ 私聊/群聊 |
| **Discord** | `discord.webhook/verify_ssl` | Webhook URL | 开发者社区 |

### 监控/推送服务类

| 平台 | 配置字段 | 认证方式 | 适用场景 |
|------|---------|--------|--------|
| **钉钉** | `dingrobot.webhook/secret` | Webhook + 签名密钥 | 企业通知 |
| **飞书** | `feishubot.webhook` | Webhook URL | 企业通知 |
| **PushDeer** | `pushdeer.api_url/token` | Token | 自托管推送 |
| **Bark** | `bark.api_url/token` | Token | iOS 通知 |
| **Gotify** | `gotify.api_url/token/priority` | Token + 优先级 | 自托管监控 |

### 聚合推送/第三方服务

| 平台 | 配置字段 | 认证方式 | 适用场景 |
|------|---------|--------|--------|
| **Server酱 (FTQQ)** | `ftqq.push_token` | SendKey | 微信通知 |
| **PushPlus** | `pushplus.push_token` | Push Token | 微信/钉钉聚合 |
| **IFTTT** | `ifttt.event/key` | Event + Webhook Key | 自动化工作流 |
| **Qmsg** | `qmsg.key` | 推送 Key | QQ 通知 |
| **邮件** | `smtp.*` | SMTP 服务器配置 | 企业邮箱通知 |

### 自定义服务

| 平台 | 配置字段 | 认证方式 | 适用场景 |
|------|---------|--------|--------|
| **通用 Webhook** | `webhook.webhook_url` | 自定义 URL | 集成任何服务 |

## 配置文件格式

配置文件位置：`push.yaml`（项目根目录）

### 配置示例

```yaml
setting:
  push_server: "telegram,dingrobot,ftqq"  # 多渠道逗号分隔
  log_level: "INFO"                       # DEBUG/INFO/WARN/ERROR

telegram:
  api_url: "api.telegram.org"             # 支持反代地址
  bot_token: "123456:ABC-DEF1234ghIkl-zyx"
  chat_id: "123456789"

dingrobot:
  webhook: "https://oapi.dingtalk.com/robot/send?access_token=xxx"
  secret: "SECxxx"                        # 可选：签名密钥

ftqq:
  push_token: "SCTxxx"

# ... 其他平台配置
```

### 环境变量支持

通过 Viper 库，可使用环境变量覆盖配置：

```bash
export WATCHDUCKER_TELEGRAM_BOT_TOKEN="your_token"
export WATCHDUCKER_TELEGRAM_CHAT_ID="your_chat_id"
```

## 关键数据结构

### Config 结构体

```go
type Config struct {
    Setting struct {
        PushServer string  // 启用的推送服务列表，以逗号分隔
        LogLevel   string  // 日志级别（DEBUG/INFO/WARN/ERROR）
    }

    Telegram   struct { ... }  // Telegram 配置
    Ftqq       struct { ... }  // Server酱配置
    Pushplus   struct { ... }  // PushPlus 配置
    Wecom      struct { ... }  // 企业微信配置
    Dingrobot  struct { ... }  // 钉钉配置
    Feishu     struct { ... }  // 飞书配置
    // ... 其他 13 个平台的配置
}
```

## 公开 API

### Send 函数

```go
// Send 发送通知消息到所有已配置的推送平台
//
// 参数：
//   - title: 消息标题（必填）
//   - msg:   消息内容（必填）
//
// 特性：
//   - 自动加载 push.yaml 配置
//   - 缺少配置文件时优雅跳过
//   - 单个平台失败不影响其他平台
//   - 自动记录每条推送的结果
func Send(title, msg string)
```

#### 使用示例

```go
import "watchducker/pkg/notify"

// 在容器更新成功后调用
notify.Send(
    "容器更新完成",
    "nginx 容器从 1.20.0 已更新至 1.21.0",
)
```

#### 调用流程

```
1. 检查 push.yaml 是否存在
   ├─ 不存在 → 记录日志并返回
   └─ 存在   → 继续

2. 加载配置文件
   ├─ 失败 → 记录错误并返回
   └─ 成功 → 继续

3. 解析 PushServer 配置
   ├─ 为空 → 记录日志并返回
   └─ 有值 → 继续

4. 遍历各推送平台
   ├─ 调用对应函数（async-like）
   ├─ 捕获异常与错误
   ├─ 记录成功/失败日志
   └─ 继续处理其他平台
```

## 内部工具函数

### HTTP 工具

```go
// postJSON 发送 JSON POST 请求
// 返回响应体，错误时返回 nil
func postJSON(url string, body interface{}) ([]byte, error)

// postForm 发送表单 POST 请求
// 返回响应体，错误时返回 nil
func postForm(url string, data url.Values) ([]byte, error)
```

### 加密工具

```go
// 钉钉签名（HMAC-SHA256）
h := hmac.New(sha256.New, []byte(secret))
h.Write([]byte(stringToSign))
sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
```

## 与其他模块的交互

```
┌─────────────────────────────────────────────────────────┐
│                    internal/core                        │
│         (容器更新后调用 notify.Send)                    │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
        ┌──────────────────────────────┐
        │   pkg/notify.Send()          │
        └──────────────────────────────┘
                       │
         ┌─────────────┼─────────────┐
         ▼             ▼             ▼
    [Telegram]   [Dingrobot]  [Gotify]
                      │
                      ▼
            (外部通知平台 API)
```

### 调用示例（来自 internal/core）

```go
// 更新容器后
if success {
    notify.Send(
        "容器更新成功",
        fmt.Sprintf("%s: %s → %s", containerName, oldTag, newTag),
    )
}
```

## 错误处理策略

### 设计原则

1. **错误隔离**：单个推送失败不阻塞其他推送
2. **优雅降级**：配置缺失时自动跳过
3. **详细日志**：所有错误都通过 logger 记录

### 错误场景

| 场景 | 处理方式 | 日志输出 |
|------|--------|--------|
| 配置文件不存在 | 跳过推送 | `Info: 未找到推送配置文件，跳过推送` |
| 配置解析失败 | 返回错误 | `Error: 配置解析失败` |
| 网络连接失败 | 记录错误 | `Error: XXX 失败: connection refused` |
| 推送成功 | 继续下一个 | `Info: XXX 成功` |

## 安全考虑

### 敏感信息保护

- ⚠️ **Token 和密钥**应存储在 `push.yaml` 中（**不要提交到 Git**）
- 建议在 `.gitignore` 中添加：
  ```
  push.yaml
  *.key
  *.secret
  ```

### 网络通信

- HTTP 调用使用标准的 `net/http` 库
- Telegram 支持代理转接（`api_url` 字段支持自定义）
- SMTP 使用 PlainAuth（对于 TLS 端口应配置 `587` 或 `465`）

### API 调用限制

- 建议在 `Send()` 调用前进行速率限制
- 某些平台（如 Telegram、钉钉）有 API 频率限制
- 考虑实现重试机制（当前版本无重试）

## 测试状态与建议

### 当前状态

- ❌ 无单元测试
- ❌ 无集成测试
- ⚠️ 仅在 `core/operator.go` 中通过 `notify.Send()` 调用测试

### 建议的测试用例

#### 单元测试

```go
// 配置加载测试
TestLoadConfig_Success()      // 有效配置
TestLoadConfig_NotFound()     // 配置不存在
TestLoadConfig_InvalidYAML()  // 无效 YAML 格式

// HTTP 工具测试
TestPostJSON_Success()        // 成功发送 JSON
TestPostJSON_NetworkError()   // 网络错误

// 各推送平台测试
TestTelegram_Success()
TestTelegram_InvalidToken()
TestDingrobot_WithSecret()
// ... 等等
```

#### 集成测试

```go
// 完整流程测试
TestSend_MultipleServers()    // 多渠道推送
TestSend_PartialFailure()     // 部分推送失败
TestSend_NoConfigFile()       // 无配置文件
```

#### 模拟测试

```go
// 使用 httptest 模拟 HTTP 服务器
server := httptest.NewServer(...)
defer server.Close()
```

## 性能与可靠性

### 当前特性

- **顺序处理**：逐个调用各推送平台（非并发）
- **无重试**：单次失败立即记录错误
- **无超时控制**：HTTP 调用使用默认超时
- **配置缓存**：每次 `Send()` 都重新加载配置

### 优化建议

1. **并发推送**：使用 goroutine 池并发调用各平台
   ```go
   go telegram(title, msg)
   go dingrobot(title, msg)
   // ...
   ```

2. **添加超时**：为 HTTP Client 设置超时
   ```go
   client := &http.Client{Timeout: 10 * time.Second}
   ```

3. **实现重试机制**：失败自动重试（指数退避）

4. **配置持久化**：缓存配置以避免重复读取磁盘

## 扩展新的通知平台

### 步骤

1. **在 `Config` 中添加结构体**
   ```go
   NewPlatform struct {
       Field1 string `mapstructure:"field1"`
       Field2 string `mapstructure:"field2"`
   } `mapstructure:"newplatform"`
   ```

2. **实现推送函数**
   ```go
   func newplatform(title, msg string) {
       // 验证配置
       // 构建请求
       // 发送请求
       // 处理错误和日志
   }
   ```

3. **在 `Send()` 中添加 case 分支**
   ```go
   case "newplatform":
       newplatform(title, msg)
   ```

4. **更新 `push.yaml.example`**

5. **编写文档和测试**

## 依赖关系

### 外部依赖

```go
import (
    "net/http"        // HTTP 请求
    "net/smtp"        // SMTP 邮件
    "crypto/hmac"     // HMAC 签名
    "encoding/json"   // JSON 序列化
    "github.com/spf13/viper"  // 配置管理
)
```

### 内部依赖

```go
import (
    "watchducker/pkg/logger"  // 日志服务
)
```

## 常见问题

### Q1: 如何不使用通知功能？
**A**: 删除或重命名 `push.yaml` 文件，模块会自动跳过推送。

### Q2: 如何同时使用多个推送平台？
**A**: 在 `push.yaml` 中配置多个平台，并用逗号分隔：
```yaml
setting:
  push_server: "telegram,dingrobot,ftqq,bark"
```

### Q3: 推送失败是否会影响容器更新？
**A**: 否。推送是独立的，失败只记录日志不影响主流程。

### Q4: 如何调试推送问题？
**A**: 在 `push.yaml` 中设置：
```yaml
setting:
  log_level: "DEBUG"
```
然后检查日志输出。

### Q5: Telegram API 被限制怎么办？
**A**: 修改 `api_url` 使用代理或反代服务器：
```yaml
telegram:
  api_url: "your-proxy.com"  # 代替 api.telegram.org
```

## 相关文档链接

- [根模块文档](../CLAUDE.md)
- [core 模块（集成 notify）](../internal/core/CLAUDE.md)
- [logger 模块（日志依赖）](../logger/CLAUDE.md)
- [push.yaml 配置示例](../../push.yaml.example)

## 变更记录

| 日期 | 版本 | 变更说明 |
|------|------|--------|
| 2025-11-13 | 1.0 | 初始文档创建，覆盖 15 种推送平台 |
| 2025-11-13 | 1.0 | 记录测试缺失和优化建议 |

---

**文档统计**: 15 种推送平台、520+ 行代码、完整 API 说明
**关键指标**: 无单元测试、支持多渠道推送、错误隔离设计
**下一步**: 建议编写单元测试和并发优化实现
