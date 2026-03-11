# AI Bridge 插件说明

**版本**: 1.1.0
**作者**: [西风](https://xifeng.net)
**插件主页**: [https://xifeng.net/global-ai-bridge-wp-plugins](https://xifeng.net/global-ai-bridge-wp-plugins)

## 插件简介

AI Bridge 是一个 WordPress 插件，用于把 WordPress 侧的 AI 请求统一转发到你自己的代理网关，再由代理网关连接 OpenAI、Gemini、Claude、DeepSeek 等模型服务。

这个插件适合以下场景：

- 服务器位于中国大陆或香港，直连海外 AI API 不稳定
- 海外服务器需要回国访问国内 AI 服务
- 希望把 AI 服务的 API Key 集中保存在 VPS 或网关服务上
- 多个 WordPress 站点共用同一个 AI 代理层
- 需要在代理层做统一路由、日志、缓存或后续计费控制

## 当前已实现能力

### 核心代理

- WordPress 插件主入口
- 后台设置页面（工具 -> AI Bridge）
- 云节点 / 自定义地址双连接模式
- 内置美国节点和中国节点
- 出国 / 回国模式切换
- 代理地址配置
- 站点访问 token 配置和自动生成
- Provider API Token 配置
- 默认 Provider 配置
- 默认模型配置
- 请求超时配置
- 日志开关
- 缓存偏好标记
- PHP 调用入口 `gab_send_ai_request()`

### HTTP 请求拦截

- 自动拦截 WordPress 发往已知 AI 提供商的 HTTP 请求
- 通过主机名和路径模式匹配拦截目标
- 被拦截请求自动改走 AI Bridge 代理
- 直接代理请求改写（支持主题或插件直接使用 AI Bridge 地址）
- AI Bridge 响应自动转换为 OpenAI 兼容格式

### 兼容性

- WordPress Connectors 兼容层，保留代理环境下的 API Key
- Lared 主题兼容，自动注入 Provider API Token

### 监控与统计

- 节点测速按钮
- 使用统计面板（请求数、Token 用量、平均延迟、成功率）
- 7 天请求趋势和 Token 趋势图表
- 最近请求日志，支持按提供商筛选
- 连接器接管日志

## 工作方式

AI Bridge 的基础调用链路如下：

1. WordPress 内部功能或其他插件触发 AI 请求
2. AI Bridge 组装统一请求体
3. 请求发送到你配置的代理网关地址
4. 代理网关根据 provider 和 model 决定路由
5. 上游模型返回结果
6. 代理网关把结果返回给 WordPress
7. 插件记录最近请求日志

## 当前配置项

插件后台路径：

`WordPress 后台 -> 设置 -> AI Bridge`

当前支持配置项：

- `Traffic Mode`
  可选 `Outbound / 出国` 和 `Inbound / 回国`
- `Proxy Endpoint`
  你的代理接口地址，例如 `https://proxy.example.com/v1/chat/completions`
- `Site Token`
  插件发送给代理的桥接访问 Token
- `Provider API Token`
  用户自己的 OpenAI、Claude、Gemini 或其他模型服务 token
- `Default Provider`
  默认提供商，目前可选 `openai`、`claude`、`google`、`gemini`、`qwen`、`baidu`、`deepseek`、`doubao`、`kimi`、`minimax`
- `Default Model`
  默认模型名称，例如 `gpt-4.1-mini`
- `Request Timeout`
  请求超时时间，单位秒
- `Enable Logging`
  是否记录最近请求日志
- `Enable Cache Flag`
  是否向代理发送缓存偏好标记

## 当前请求结构

插件向代理发送的是 JSON POST 请求，主体结构大致如下：

```json
{
  "provider": "openai",
  "model": "gpt-4.1-mini",
  "messages": [
    {
      "role": "user",
      "content": "Hello"
    }
  ],
  "stream": false,
  "meta": {
    "source": "wordpress-ai-bridge",
    "site": "https://example.com"
  },
  "temperature": 0.7,
  "max_tokens": 512,
  "cache": true
}
```

请求头大致如下：

```http
Content-Type: application/json
Accept: application/json
Authorization: Bearer your-site-token
X-AIBRIDGE-PROVIDER-TOKEN: user-provider-token
```

## 回国 / 出国模式

- `Outbound / 出国`
  适用于中国大陆或香港站点访问海外 AI 服务，默认可选 `US`、`SG`、`JP`
- `Inbound / 回国`
  适用于海外服务器回国访问国内 AI 服务，默认可选 `CN`

当前内置中国节点：

```text
https://cn-aibridge.yite.net/v1/chat/completions
```

## 日志能力

当前版本会在 WordPress option 中保存最近 100 条请求日志，日志字段包括：

- 请求时间
- Provider
- Model
- 状态码或错误状态
- 延迟时间
- Token 数量

这部分日志主要用于开发和联调，不适合作为长期审计系统。

## Provider 与模型

插件的 `provider` 和 `model` 是分开的：

- `provider` 决定后端路由到哪个厂商
- `model` 由用户或业务代码自由填写

后端不应该把模型名写死。只要对应 provider 支持该模型，网关就可以原样透传模型名称。

## 适用边界

当前版本已经具备完整的代理转发、HTTP 拦截和兼容性能力，可用于生产环境。以下能力尚未实现：

- WordPress REST API 对外接口
- 数据库存储级别的详细调用记录
- 更细粒度的权限和站点签名机制
- 流式返回处理
- 完整统计报表和导出

## 目录说明

- `global-ai-bridge.php`
  插件主文件
- `includes/class-gab-plugin.php`
  插件启动和核心实例管理
- `includes/class-gab-settings.php`
  后台设置页、配置存储和统计面板
- `includes/class-gab-proxy-client.php`
  代理请求客户端
- `includes/class-gab-http-interceptor.php`
  HTTP 请求拦截器，自动拦截发往 AI 提供商的请求
- `includes/class-gab-connectors-compat.php`
  WordPress Connectors 兼容层
- `includes/class-gab-theme-compat.php`
  主题兼容层（当前支持 Lared 主题）
- `includes/class-gab-logger.php`
  日志记录器

## 下一步开发建议

建议按以下顺序继续完善：

1. 添加“测试连接”按钮
2. 增加 WordPress REST API 接口
3. 增加代理响应标准化处理
4. 增加与具体 AI 插件或编辑器功能的适配层
5. 增加数据库日志表而不是仅使用 option 存储
