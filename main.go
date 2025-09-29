package main

import (
	"embed"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"superqi-proxy/api"
	"superqi-proxy/superqi"
	"superqi-proxy/util"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
)

//go:embed res/frontend/*
var frontendFS embed.FS

// localhostOnly middleware restricts access to localhost only
func localhostOnly(c *fiber.Ctx) error {
	clientIP := c.IP()

	// Parse the client IP
	ip := net.ParseIP(clientIP)
	if ip == nil {
		// If we can't parse the IP, check if it's a string representation
		if clientIP == "127.0.0.1" || clientIP == "::1" || strings.HasPrefix(clientIP, "127.") {
			return c.Next()
		}
		slog.Debug("Access denied to /test endpoint", "ip", clientIP, "reason", "invalid IP format")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	// Check if IP is localhost
	if ip.IsLoopback() || clientIP == "127.0.0.1" || clientIP == "::1" {
		return c.Next()
	}

	slog.Debug("Access denied to /test endpoint", "ip", clientIP, "reason", "not localhost")
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
}

func main() {
	writer := os.Stdout
	slog.New(tint.NewHandler(writer, nil))
	slog.SetDefault(slog.New(
		tint.NewHandler(writer, &tint.Options{
			AddSource:  true,
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
		}),
	))

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := superqi.InitSuperQiClient(); err != nil {
		log.Fatal("Error init SuperQi client", err)
	}

	startWebServer(writer)
}

func startWebServer(writer *os.File) {
	app := fiber.New(fiber.Config{
		AppName:               "SuperQi-Proxy",
		ServerHeader:          "superqi-proxy",
		StrictRouting:         true,
		CaseSensitive:         true,
		Concurrency:           512 * 1024,       // default: 256 * 1024
		ReadTimeout:           60 * time.Second, // default: unlimited
		WriteTimeout:          60 * time.Second, // default: unlimited
		DisableKeepalive:      true,
		DisableStartupMessage: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
				slog.Debug("Caught response error:", slog.String("Message", e.Message))
			} else {
				slog.Warn("Caught generic error:", "Err", err)
			}

			return ctx.Status(code).JSON(fiber.Map{"message": "server error, check logs"})
		},
	})

	app.Use("/_", localhostOnly)
	app.Use("/_", filesystem.New(filesystem.Config{
		Root:       http.FS(frontendFS),
		PathPrefix: "res/frontend",
		Browse:     false,
		Index:      "index.html",
	}))

	app.Use(requestid.New())
	app.Use(recover2.New())
	app.Use(healthcheck.New())
	app.Use(logger.New(logger.Config{
		Output:       writer,
		TimeInterval: 500 * time.Millisecond,
		Format:       "${time} ${method} [${status}]${latency} - ${path}\n",
		TimeFormat:   time.DateTime,
	}))

	// API routes
	apiGroup := app.Group("/api")
	api.InitAPI(apiGroup)

	port := util.GetEnvOrDefaultString("PORT", "1999")
	slog.Info("Starting App...", "port", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
