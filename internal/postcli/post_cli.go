package postcli

import (
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gitlab.com/pet-pr-social-network/api-gateway/config"
)

func NewConn(cfg config.PostCli) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(net.JoinHostPort(cfg.Host, cfg.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrapf(err, "grpc.Dial, addr %s", net.JoinHostPort(cfg.Host, cfg.Port))
	}

	return conn, nil
}
