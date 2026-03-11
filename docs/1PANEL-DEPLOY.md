# 1Panel 部署指南

这份文档按 1Panel 当前常见流程整理，目标是把 AI Bridge Go 后端部署到 1Panel 管理的服务器中，并通过域名反向代理暴露给 WordPress 插件使用。

## 你将用到的文件

- [server/docker-compose.yml](/Users/gentpan/Local%20Sites/pan/app/public/wp-content/plugins/global-ai-bridge/server/docker-compose.yml)
- [server/.env.production.example](/Users/gentpan/Local%20Sites/pan/app/public/wp-content/plugins/global-ai-bridge/server/.env.production.example)
- [server/nginx/aibridge.conf](/Users/gentpan/Local%20Sites/pan/app/public/wp-content/plugins/global-ai-bridge/server/nginx/aibridge.conf)

## 一、准备服务器

建议起步配置：

- `2 vCPU`
- `4 GB RAM`
- `40 GB SSD`

开放端口建议：

- `80`
- `443`
- `8090` 或你的 1Panel 面板端口

如果你还没有安装 1Panel，官方文档当前提供在线安装脚本，并说明安装完成后可通过 `http://服务器IP:目标端口/安全入口` 登录，且可用 `1pctl user-info` 获取入口信息。
来源：
- [1Panel 在线安装文档](https://1panel.cn/docs/v2/installation/online_installation/)

## 二、安装 1Panel

SSH 登录服务器后执行：

```bash
bash -c "$(curl -sSL https://resource.fit2cloud.com/1panel/package/v2/quick_start.sh)"
```

安装完成后：

1. 记录终端输出的面板地址、账号和密码
2. 如果忘记入口地址，执行：

```bash
1pctl user-info
```

3. 浏览器打开面板地址并登录

## 三、把项目文件传到服务器

推荐把仓库放到例如这个目录：

```text
/opt/aibridge
```

你可以用任意一种方式上传：

- `git clone`
- SFTP
- 1Panel 文件管理上传

最终你需要保证服务器上存在：

```text
/opt/aibridge/server/docker-compose.yml
/opt/aibridge/server/.env
```

## 四、生成生产环境变量文件

进入后端目录：

```bash
cd /opt/aibridge/server
cp .env.production.example .env
```

然后编辑 `.env`，至少填写这些值：

- `SITE_TOKEN`
- `METRICS_TOKEN`
- `DEFAULT_PROVIDER`
- 对应 provider 的 `BASE_URL`
- 对应 provider 的 `DEFAULT_MODEL`

一个最小 OpenAI 示例：

```env
SITE_TOKEN=replace-with-a-very-long-random-site-token
METRICS_TOKEN=replace-with-another-random-token
DEFAULT_PROVIDER=openai
DEFAULT_MODEL=gpt-4.1-mini
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_DEFAULT_MODEL=gpt-4.1-mini
```

如果你还想启用 Claude：

```env
ANTHROPIC_BASE_URL=https://api.anthropic.com/v1
ANTHROPIC_DEFAULT_MODEL=claude-3-7-sonnet-latest
```

如果你还想启用 Google/Gemini：

```env
GOOGLE_BASE_URL=https://generativelanguage.googleapis.com/v1beta
GOOGLE_DEFAULT_MODEL=gemini-2.0-flash
```

## 五、在 1Panel 里创建 Compose 编排

1Panel 官方文档当前说明，创建 Compose 支持三种方式：

- 编辑
- 路径选择
- 编排模板

来源：
- [1Panel 编排文档](https://1panel.cn/docs/v1/user_manual/containers/compose/)

推荐你用“路径选择”。

具体步骤：

1. 登录 1Panel
2. 进入 `容器`
3. 进入 `编排`
4. 点击 `创建编排`
5. 选择 `路径选择`
6. 选择服务器上的：

```text
/opt/aibridge/server/docker-compose.yml
```

7. 确认创建
8. 启动编排

如果你更习惯在 1Panel 里直接粘贴内容，也可以选“编辑”，把 [server/docker-compose.yml](/Users/gentpan/Local%20Sites/pan/app/public/wp-content/plugins/global-ai-bridge/server/docker-compose.yml) 内容粘进去。

## 六、检查容器是否正常运行

在 1Panel 里：

1. 进入 `容器`
2. 查看容器列表
3. 找到 `aibridge`
4. 查看日志是否有类似：

```text
ai bridge gateway listening on :8080
```

也可以在服务器上手动检查：

```bash
cd /opt/aibridge/server
docker compose ps
docker compose logs -f
```

## 七、在 1Panel 里创建反向代理网站

1Panel 当前支持“反向代理”类型网站，并可把域名请求转发到已有服务地址。
来源：
- [1Panel 创建网站文档](https://1panel.cn/docs/v2/user_manual/websites/website_create/)

具体步骤：

1. 进入 `网站`
2. 点击 `创建网站`
3. 选择 `反向代理`
4. 主域名填写你的网关域名，例如：

```text
us-aibridge.bluecdn.com
```

5. 代理地址填写：

```text
http://127.0.0.1:8080
```

6. 开启 HTTPS
7. 选择或申请 SSL 证书
8. 保存

这样外部访问：

```text
https://us-aibridge.bluecdn.com/v1/chat/completions
```

就会转发到你的 Go 服务。

## 八、验证接口

先验证健康检查：

```bash
curl http://127.0.0.1:8080/healthz
```

再验证外部域名：

```bash
curl https://us-aibridge.bluecdn.com/healthz
```

然后测试 AI 接口：

```bash
curl -X POST https://us-aibridge.bluecdn.com/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer 你的_SITE_TOKEN' \
  -H 'X-AIBRIDGE-PROVIDER-TOKEN: 你的_provider_token' \
  -d '{
    "provider": "openai",
    "model": "gpt-4.1-mini",
    "messages": [
      { "role": "user", "content": "hello" }
    ]
  }'
```

## 九、在 WordPress 插件里填写

如果你用的是你的云节点，在插件后台填写：

- `Connection Mode`: `Use AI Bridge Cloud` 或 `Use Custom Endpoint`
- 如果是自定义方式：
  - `Custom Endpoint`: `https://us-aibridge.bluecdn.com/v1/chat/completions`
  - `Site Token`: 你在 `.env` 里配置的桥接访问 token
  - `Provider API Token`: 用户自己的 OpenAI / Claude / Gemini token

## 十、安全建议

至少做到这些：

- `SITE_TOKEN` 和 `METRICS_TOKEN` 使用高强度随机字符串
- 不要在服务端日志里记录 `X-AIBRIDGE-PROVIDER-TOKEN`
- 用户自己的 provider token 只保存在插件设置中
- 开启 HTTPS
- 在 1Panel 或云防火墙只开放 `80/443`
- 不要把 `8080` 直接暴露公网
- 给 `/metrics` 单独设置 token
- 定期轮换 `SITE_TOKEN`

## 十一、推荐运行方式

第一版推荐只上一个节点：

- `us-aibridge.bluecdn.com`

跑稳以后再复制同样流程部署：

- `sg-aibridge.bluecdn.com`
- `jp-aibridge.bluecdn.com`

## 十二、1Panel 官方文档参考

- [在线安装](https://1panel.cn/docs/v2/installation/online_installation/)
- [编排](https://1panel.cn/docs/v1/user_manual/containers/compose/)
- [创建网站](https://1panel.cn/docs/v2/user_manual/websites/website_create/)
- [应用安装](https://1panel.cn/docs/v2/user_manual/appstore/install/)
