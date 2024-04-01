package v1

import (
	"4microservice/api_gateway/api/tokens"
	"4microservice/api_gateway/config"
	"4microservice/api_gateway/pkg/logger"
	"4microservice/api_gateway/services"
	"4microservice/api_gateway/storage/repo"

	"github.com/nats-io/nats.go"
)

type HandlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	jwthandler     tokens.JWTHandler
	cfg            config.Config
	reds           repo.RedisStorageI
	natsConn       *nats.Conn
}

type HandlerV1Config struct {
	Log            logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Reds           repo.RedisStorageI
	NatsConn       *nats.Conn
}

func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		log:            c.Log,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		reds:           c.Reds,
		natsConn:       c.NatsConn,
	}
}
