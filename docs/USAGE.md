# AI Bridge 使用说明

## 1. 安装

将插件目录放入 WordPress 插件目录：

```text
wp-content/plugins/global-ai-bridge
```

然后在 WordPress 后台启用 `AI Bridge` 插件。

## 2. 基础配置

启用后进入：

`设置 -> AI Bridge`

填写以下内容：

### Traffic Mode

可选：

- `Outbound / 出国`
- `Inbound / 回国`

出国模式通常用于中国大陆或香港站点访问海外 AI。回国模式通常用于海外服务器访问国内 AI。

### Proxy Endpoint

填写你的代理接口地址，例如：

```text
https://proxy.example.com/v1/chat/completions
```

如果使用内置中国节点，可填：

```text
https://cn-aibridge.yite.net/v1/chat/completions
```

### Site Token

填写你的代理网关分配给当前 WordPress 站点的访问 token。

### Provider API Token

填写用户自己的模型服务 token，例如：

- OpenAI API key
- Claude API key
- Gemini API key
- 其他兼容 provider 的 token

### Default Provider

可选值：

- `openai`
- `claude`
- `google`
- `gemini`
- `qwen`
- `baidu`
- `deepseek`
- `doubao`
- `kimi`
- `minimax`

### Default Model

填写默认模型名称，例如：

```text
gpt-4.1-mini
```

模型名称支持自定义填写，后端会把它透传给对应 provider。

### Request Timeout

请求超时秒数，建议从 `30` 或 `60` 开始。

### Enable Logging

开启后，插件会在设置页下方显示最近请求日志。

### Enable Cache Flag

开启后，插件会向代理请求体发送 `cache: true` 标记，具体是否缓存由你的代理服务决定。

## 3. 代理接口预期

当前插件通过 `POST` 方式把 JSON 发到你配置的代理地址。

代理至少需要支持：

- 接收 JSON 请求体
- 读取 `Authorization: Bearer <token>`
- 读取 `X-AIBRIDGE-PROVIDER-TOKEN`
- 按 `provider` 和 `model` 字段路由到上游模型
- 返回 JSON 响应

建议代理返回结构兼容常见 LLM API 响应，并尽量包含 `usage.total_tokens`，这样插件可以记录 token 数。

## 4. 在主题或插件中调用

当前版本提供了一个简单调用函数：

```php
$response = gab_send_ai_request(
    array(
        array(
            'role'    => 'system',
            'content' => 'You are a helpful assistant.',
        ),
        array(
            'role'    => 'user',
            'content' => 'Write a short excerpt for my blog post.',
        ),
    ),
    array(
        'provider'    => 'openai',
        'model'       => 'gpt-4.1-mini',
        'temperature' => 0.7,
        'max_tokens'  => 300,
    )
);

if ( is_wp_error( $response ) ) {
    error_log( $response->get_error_message() );
} else {
    // 根据你的代理返回结构读取结果。
    error_log( wp_json_encode( $response ) );
}
```

实际转发时，插件会自动附带：

- 桥接访问 `site token`
- 用户自己的 `provider api token`

## 5. 请求参数说明

`gab_send_ai_request( $messages, $args )` 当前支持的主要参数：

### `$messages`

标准聊天消息数组，例如：

```php
array(
    array(
        'role'    => 'user',
        'content' => 'Hello',
    ),
)
```

### `$args`

可选参数包括：

- `provider`
- `model`
- `temperature`
- `max_tokens`
- `stream`

如果未传 `provider` 或 `model`，会自动使用后台默认配置。

## 6. 日志查看

打开：

`设置 -> AI Bridge`

页面下方可以看到最近请求日志，包括：

- 时间
- Provider
- Model
- 状态
- 延迟
- Tokens

## 7. 常见问题

### 配置后请求失败

优先检查：

- `Proxy Endpoint` 是否可访问
- `Site Token` 是否正确
- `Provider API Token` 是否正确
- 代理服务是否支持当前请求结构
- WordPress 服务器是否可以访问你的代理地址

### 后台有日志但没有内容结果

说明插件已经发出请求，但你的代理返回格式可能与当前业务代码预期不一致。当前版本会原样返回 JSON 数组，不会强制转换成统一文本结果。

### Token 数始终为 0

说明代理响应中没有包含插件当前识别的 token 字段。当前版本优先读取：

- `usage.total_tokens`
- `usage.totalTokenCount`

## 8. 当前版本限制

当前版本还不包含：

- 流式输出处理
- 前端可视化测试面板
- 统一 REST API 调用入口
- 多站点权限管理
- 完整统计报表

这些能力可以在下一阶段继续开发。
