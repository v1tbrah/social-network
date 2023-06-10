package api

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/api-gateway/config"

	"github.com/v1tbrah/feed-service/fpbapi"

	"github.com/v1tbrah/post-service/ppbapi"

	"github.com/v1tbrah/relation-service/rpbapi"

	"github.com/v1tbrah/user-service/upbapi"
)

type API struct {
	server                *http.Server
	userServiceClient     upbapi.UserServiceClient
	postServiceClient     ppbapi.PostServiceClient
	relationServiceClient rpbapi.RelationServiceClient
	feedServiceClient     fpbapi.FeedServiceClient
}

// New returns new API.
func New(cfg config.Config,
	userServiceClient upbapi.UserServiceClient,
	postServiceClient ppbapi.PostServiceClient,
	relationServiceClient rpbapi.RelationServiceClient,
	feedServiceClient fpbapi.FeedServiceClient) (newAPI *API) {
	newAPI = &API{
		server: &http.Server{
			Addr: net.JoinHostPort(cfg.HTTPHost, cfg.HTTPPort),
		},
		userServiceClient:     userServiceClient,
		postServiceClient:     postServiceClient,
		relationServiceClient: relationServiceClient,
		feedServiceClient:     feedServiceClient,
	}

	newAPI.server.Handler = newAPI.newRouter()

	return newAPI
}

// StartServing starts listening the API.
func (a *API) StartServing(shutdown <-chan os.Signal) (err error) {
	ended := make(chan struct{})

	go func() {
		log.Info().Str("Addr", a.server.Addr).Msg("Starting HTTP server")
		if err = a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			err = errors.Wrapf(err, "server.ListenAndServe, addr %s", a.server.Addr)
		}
		ended <- struct{}{}
	}()

	select {
	case <-shutdown:
		return err
	case <-ended:
		return err
	}
}

// Shutdown gracefully shuts down the API without interrupting any active connections.
func (a *API) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
