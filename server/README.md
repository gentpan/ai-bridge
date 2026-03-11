# AI Bridge Gateway

This is the Go backend for AI Bridge. It receives requests from the WordPress plugin, validates a bridge site token, forwards the request to an upstream AI provider using the user's own provider token, and returns a normalized JSON response.

## Provider Coverage

Native adapters implemented now:

- OpenAI-compatible chat completions
- Claude Messages API
- Google / Gemini `generateContent`

Compatible adapter slots available now:

- Qwen
- Baidu
- DeepSeek
- Doubao
- Kimi
- MiniMax

The compatible adapter expects either:

- an OpenAI-compatible `POST /chat/completions` interface, or
- your own relay that normalizes the provider into that shape

For providers with native protocols that differ from this shape, you should either place a relay in front of them or add another native adapter in `internal/providers/`.

## Endpoints

- `GET /healthz`
- `POST /v1/chat/completions`

## Quick Start

1. Copy `.env.example` into your environment
2. Set `SITE_TOKEN`
3. Configure the provider base URLs you want to expose
4. Start the server:

```bash
go run ./cmd/server
```

## Example Request

```bash
curl -X POST http://127.0.0.1:8080/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer replace-with-a-random-site-token' \
  -H 'X-AIBRIDGE-PROVIDER-TOKEN: sk-user-provider-token' \
  -d '{
    "provider": "openai",
    "model": "gpt-4.1-mini",
    "messages": [
      {"role":"user","content":"Hello"}
    ]
  }'
```

## Response Shape

```json
{
  "id": "chatcmpl_xxx",
  "provider": "openai",
  "model": "gpt-4.1-mini",
  "content": "Hello",
  "raw": {
    "latency_ms": 842
  },
  "usage": {
    "prompt_tokens": 9,
    "completion_tokens": 12,
    "total_tokens": 21
  }
}
```

## BYOK flow

The gateway is designed for bring-your-own-key mode:

- the plugin stores the user's own provider token
- the plugin sends that token in `X-AIBRIDGE-PROVIDER-TOKEN`
- the gateway forwards the request with that token
- the gateway does not need user provider tokens stored in `.env`

## Custom Models

Models are not hardcoded in the server. The request `model` is passed through to the selected provider. If the request does not include a model, the server falls back to that provider's configured default model.
