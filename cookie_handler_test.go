package traefik_cookie_handler_plugin_test


import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	cookieHandler "github.com/vaspapadopoulos/traefik-cookie-handler-plugin"
)

func CookieHandlerDemo(t *testing.T) {
	cfg := cookieHandler.CreateConfig()
	cfg.Url = "https://my-api.my-domain.com"
	cfg.Method = http.MethodGet
	cfg.ResponseCookies = []string{"JWT-SESSION", "XSRF-TOKEN"}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	plugin, err := cookieHandler.New(ctx, next, cfg, "cookie-handler-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, cfg.Method, cfg.Url, nil)
	if err != nil {
		t.Fatal(err)
	}

	plugin.ServeHTTP(recorder, req)
}
