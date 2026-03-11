# AI Bridge for WordPress

AI Bridge 是一个 WordPress 插件，用于将 WordPress 的 AI 请求统一转发到你自己的代理网关，再由代理网关连接 OpenAI、Gemini、Claude、DeepSeek 等全球 AI 模型服务。

适用于服务器位于中国大陆或香港等直连海外 AI API 不稳定的环境，也支持海外服务器回国访问国内 AI 服务。

- **作者**: 西风
- **作者主页**: [https://xifeng.net](https://xifeng.net)
- **插件主页**: [https://xifeng.net/global-ai-bridge-wp-plugins](https://xifeng.net/global-ai-bridge-wp-plugins)
- **GitHub**: [https://github.com/gentpan/ai-bridge](https://github.com/gentpan/ai-bridge)
- **当前版本**: 1.1.0

## 文档

- [插件说明](./PLUGIN.md)
- [使用指南](./USAGE.md)
- [更新日志](./CHANGELOG.md)
- [后端说明](../server/README.md)
- [部署指南](./DEPLOY.md)
- [1Panel 部署指南](./1PANEL-DEPLOY.md)

## 功能特性

### 代理与路由

- 通过代理网关连接全球 AI 模型（OpenAI / Gemini / Claude / DeepSeek 等）
- 内置云节点（美国 / 中国大陆），开箱即用
- 支持自定义代理地址，兼容自建网关或第三方中转服务
- 出国模式 / 回国模式一键切换
- 多模型路由，模型名称透传

### HTTP 请求拦截

- 自动拦截 WordPress 中发往已知 AI 提供商的 HTTP 请求
- 被拦截的请求自动改走 AI Bridge 代理
- 支持已知主机匹配和路径模式匹配
- 直接代理请求改写，支持主题或插件直接使用 AI Bridge 地址
- AI Bridge 响应自动转换为 OpenAI 兼容格式

### 兼容性

- WordPress Connectors 兼容层，解决代理环境下 API Key 验证失败问题
- Lared 主题兼容，自动注入 Provider API Token
- 兼容 WordPress AI 集成功能

### 管理与监控

- 后台设置页面（工具 -> AI Bridge）
- 节点测速按钮
- 使用统计面板（请求数、Token 用量、平均延迟、成功率）
- 7 天请求趋势和 Token 趋势图表
- 最近请求日志，支持按提供商筛选
- 连接器接管日志
- Site Token 自动生成和管理
- Provider API Token 安全存储

### 安全

- WordPress 仅存储代理端点和站点级认证数据
- Provider API Key 保存在本地，通过代理转发
- 代理层可集中管理访问控制、速率限制和日志

## 架构

AI Bridge 使用简单的中继架构：

1. WordPress 发起 AI 请求
2. AI Bridge 插件拦截或组装请求
3. 请求通过代理网关转发到上游 AI 提供商
4. 上游返回结果经代理标准化后回传给 WordPress

### 流程图

```text
+-------------------+       +---------------------+       +----------------------+
| WordPress Site    | ----> | AI Bridge Plugin    | ----> | Private Proxy Gateway|
| (admin/frontend)  |       | (intercept/proxy)   |       | (routing/auth/cache) |
+-------------------+       +---------------------+       +----------------------+
                                                                    |
                                                                    v
                                                        +--------------------------+
                                                        | AI Providers             |
                                                        | OpenAI / Gemini / Claude |
                                                        | DeepSeek / Others        |
                                                        +--------------------------+
```

## 核心组件

### WordPress 插件层

- `global-ai-bridge.php` — 插件主入口
- `includes/class-gab-plugin.php` — 核心实例管理
- `includes/class-gab-settings.php` — 后台设置与统计面板
- `includes/class-gab-proxy-client.php` — 代理请求客户端
- `includes/class-gab-http-interceptor.php` — HTTP 请求拦截器
- `includes/class-gab-connectors-compat.php` — WordPress Connectors 兼容层
- `includes/class-gab-theme-compat.php` — 主题兼容层
- `includes/class-gab-logger.php` — 日志记录器

### 代理网关层（Go 后端）

- 接收 WordPress 请求并验证站点令牌
- 按 provider 和 model 路由到上游 AI 提供商
- 标准化不同提供商的请求和响应格式
- 支持重试、超时、日志和可选缓存

## 适用场景

- WordPress 站点部署在中国大陆或香港，直连海外 AI API 不稳定
- 海外服务器需要回国访问国内 AI 服务
- 多个 WordPress 站点共用同一个 AI 代理层
- 需要在代理层做统一路由、日志、缓存或计费控制
- 希望把 AI 服务的 API Key 集中管理

## 许可

本插件由 [西风](https://xifeng.net) 开发和维护。
