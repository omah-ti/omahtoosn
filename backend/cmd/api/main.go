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
	swagger "github.com/gofiber/swagger"

	docs "github.com/omah-ti/omahtoosn/backend/docs"
	"github.com/omah-ti/omahtoosn/backend/internal/modules/auth"
	"github.com/omah-ti/omahtoosn/backend/internal/modules/tryout"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/config"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/db"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/httpx"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/logx"
	appmiddleware "github.com/omah-ti/omahtoosn/backend/internal/platform/middleware"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/security"
)

//go:generate swag init -g cmd/api/main.go -d ../.. -o ../../docs

// @title TO OSN Backend API
// @version 1.0
// @description Dokumentasi API backend TO OSN.
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cfg := config.Load()

	docs.SwaggerInfo.Title = cfg.AppName + " API"
	docs.SwaggerInfo.Version = cfg.AppVersion
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// db connection
	pool, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db open failed: %v", err)
	}
	defer pool.Close()

	// security (jwt)
	tokenConfig := security.TokenConfig{
		AccessSecret:    cfg.JWTSecret,
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	}
	tokenProvider := security.NewTokenProvider(tokenConfig)

	// auth module
	authRepo := auth.NewRepository(pool)
	authService := auth.NewService(authRepo, tokenProvider)
	authHandler := auth.NewHandler(authService)

	// tryout module
	tryoutRepo := tryout.NewRepository()
	tryoutService := tryout.NewService(pool, tryoutRepo)
	tryoutHandler := tryout.NewHandler(tryoutService)

	// app init
	app := fiber.New(fiber.Config{
		AppName:               cfg.AppName,
		DisableStartupMessage: true,
		ErrorHandler:          httpx.ErrorHandler,
	})

	// middleware
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(logx.Middleware())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowOrigins,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Device-ID",
		AllowMethods:     "GET,POST,PUT,OPTIONS",
	}))

	app.Get("/health", healthHandler(cfg))
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Static("/question-assets", "./assets/questions")

	v1 := app.Group("/api/v1")

	authGroup := v1.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/refresh", authHandler.RefreshToken)

	authRequired := appmiddleware.NewAuthMiddleware(tokenProvider)
	protected := v1.Group("", authRequired)

	protected.Post("/auth/logout", authHandler.Logout)
	protected.Post("/auth/logout-all", authHandler.LogoutAll)
	protected.Get("/me", authHandler.Me)

	protected.Get("/tryouts/current", tryoutHandler.GetCurrentTryout)
	protected.Post("/tryouts/current/start", tryoutHandler.StartCurrentTryout)
	protected.Get("/attempts/current", tryoutHandler.GetCurrentAttempt)
	protected.Put("/attempts/current/answers", tryoutHandler.SaveCurrentAnswers)
	protected.Post("/attempts/current/submit", tryoutHandler.SubmitCurrentAttempt)
	protected.Get("/results/current", tryoutHandler.GetCurrentResult)
	protected.Get("/leaderboard/current", tryoutHandler.GetCurrentLeaderboard)

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
