package mediacli

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/v1tbrah/api-gateway/config"
)

func NewConn(cfg config.MediaCli) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(net.JoinHostPort(cfg.Host, cfg.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial: %w", err)
	}

	return conn, nil
}
