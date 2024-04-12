package token

import (
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

//func NewTokenHandler() func(next http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//
//		fn := func(w http.ResponseWriter, r *http.Request) {
//			token := r.Header.Get("token")
//
//			if token == "" {
//				resp.Send401Error(w, r)
//				slog.Error("user has no token!!!")
//				return
//			}
//			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
//			next.ServeHTTP(ww, r)
//		}
//
//		return http.HandlerFunc(fn)
//	}
//}

func NewTokenHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("token")

		if token == "" {
			slog.Error("user has no token!!!")
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized access")
		}

		return c.Next()
	}
}
