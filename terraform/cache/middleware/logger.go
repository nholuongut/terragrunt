package middleware

import (
	"github.com/nholuongut/terragrunt/pkg/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Logger(logger log.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger := logger.WithField("uri", v.URI).WithField("status", v.Status)
			if v.Error != nil {
				logger.Errorf("Cache server was unable to process the received request, %s", v.Error.Error())
			} else {
				logger.Tracef("Cache server received request")
			}
			return nil
		},
	})
}
