package main

import (
	"context"
	"echo-sample2/internal/todo/handler"
	"echo-sample2/internal/todo/infrastructure"
	"echo-sample2/internal/tracing"
	customValidator "echo-sample2/internal/validator"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

	e.Validator = &customValidator.CustomValidator{Validator: validator.New()}

	// TODO: DB接続情報を適切に設定する
	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect database")
	}

	todoRepo := infrastructure.NewTodoRepository(db)
	_ = todoRepo // 今後handler等で使用する

	todoHandler := handler.NewTodoHandler()
	todoHandler.RegisterTodoRoutes(e)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := e.Start(":1323"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
