# Changelog

All notable changes to this project will be documented in this file.

## 1.1.0 - 2026-03-11

### Added

- HTTP 请求拦截器（`GAB_HTTP_Interceptor`），自动拦截 WordPress 中发往已知 AI 提供商的 HTTP 请求并通过 AI Bridge 代理转发
- WordPress Connectors 兼容层（`GAB_Connectors_Compat`），解决通过代理时 API Key 验证失败的问题
- Lared 主题兼容层（`GAB_Theme_Compat`），自动为 Lared 主题注入 Provider API Token
- 云节点 / 自定义地址双连接模式，支持一键切换
- 内置美国节点（us-aibridge.bluecdn.com）和中国节点（cn-aibridge.yite.net）
- 出国 / 回国流量方向切换
- 云节点区域选择（美国 / 中国大陆）
- 节点测速按钮，支持在后台直接测试当前节点延迟
- 使用统计面板，展示请求数、Token 用量、平均延迟、成功请求数
- 7 天请求趋势和 Token 趋势图表
- 连接器接管日志，记录被拦截的第三方 AI 请求
- 按提供商筛选日志
- Site Token 自动生成和重新生成按钮
- Provider API Token 输入和管理
- 提供商选择时显示对应 API Token 获取链接
- 直接代理请求处理，支持主题或插件直接以 AI Bridge 地址作为 endpoint 时的自动改写
- AI Bridge 响应到 OpenAI 兼容格式的标准化转换

### Changed

- 设置页面从 Settings 菜单移至 Tools 菜单
- 设置页面重新设计，新增品牌头部、统计面板和可视化图表
- 插件作者更新为 西风（xifeng.net）
- 插件链接更新为 https://xifeng.net/global-ai-bridge-wp-plugins

## 0.1.0 - 2026-03-06

### Added

- Initial WordPress plugin bootstrap
- Admin settings page under `Settings -> AI Bridge`
- Proxy endpoint, site token, default provider, default model, timeout, logging, and cache settings
- Proxy request client based on `wp_remote_post()`
- Simple request logging stored in WordPress options
- Public helper function `gab_send_ai_request()`
- Project documentation files: `README.md`, `PLUGIN.md`, `USAGE.md`, `CHANGELOG.md`

### Notes

- This is the first working scaffold version of the plugin
- Deep integration with specific WordPress AI plugins is not implemented yet
