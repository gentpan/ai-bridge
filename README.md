<div align="center">

# AI Bridge

**让 WordPress 连接全球 AI —— 你的专属 AI API 反向代理网关**

<p>
  <img src="https://img.shields.io/badge/WordPress-5.0%2B-21759b?style=for-the-badge&logo=wordpress&logoColor=white" alt="WordPress">
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/PHP-7.4%2B-777BB4?style=for-the-badge&logo=php&logoColor=white" alt="PHP">
  <img src="https://img.shields.io/github/v/release/gentpan/ai-bridge?style=for-the-badge" alt="Release">
  <img src="https://img.shields.io/github/license/gentpan/ai-bridge?style=for-the-badge" alt="License">
</p>

<p>
  <a href="https://github.com/gentpan/ai-bridge-php">PHP 单文件版</a> ·
  <a href="https://github.com/gentpan/ai-bridge-go">Go 独立版</a>
</p>

</div>

---

## 📖 概述

AI Bridge 解决一个核心痛点：**中国大陆及香港的服务器无法直接访问海外 AI API**（OpenAI、Claude、Gemini 等）。

通过在海外的 VPS 上部署一台轻量级反向代理网关，国内/港区的 WordPress 或其他应用可以将 AI 请求发送到你的海外网关，由网关转发至 AI 服务商。

**安全保障：** API Key 仅在你的服务器上流转，不经任何第三方平台。

---

## 📦 项目结构

本仓库分为两个独立组件：

| 组件 | 路径 | 说明 |
|------|------|------|
| WordPress 插件 | `global-ai-bridge.php` + `includes/` | 安装在 WordPress 站点上，拦截 AI 请求并转发到海外网关 |
| Go 后端代理 | `server/` | 部署在海外 VPS 上，接收转发请求并代理到 AI 服务商 |

---

## 🚀 快速开始

### WordPress 插件

1. 下载 `global-ai-bridge.php` 和 `includes/` 目录
2. 上传到 WordPress 的 `/wp-content/plugins/global-ai-bridge/`
3. 在 WordPress 后台启用插件
4. 在设置页面填入你的网关地址和 API Key

### Go 后端代理

请移步 [`server/`](server/) 目录查看部署文档。

---

## 🔧 工作原理

```
┌─────────────────────┐         ┌──────────────────┐         ┌────────────────┐
│  WordPress 站点      │  POST   │  AI Bridge 网关  │  POST   │  AI 服务商      │
│  (国内/香港服务器)    │ ──────→ │  (海外 VPS)      │ ──────→ │  (OpenAI 等)   │
│                     │         │                  │         │                │
│  AI Bridge 插件      │         │  1. 读取 provider │         │                │
│  拦截请求并发送       │         │  2. 查转发地址表   │         │                │
│  到网关地址          │         │  3. 转发请求+响应  │         │                │
└─────────────────────┘         └──────────────────┘         └────────────────┘
```

---

## 🔗 关联项目

| 仓库 | 说明 |
|------|------|
| [ai-bridge-php](https://github.com/gentpan/ai-bridge-php) | PHP 单文件版后端正向代理（适合共享主机，无 Docker） |
| [ai-bridge-go](https://github.com/gentpan/ai-bridge-go) | Go 版后端代理（高性能推荐，Docker 部署） |

---

## 📄 License

MIT
