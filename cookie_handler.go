// Package traefik_cookie_handler_plugin Traefik Cookie Handler Plugin.
package traefik_cookie_handler_plugin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Config the plugin configuration.
type Config struct {
	Url             string   `json:"url,omitempty"`
	Method          string   `json:"method,omitempty"`
	ResponseCookies []string `json:"responseCookies,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// CookieHandler a Cookie Handler plugin.
type CookieHandler struct {
	next            http.Handler
	name            string
	url             string
	method          string
	responseCookies []string
}

// New creates a new Cookie Handler plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.Url == "" {
		return nil, fmt.Errorf("URL cannot be empty")
	}
	// TODO validate url

	if config.Method == "" {
		return nil, fmt.Errorf("method cannot be empty")
	}
	// TODO validate method

	if len(config.ResponseCookies) == 0 {
		return nil, fmt.Errorf("responseCookies cannot be empty")
	}
	// TODO check if any cookie contains whitespaces or other invalid characters

	return &CookieHandler{
		next:            next,
		name:            name,
		url:             config.Url,
		method:          config.Method,
		responseCookies: config.ResponseCookies,
	}, nil
}

func (middleware *CookieHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	parsedUrl, err := url.Parse(strings.TrimSpace(middleware.url))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	mediumReq, err := http.NewRequest(middleware.method, parsedUrl.String(), nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	client := http.Client{}
	resp, err := client.Do(mediumReq)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	headerValues := resp.Header.Values("set-cookie")
	cookieHeader := ""
	for _, headerValue := range headerValues {
		cookieHeader = cookieHeader + headerValue + ";"

	}
	req.Header.Set("Cookie", cookieHeader)

	middleware.next.ServeHTTP(rw, req)
}
