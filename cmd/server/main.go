package main

import (
	"context"
	"echo-sample2/internal/domains/todo"
	"echo-sample2/internal/tracing"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	logger := zerolog.New(os.Stdout)

	tp, tpErr := tracing.InitTracerProvider()
	if tpErr != nil {
		logger.Fatal().Err(tpErr)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Err(err)
		}
	}()

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	e := echo.New()

	e.Use(otelecho.Middleware("echo-sample2"))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))

	todoHandler := todo.NewTodoHandler()
	todoHandler.RegisterTodoRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}
