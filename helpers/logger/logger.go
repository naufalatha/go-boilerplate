package logger

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/naufalatha/go-boilerplate/models"
	"github.com/rs/zerolog"
)

const (
	ErrorTag = "[ERROR]"
	DebugTag = "[DEBUG]"
	TraceTag = "[TRACE]"
	InfoTag  = "[INFO]"
)

type Logger struct {
	log *zerolog.Logger
	env string
}

var logger Logger

func InitLogger(zerolog *zerolog.Logger, env string) {
	logger = Logger{
		log: zerolog,
		env: env,
	}
}

func handleArguments(args ...any) string {
	var arguments string
	for _, arg := range args {
		arguments += fmt.Sprintf("%v ", arg)
	}
	return arguments
}

func Debug(message string, args ...any) {
	fmt.Printf("%s \033[36m%s\033[0m %s %v\n", time.Now().Format(models.DEFAULT_TIME_FORMAT), DebugTag, message, handleArguments(args...))
}

func Trace(message string, args ...any) {
	logger.log.Trace().Caller(1).Msg(fmt.Sprint(message, handleArguments(args...)))
}

func TraceCtx(ctx *fiber.Ctx, message string, args ...any) {
	claim := getClaim(ctx)
	prefix := logger.log.Trace().Interface("jwt_id", claim["jti"])
	if claim["sub"] != nil {
		prefix.Interface("user_id", claim["sub"])
	}
	prefix.Msg(fmt.Sprint(message, handleArguments(args...)))
}

func Info(message string, args ...any) {
	logger.log.Info().Msg(fmt.Sprint(message, handleArguments(args...)))
}

func InfoCtx(ctx *fiber.Ctx, message string, args ...any) {
	claim := getClaim(ctx)
	prefix := logger.log.Info().Interface("jwt_id", claim["jti"])
	if claim["sub"] != nil {
		prefix.Interface("user_id", claim["sub"])
	}
	prefix.Msg(fmt.Sprint(message, handleArguments(args...)))
}

func Error(message string, err ...any) {
	logger.log.Error().Caller(1).Msg(fmt.Sprint(message, handleArguments(err...)))
}
func ErrorCaller(message string, caller int, err ...any) {
	logger.log.Error().Caller(caller).Msg(fmt.Sprint(message, handleArguments(err...)))
}

func ErrorStack(message string, err ...any) {
	logger.log.Error().Caller(1).Msg(fmt.Sprint(message, handleArguments(err...), string(debug.Stack())))
}

func ErrorCtx(ctx *fiber.Ctx, message string, err ...any) {
	claim := getClaim(ctx)
	prefix := logger.log.Error().Caller(1).Interface("jwt_id", claim["jti"])
	if claim["sub"] != nil {
		prefix.Interface("customer_id", claim["sub"])
	}
	prefix.Msg(fmt.Sprint(message, handleArguments(err...)))
}

func getClaim(c *fiber.Ctx) jwt.MapClaims {
	user := c.Locals("user")
	claims := jwt.MapClaims{}
	if user != nil {
		token := user.(*jwt.Token)
		claims = token.Claims.(jwt.MapClaims)
	}
	return claims
}
