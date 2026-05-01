package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/omah-ti/omahtoosn/backend/internal/modules/auth"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/config"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/db"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/httpx"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/logx"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/middleware"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/security"
)

func main() {
	cfg := config.Load()
	pool, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db open failed: %v", err)
	}
	defer pool.Close()

	app := fiber.New(fiber.Config{
		AppName:               cfg.AppName,
		DisableStartupMessage: true,
		ErrorHandler:          httpx.ErrorHandler,
	})

	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(logx.Middleware())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowOrigins,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Device-ID",
		AllowMethods:     "GET,POST,PUT,OPTIONS",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return httpx.Success(c, fiber.StatusOK, "service is healthy", fiber.Map{
			"service":   cfg.AppName,
			"version":   cfg.AppVersion,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	tokenConfig := security.TokenConfig{
		AccessSecret:    "rahasia-negara-jangan-bocor",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	}
	tokenProvider := security.NewTokenProvider(tokenConfig)

	authRepo := auth.NewRepository(pool)
	authSvc := auth.NewService(authRepo, tokenProvider)
	authHandler := auth.NewHandler(authSvc)

	authMW := middleware.NewAuthMiddleware(tokenProvider)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	authHandler.RegisterRoutes(v1, authMW)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	log.Printf("service started port=%s env=%s", cfg.AppPort, cfg.AppEnv)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Printf("service stopped: %v", err)
	}
}
