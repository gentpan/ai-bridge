package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	nethttp "net/http"
	"net/url"
	"strings"
	"time"

	"global-ai-bridge/server/internal/config"
	"global-ai-bridge/server/internal/gateway"
	"global-ai-bridge/server/internal/service"
)

type Server struct {
	config      config.Config
	chatService *service.ChatService
	metrics     *Metrics
	rateLimiter *rateLimiter
}

func NewServer(cfg config.Config) nethttp.Handler {
	server := &Server{
		config:      cfg,
		chatService: service.NewChatService(cfg),
		metrics:     NewMetrics(),
		rateLimiter: newRateLimiter(cfg.RateLimitPerMinute),
	}

	mux := nethttp.NewServeMux()
	mux.HandleFunc("/healthz", server.handleHealth)
	mux.HandleFunc("/metrics", server.handleMetrics)
	mux.HandleFunc("/v1/connectors/proxy", server.handleConnectorProxy)
	mux.HandleFunc("/v1/chat/completions", server.handleChat)

	return server.withAccessLog(mux)
}

func (s *Server) handleHealth(writer nethttp.ResponseWriter, request *nethttp.Request) {
	writeJSON(writer, nethttp.StatusOK, map[string]any{
		"ok":           true,
		"time":         time.Now().UTC().Format(time.RFC3339),
		"node_name":    s.config.NodeName,
		"traffic_mode": s.config.NodeTrafficMode,
	})
}

func (s *Server) handleMetrics(writer nethttp.ResponseWriter, request *nethttp.Request) {
	if s.config.MetricsToken != "" {
		header := strings.TrimSpace(request.Header.Get("Authorization"))
		expected := "Bearer " + s.config.MetricsToken
		if header != expected {
			writeError(writer, nethttp.StatusUnauthorized, "unauthorized", "invalid metrics token")
			return
		}
	}

	writeJSON(writer, nethttp.StatusOK, s.metrics.Snapshot())
}

func (s *Server) handleChat(writer nethttp.ResponseWriter, request *nethttp.Request) {
	if request.Method != nethttp.MethodPost {
		s.metrics.RecordRequest(nethttp.StatusMethodNotAllowed, "POST required")
		writeError(writer, nethttp.StatusMethodNotAllowed, "method_not_allowed", "POST required")
		return
	}

	if !s.authorize(request) {
		s.metrics.RecordAuthFailure()
		s.metrics.RecordRequest(nethttp.StatusUnauthorized, "invalid site token")
		writeError(writer, nethttp.StatusUnauthorized, "unauthorized", "invalid site token")
		return
	}

	if !s.rateLimiter.Allow(request) {
		s.metrics.RecordRateLimited()
		s.metrics.RecordRequest(nethttp.StatusTooManyRequests, "rate limit exceeded")
		writeError(writer, nethttp.StatusTooManyRequests, "rate_limited", "rate limit exceeded")
		return
	}

	defer request.Body.Close()

	var chatRequest gateway.ChatRequest
	if err := json.NewDecoder(request.Body).Decode(&chatRequest); err != nil {
		s.metrics.RecordRequest(nethttp.StatusBadRequest, "request body must be valid JSON")
		writeError(writer, nethttp.StatusBadRequest, "invalid_json", "request body must be valid JSON")
		return
	}

	if len(chatRequest.Messages) == 0 {
		s.metrics.RecordRequest(nethttp.StatusBadRequest, "messages is required")
		writeError(writer, nethttp.StatusBadRequest, "missing_messages", "messages is required")
		return
	}

	chatRequest.ProviderToken = strings.TrimSpace(request.Header.Get("X-AIBRIDGE-PROVIDER-TOKEN"))

	if chatRequest.ProviderToken == "" {
		s.metrics.RecordRequest(nethttp.StatusBadRequest, "provider token is required")
		writeError(writer, nethttp.StatusBadRequest, "missing_provider_token", "provider token is required")
		return
	}

	requestTrafficMode := strings.TrimSpace(chatRequest.Meta.TrafficMode)
	if requestTrafficMode == "" {
		requestTrafficMode = "outbound"
	}
	if requestTrafficMode != s.config.NodeTrafficMode {
		s.metrics.RecordRequest(nethttp.StatusBadRequest, "traffic mode mismatch")
		writeError(writer, nethttp.StatusBadRequest, "traffic_mode_mismatch", "request traffic mode does not match this node")
		log.Printf(
			"chat_rejected node=%s node_mode=%s request_mode=%s provider=%s reason=traffic_mode_mismatch",
			s.config.NodeName,
			s.config.NodeTrafficMode,
			requestTrafficMode,
			chatRequest.Provider,
		)
		return
	}

	response, err := s.chatService.Chat(request.Context(), chatRequest)
	if err != nil {
		status := nethttp.StatusBadGateway
		code := "upstream_error"

		if errors.Is(err, service.ErrUnsupportedProvider) {
			status = nethttp.StatusBadRequest
			code = "unsupported_provider"
		}

		s.metrics.RecordRequest(status, err.Error())
		log.Printf(
			"chat_error node=%s node_mode=%s request_mode=%s provider=%s model=%s status=%d err=%q",
			s.config.NodeName,
			s.config.NodeTrafficMode,
			requestTrafficMode,
			chatRequest.Provider,
			chatRequest.Model,
			status,
			err.Error(),
		)
		writeError(writer, status, code, err.Error())
		return
	}

	s.metrics.RecordRequest(nethttp.StatusOK, "")
	log.Printf(
		"chat_ok node=%s node_mode=%s request_mode=%s provider=%s model=%s tokens=%d",
		s.config.NodeName,
		s.config.NodeTrafficMode,
		requestTrafficMode,
		response.Provider,
		response.Model,
		response.Usage.TotalTokens,
	)
	writeJSON(writer, nethttp.StatusOK, response)
}

type connectorProxyRequest struct {
	TargetURL string            `json:"target_url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
	Timeout   int               `json:"timeout"`
	Meta      gateway.Meta      `json:"meta"`
}

func (s *Server) handleConnectorProxy(writer nethttp.ResponseWriter, request *nethttp.Request) {
	if request.Method != nethttp.MethodPost {
		writeError(writer, nethttp.StatusMethodNotAllowed, "method_not_allowed", "POST required")
		return
	}

	if !s.authorize(request) {
		writeError(writer, nethttp.StatusUnauthorized, "unauthorized", "invalid site token")
		return
	}

	defer request.Body.Close()

	var payload connectorProxyRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, nethttp.StatusBadRequest, "invalid_json", "request body must be valid JSON")
		return
	}

	target, err := url.Parse(payload.TargetURL)
	if err != nil || target.Host == "" {
		writeError(writer, nethttp.StatusBadRequest, "invalid_target_url", "target_url is invalid")
		return
	}

	host := strings.ToLower(target.Hostname())
	if _, ok := s.config.AllowedProxyHosts[host]; !ok {
		writeError(writer, nethttp.StatusBadRequest, "target_not_allowed", "target host is not allowed")
		return
	}

	requestTrafficMode := strings.TrimSpace(payload.Meta.TrafficMode)
	if requestTrafficMode == "" {
		requestTrafficMode = "outbound"
	}
	if requestTrafficMode != s.config.NodeTrafficMode {
		writeError(writer, nethttp.StatusBadRequest, "traffic_mode_mismatch", "request traffic mode does not match this node")
		return
	}

	method := strings.ToUpper(strings.TrimSpace(payload.Method))
	if method == "" {
		method = nethttp.MethodGet
	}

	upstreamRequest, err := nethttp.NewRequestWithContext(request.Context(), method, payload.TargetURL, bytes.NewBufferString(payload.Body))
	if err != nil {
		writeError(writer, nethttp.StatusBadRequest, "build_upstream_request_failed", "could not build upstream request")
		return
	}

	for key, value := range payload.Headers {
		headerKey := nethttp.CanonicalHeaderKey(key)
		if strings.EqualFold(headerKey, "Host") || strings.EqualFold(headerKey, "Content-Length") {
			continue
		}
		upstreamRequest.Header.Set(headerKey, value)
	}

	client := &nethttp.Client{Timeout: s.config.RequestTimeout}
	if payload.Timeout > 0 {
		client.Timeout = time.Duration(payload.Timeout) * time.Second
	}

	upstreamResponse, err := client.Do(upstreamRequest)
	if err != nil {
		writeError(writer, nethttp.StatusBadGateway, "upstream_request_failed", err.Error())
		return
	}
	defer upstreamResponse.Body.Close()

	upstreamBody, err := io.ReadAll(io.LimitReader(upstreamResponse.Body, 20*1024*1024))
	if err != nil {
		writeError(writer, nethttp.StatusBadGateway, "read_upstream_failed", "could not read upstream response")
		return
	}

	responseHeaders := map[string]string{}
	for key, values := range upstreamResponse.Header {
		if len(values) > 0 {
			responseHeaders[key] = strings.Join(values, ", ")
		}
	}

	writeJSON(writer, nethttp.StatusOK, map[string]any{
		"status":  upstreamResponse.StatusCode,
		"headers": responseHeaders,
		"body":    string(upstreamBody),
	})
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}

func (s *Server) authorize(request *nethttp.Request) bool {
	header := strings.TrimSpace(request.Header.Get("Authorization"))
	if header == "" || !strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return false
	}

	token := strings.TrimSpace(header[7:])
	_, ok := s.config.SiteTokens[token]
	return ok
}

func (s *Server) withAccessLog(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(writer nethttp.ResponseWriter, request *nethttp.Request) {
		start := time.Now()
		next.ServeHTTP(writer, request)
		log.Printf("%s %s %s", request.Method, request.URL.Path, time.Since(start))
	})
}

func writeJSON(writer nethttp.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(payload)
}

func writeError(writer nethttp.ResponseWriter, status int, code, message string) {
	writeJSON(writer, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
