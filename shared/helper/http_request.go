package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultTimeout = 30 * time.Second

type (
	HTTPClient struct {
		client  *http.Client
		baseURL string
		headers map[string]string
	}

	RequestOptions struct {
		Context     context.Context
		Headers     map[string]string
		QueryParams map[string]string
	}

	HttpResponse struct {
		StatusCode int
		Body       []byte
		Headers    http.Header
	}

	FormData struct {
		Fields map[string]string
		Files  []*multipart.FileHeader
	}
)

type RequestOption func(*RequestOptions)

func NewHTTPClient(baseURL string, opts ...func(*http.Client)) *HTTPClient {
	client := &http.Client{
		Timeout: defaultTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(client)
	}
	return &HTTPClient{
		client:  client,
		baseURL: strings.TrimRight(baseURL, "/"),
		headers: make(map[string]string),
	}
}

func (c *HTTPClient) WithHeader(key, value string) *HTTPClient {
	c.headers[key] = value
	return c
}

func (c *HTTPClient) WithHeaders(headers map[string]string) *HTTPClient {
	maps.Copy(c.headers, headers)
	return c
}

func (c *HTTPClient) Get(path string, opts ...RequestOption) (*HttpResponse, error) {
	return c.doRequest(http.MethodGet, path, nil, opts...)
}

func (c *HTTPClient) Post(path string, body any, opts ...RequestOption) (*HttpResponse, error) {
	return c.sendBody(http.MethodPost, path, body, opts...)
}

func (c *HTTPClient) Patch(path string, body any, opts ...RequestOption) (*HttpResponse, error) {
	return c.sendBody(http.MethodPatch, path, body, opts...)
}

func (c *HTTPClient) Delete(path string, opts ...RequestOption) (*HttpResponse, error) {
	return c.doRequest(http.MethodDelete, path, nil, opts...)
}

func (c *HTTPClient) Upload(method, path string, form FormData, opts ...RequestOption) (*HttpResponse, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for k, v := range form.Fields {
		if err := writer.WriteField(k, v); err != nil {
			return nil, fmt.Errorf("write field %s: %w", k, err)
		}
	}

	for _, file := range form.Files {
		f, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("open file %s: %w", file.Filename, err)
		}
		defer f.Close()

		part, err := writer.CreateFormFile(file.Filename, file.Filename)
		if err != nil {
			return nil, fmt.Errorf("create form file %s: %w", file.Filename, err)
		}
		if _, err := io.Copy(part, f); err != nil {
			return nil, fmt.Errorf("copy file %s: %w", file.Filename, err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	opts = append([]RequestOption{
		WithHeaders(map[string]string{
			"Content-Type": writer.FormDataContentType(),
		}),
	}, opts...)

	return c.doRequest(method, path, buf.Bytes(), opts...)
}

func WithHeaders(headers map[string]string) RequestOption {
	return func(o *RequestOptions) {
		if o.Headers == nil {
			o.Headers = make(map[string]string)
		}
		maps.Copy(o.Headers, headers)
	}
}

func WithQueryParams(params map[string]string) RequestOption {
	return func(o *RequestOptions) {
		if o.QueryParams == nil {
			o.QueryParams = make(map[string]string)
		}
		maps.Copy(o.QueryParams, params)
	}
}

func WithContext(ctx context.Context) RequestOption {
	return func(o *RequestOptions) {
		o.Context = ctx
	}
}

func (c *HTTPClient) sendBody(method, path string, body any, opts ...RequestOption) (*HttpResponse, error) {
	var (
		payload []byte
		err     error
	)

	switch v := body.(type) {
	case nil:
	case []byte:
		payload = v
	case string:
		payload = []byte(v)
	default:
		payload, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		opts = append([]RequestOption{
			WithHeaders(map[string]string{"Content-Type": "application/json"}),
		}, opts...)
	}

	return c.doRequest(method, path, payload, opts...)
}

func (c *HTTPClient) doRequest(method, path string, body []byte, opts ...RequestOption) (*HttpResponse, error) {
	settings := &RequestOptions{
		Context:     context.Background(),
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}
	for _, apply := range opts {
		apply(settings)
	}

	reqURL := c.buildURL(path, settings.QueryParams)

	req, err := http.NewRequestWithContext(settings.Context, method, reqURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	for k, v := range settings.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	return &HttpResponse{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}, nil
}

func (c *HTTPClient) buildURL(path string, query map[string]string) string {
	u := c.baseURL + path
	if len(query) == 0 {
		return u
	}

	params := url.Values{}
	for k, v := range query {
		params.Set(k, v)
	}
	if strings.Contains(u, "?") {
		return u + "&" + params.Encode()
	}
	return u + "?" + params.Encode()
}

func (r *HttpResponse) JSON(v any) error {
	if len(r.Body) == 0 {
		return fmt.Errorf("empty response body")
	}
	return json.Unmarshal(r.Body, v)
}

func (r *HttpResponse) String() string {
	return string(r.Body)
}

func (r *HttpResponse) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

func (r *HttpResponse) Error() error {
	if r.IsSuccess() {
		return nil
	}
	return fmt.Errorf("status %d: %s", r.StatusCode, r.String())
}
