package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"

	"github.com/gofiber/contrib/fiberzerolog"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/naufalatha/go-boilerplate/config"
	"github.com/naufalatha/go-boilerplate/models"
	"github.com/rs/zerolog"
)

func SetupFiberConfig(config *config.Configuration) fiber.Config {
	fiberConfig := fiber.Config{
		Immutable:    true,
		ErrorHandler: FiberErrorHandler,
		JSONEncoder:  JSONEncoder,
	}
	if config.AppEnv != models.ENV_LOCAL {
		fiberConfig.Prefork = true
		fiberConfig.DisableStartupMessage = true
	}
	return fiberConfig
}

// Custom JSON Encoder to make outgoing response match our convention
func JSONEncoder(value interface{}) ([]byte, error) {
	values := new(bytes.Buffer)
	err := json.NewEncoder(values).Encode(value)
	snakeCase := regexp.MustCompile(`\"([^\"]*?)\"\s*?:`).ReplaceAllFunc(
		values.Bytes(),
		func(match []byte) []byte {
			return bytes.ToLower(regexp.MustCompile(`(\w[^A-Z])([A-Z])`).ReplaceAll(
				match,
				[]byte(`${1}_${2}`),
			))
		},
	)
	return snakeCase, err
}

// Deafult error handling, used when handler throws unhandled error
func FiberErrorHandler(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Message:    "Internal Server Error",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	return ctx.Status(code).JSON(models.Response{
		Message:    e.Message,
		StatusCode: code,
	})
}

// Default registers list of middleware to *fiber.App
func Default(app *fiber.App, conf *config.Configuration, logger *zerolog.Logger) {
	app.Use(cors.New())
	if conf.AppLogRequest { // Request Logger
		app.Use(fiberzerolog.New(fiberzerolog.Config{
			Logger: logger,
		}))
	}
}

// UseJWT registers JWT auth guard middleware to *fiber.App
func UseJWT(app *fiber.App, conf *config.Configuration) {
	app.Use(jwtware.New(jwtware.Config{ // JWT Authentication
		SigningKey: jwtware.SigningKey{
			JWTAlg: conf.JWTAlgorithm,
			Key:    conf.RSAPublicKey,
		},
	}))
}

// middleware JWT untuk mengamankan endpoint tertentu
func JWTGuard(conf *config.Configuration) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: conf.JWTAlgorithm,
			Key:    conf.RSAPublicKey,
		},
	})
}
