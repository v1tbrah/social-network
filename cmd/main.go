package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/pet-pr-social-network/api-gateway/internal/api"
	"gitlab.com/pet-pr-social-network/api-gateway/internal/config"
	psclient "gitlab.com/pet-pr-social-network/api-gateway/internal/post-service-client"
	rsclient "gitlab.com/pet-pr-social-network/api-gateway/internal/relation-service-client"
	usclient "gitlab.com/pet-pr-social-network/api-gateway/internal/user-service-client"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
	"gitlab.com/pet-pr-social-network/relation-service/rpbapi"
	upbapi "gitlab.com/pet-pr-social-network/user-service/pbapi"
)

func main() {
	newConfig := config.NewDefaultConfig()
	zerolog.SetGlobalLevel(newConfig.LogLvl)

	if err := newConfig.ParseEnv(); err != nil {
		log.Fatal().Err(err).Msg("config.ParseEnv")
	}
	zerolog.SetGlobalLevel(newConfig.LogLvl)

	userServiceConn, err := usclient.NewConn(newConfig.UserServiceClient)
	if err != nil {
		log.Fatal().Err(err).Msg("usclient.NewConn")
	}
	newUserServiceClient := upbapi.NewUserServiceClient(userServiceConn)

	postServiceConn, err := psclient.NewConn(newConfig.PostServiceClient)
	if err != nil {
		log.Fatal().Err(err).Msg("usclient.NewConn")
	}
	newPostServiceClient := ppbapi.NewPostServiceClient(postServiceConn)

	relationServiceConn, err := rsclient.NewConn(newConfig.RelationServiceClient)
	if err != nil {
		log.Fatal().Err(err).Msg("rsclient.NewConn")
	}
	newRelationServiceClient := rpbapi.NewRelationServiceClient(relationServiceConn)

	newAPI := api.New(newConfig, newUserServiceClient, newPostServiceClient, newRelationServiceClient)

	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	errServingCh := make(chan error)
	go func() {
		errServing := newAPI.StartServing(shutdownSig)
		errServingCh <- errServing
	}()

	select {
	case <-shutdownSig:
		close(shutdownSig)
	case errServing := <-errServingCh:
		if errServing != nil {
			log.Error().Err(errServing).Msg("newAPI.StartServing")
		}
	}

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownCtxCancel()
	if err = newAPI.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("newAPI.Shutdown")
	} else {
		log.Info().Msg("HTTP server gracefully stopped")
	}
	if err = userServiceConn.Close(); err != nil {
		log.Error().Err(err).Msg("userServiceConn.Close")
	}
	if err = postServiceConn.Close(); err != nil {
		log.Error().Err(err).Msg("postServiceConn.Close")
	}
	if err = relationServiceConn.Close(); err != nil {
		log.Error().Err(err).Msg("relationServiceConn.Close")
	}
}
