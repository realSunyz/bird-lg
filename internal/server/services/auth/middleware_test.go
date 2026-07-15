package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bird-lg/server/internal/server/platform/config"
	"github.com/gofiber/fiber/v3"
)

func TestSSOJWTMiddlewareRequiresSSOIdentity(t *testing.T) {
	t.Parallel()

	const secret = "test-secret"
	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{name: "captcha token", token: GenerateJWT(secret), wantStatus: fiber.StatusForbidden},
		{name: "SSO token without subject", token: GenerateJWTWithSub(secret, AuthTypeLogto, ""), wantStatus: fiber.StatusForbidden},
		{name: "SSO token with subject", token: GenerateJWTWithSub(secret, AuthTypeLogto, "user-1"), wantStatus: fiber.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			app := fiber.New()
			app.Get("/", SSOJWTMiddleware(&config.Config{JWTSecret: secret}), func(c fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "http://example.test/", nil)
			req.AddCookie(&http.Cookie{Name: SessionCookieName(), Value: tt.token})
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test() error = %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}
