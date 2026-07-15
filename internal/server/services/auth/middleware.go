package auth

import (
	"bird-lg/server/internal/server/platform/config"
	errx "bird-lg/server/internal/server/platform/errors"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/golang-jwt/jwt/v5"
)

const sessionCookieName = "token"

func ToolJWTMiddleware(cfg *config.Config) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.JWTSecret)},
		Extractor:  extractors.FromCookie(sessionCookieName),
		Next: func(_ fiber.Ctx) bool {
			return cfg.TurnstileSecretKey == ""
		},
		ErrorHandler: func(c fiber.Ctx, _ error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": errx.FormatPublicError(errx.ErrCodeAuthUnauthorized, "Authentication required"),
			})
		},
	})
}

func SSOJWTMiddleware(cfg *config.Config) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.JWTSecret)},
		Extractor:  extractors.FromCookie(sessionCookieName),
		SuccessHandler: func(c fiber.Ctx) error {
			token := jwtware.FromContext(c)
			if token == nil {
				return ssoAuthError(c)
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || ClaimString(claims, "auth_type") != AuthTypeLogto || ClaimString(claims, "sub") == "" {
				return ssoAuthError(c)
			}
			return c.Next()
		},
		ErrorHandler: func(c fiber.Ctx, _ error) error { return ssoAuthError(c) },
	})
}

func ssoAuthError(c fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error": errx.FormatPublicError(errx.ErrCodeAuthSSORequired, "SSO authentication required"),
	})
}

func SessionCookieName() string {
	return sessionCookieName
}
