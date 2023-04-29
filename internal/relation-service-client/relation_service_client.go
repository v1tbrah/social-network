package relation_service_client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gitlab.com/pet-pr-social-network/api-gateway/internal/config"
)

func NewConn(cfg config.RelationServiceClient) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(cfg.ServHost+":"+cfg.ServPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial: %w", err)
	}

	return conn, nil
}
