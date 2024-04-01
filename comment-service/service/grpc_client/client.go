package grpcclient

import "4microservices/comment-service/config"

//GrpcClientI

type GrpcClient interface {
}

// GrpcClient
type GrpClient struct {
	cfg         config.Config
	connections map[string]interface{}
}

//New Grpc Client

func New(cfg config.Config) (*GrpClient, error) {
	return &GrpClient{
		cfg:         cfg,
		connections: map[string]interface{}{},
	}, nil
}
