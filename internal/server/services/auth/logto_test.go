package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bird-lg/server/internal/server/platform/config"
	"github.com/gofiber/fiber/v3"
)

func TestLogtoCallbackFailsClosed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		tokenStatus  int
		tokenBody    string
		userStatus   int
		userInfoBody string
	}{
		{
			name:         "token endpoint rejects code",
			tokenStatus:  http.StatusBadRequest,
			tokenBody:    `{"error":"invalid_grant"}`,
			userStatus:   http.StatusOK,
			userInfoBody: `{"sub":"user-1"}`,
		},
		{
			name:         "userinfo endpoint fails",
			tokenStatus:  http.StatusOK,
			tokenBody:    `{"access_token":"access-token"}`,
			userStatus:   http.StatusUnauthorized,
			userInfoBody: `{"error":"invalid_token"}`,
		},
		{
			name:         "userinfo has no subject",
			tokenStatus:  http.StatusOK,
			tokenBody:    `{"access_token":"access-token"}`,
			userStatus:   http.StatusOK,
			userInfoBody: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/oidc/token":
					w.WriteHeader(tt.tokenStatus)
					_, _ = w.Write([]byte(tt.tokenBody))
				case "/oidc/me":
					w.WriteHeader(tt.userStatus)
					_, _ = w.Write([]byte(tt.userInfoBody))
				default:
					http.NotFound(w, r)
				}
			}))
			defer provider.Close()

			cfg := &config.Config{
				LogtoEndpoint: provider.URL,
				LogtoAppID:    "app-id",
				JWTSecret:     "test-secret",
			}
			app := fiber.New()
			app.Get("/auth/callback", HandleLogtoCallback(cfg))

			req := httptest.NewRequest(http.MethodGet, "http://example.test/auth/callback?code=code&state=Lw", nil)
			req.AddCookie(&http.Cookie{Name: "logto_code_verifier", Value: "verifier"})
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test() error = %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != fiber.StatusUnauthorized {
				t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusUnauthorized)
			}
			for _, cookie := range resp.Cookies() {
				if cookie.Name == SessionCookieName() && cookie.Value != "" {
					t.Fatalf("callback issued an SSO session after provider failure")
				}
			}
		})
	}
}

func TestSanitizeRedirectPath(t *testing.T) {
	t.Parallel()

	for _, unsafe := range []string{"//example.com", `/\\example.com`, "/path\nnext"} {
		if got := sanitizeRedirectPath(unsafe); got != "/" {
			t.Fatalf("sanitizeRedirectPath(%q) = %q, want /", unsafe, got)
		}
	}
}
