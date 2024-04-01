package config

import (
	"github.com/spf13/cast"
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
	PostServiceHost    string
	PostServicePort    int
	CommentServiceHost string
	CommentServicePort int
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
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

	c.RPCPort = cast.ToString(getOrReturnDefault("RPC_PORT", ":8000"))

	c.PostServiceHost = cast.ToString(getOrReturnDefault("POST_SERVICE_HOST", "post_service"))
	c.PostServicePort = cast.ToInt(getOrReturnDefault("POST_SERVICE_HOST", "5000"))

	c.CommentServiceHost = cast.ToString(getOrReturnDefault("COMMENT_SERVICE_HOST", "comment_service"))
	c.CommentServicePort = cast.ToInt(getOrReturnDefault("COMMENT_SERVICE_HOST", "6000"))

	c.RedisAddr = cast.ToString(getOrReturnDefault("REDIS_ADDR", "redis:6379"))
	c.RedisPassword = cast.ToString(getOrReturnDefault("REDIS_PASSWORD", ""))
	c.RedisDB = cast.ToInt(getOrReturnDefault("REDIS_DB", 0))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	// _, exists := os.LookupEnv(key)
	// if exists {
	//     return os.Getenv(key)
	// }

	return defaultValue
}
