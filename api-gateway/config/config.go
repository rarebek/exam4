package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Environment string //develop, staging, production

	UserServiceHost string
	UserServicePort int

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	PostServiceHost string
	PostServicePort int

	CommentServiceHost string
	CommentServicePort int

	LikeServiceHost string
	LikeServicePort int

	CtxTimeOut int
	RedisPort  int
	LogLevel   string
	HTTPPort   string

	SignInKey          string
	AccessTokenTimeout int
}

func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":5555"))

	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "user_service"))
	c.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", 8000))
	c.PostServiceHost = cast.ToString(getOrReturnDefault("POST_SERVICE_HOST", "post_service"))
	c.PostServicePort = cast.ToInt(getOrReturnDefault("POST_SERVICE_PORT", 5000))
	c.CommentServiceHost = cast.ToString(getOrReturnDefault("COMMENT_SERVICE_HOST", "comment_service"))
	c.CommentServicePort = cast.ToInt(getOrReturnDefault("COMMENT_SERVICE_PORT", 6000))
	c.RedisAddr = cast.ToString(getOrReturnDefault("REDIS_ADDR", "redis"))
	c.RedisPort = cast.ToInt(getOrReturnDefault("REFIS_PORT", 6379))
	c.RedisPassword = cast.ToString(getOrReturnDefault("REDIS_PASSWORD", ""))
	c.RedisDB = cast.ToInt(getOrReturnDefault("REDIS_DB", 0))
	c.SignInKey = cast.ToString(getOrReturnDefault("SIGN_IN_KEY", "nodirbek"))
	c.AccessTokenTimeout = cast.ToInt(getOrReturnDefault("ACCESS_TOKEN_TIMEOUT", 15))

	c.CtxTimeOut = cast.ToInt(getOrReturnDefault("CTX_TIMEOUT", 7))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
