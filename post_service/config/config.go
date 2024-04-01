package config

import (
	"github.com/spf13/cast"
	"os"
)

// Config ...
type Config struct {
	Environment        string // develop, staging, production
	PostgresHost       string
	PostgresPort       int
	PostgresDatabase   string
	PostgresUser       string
	PostgresPassword   string
	LogLevel           string
	RPCPort            string
	UserServiceHost    string
	UserServicePort    int
	CommentServiceHost string
	CommentServicePort int
}

// Load loads environment vars and inflates Config
func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "db"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "exam"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "nodirbek"))
	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.RPCPort = cast.ToString(getOrReturnDefault("RPC_PORT", ":5000"))

	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "user_service"))
	c.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", "8000"))

	c.CommentServiceHost = cast.ToString(getOrReturnDefault("COMMENT_SERVICE_HOST", "localhost"))
	c.CommentServicePort = cast.ToInt(getOrReturnDefault("COMMENT_SERVICE_PORT", "6000"))
	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defaultValue
}
