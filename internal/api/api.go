package api

import (
	"context"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"gitlab.com/pet-pr-social-network/api-gateway/internal/config"
	"gitlab.com/pet-pr-social-network/feed-service/fpbapi"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
	"gitlab.com/pet-pr-social-network/relation-service/rpbapi"
	upbapi "gitlab.com/pet-pr-social-network/user-service/pbapi"
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
			Addr: cfg.HTTPServHost + ":" + cfg.HTTPServPort,
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
			log.Fatal().Err(err).Msg("HTTP server ListenAndServe")
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
