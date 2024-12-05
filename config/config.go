package config

import (
	"crypto/rsa"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Configuration struct {
	AppEnv            string        `mapstructure:"APPLICATION_ENV"`
	AppPort           string        `mapstructure:"APPLICATION_PORT"`
	AppRateLimit      bool          `mapstructure:"APPLICATION_RATE_LIMIT"`
	AppLogRequest     bool          `mapstructure:"APPLICATION_LOG_REQUEST"`
	AppDefaultTimeout time.Duration `mapstructure:"APPLICATION_DEFAULT_TIMEOUT"`
	AppURL            string        `mapstructure:"APPLICATION_URL"`

	DbName     string `mapstructure:"DATABASE_NAME"`
	DbHost     string `mapstructure:"DATABASE_HOST"`
	DbPort     string `mapstructure:"DATABASE_PORT"`
	DbUsername string `mapstructure:"DATABASE_USERNAME"`
	DbPassword string `mapstructure:"DATABASE_PASSWORD"`
	DbSSLMode  string `mapstructure:"DATABASE_SSL_MODE"`
	DbTimeout  int    `mapstructure:"DATABASE_TIMEOUT"`

	JWTAlgorithm  string        `mapstructure:"JWT_ALGO"`
	JWTDefaultExp time.Duration `mapstructure:"JWT_DEFAULT_EXPIRATION"`
	JWTPublicKey  string        `mapstructure:"JWT_PUBLIC_KEY"`
	JWTPrivateKey string        `mapstructure:"JWT_PRIVATE_KEY"`

	RSAPublicKey  *rsa.PublicKey
	RSAPrivateKey *rsa.PrivateKey
}

func LoadConfig(path string, log *zerolog.Logger) *Configuration {
	vp := viper.New()
	vp.AddConfigPath(path)
	vp.AddConfigPath(".")
	vp.SetConfigName("app")
	vp.SetConfigType("env")
	vp.AutomaticEnv()

	if err := vp.ReadInConfig(); err != nil {
		log.Error().AnErr("Error occured while loading config file", err)
		panic(err)
	}

	var config Configuration
	if err := vp.Unmarshal(&config); err != nil {
		log.Error().AnErr("Error occured while parsing config file", err)
		panic(err)
	}

	return &config
}
