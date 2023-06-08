package usercli

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gitlab.com/pet-pr-social-network/api-gateway/config"
)

func NewConn(cfg config.UserCli) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(cfg.Host+":"+cfg.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial: %w", err)
	}

	return conn, nil
}
