package middleware

import (
	"chatwithme/logger"

	"github.com/labstack/echo/v4"
)

func Logger(log logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			log.Noticef("%v from %v", req.Method, req.RemoteAddr)
			if err := next(c); err != nil {
				log.Errorf("%v; %v", err, c.Get("error"))
				return err
			}

			return nil
		}
	}
}
