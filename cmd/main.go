package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/api-gateway/config"
	"github.com/v1tbrah/api-gateway/internal/api"
	"github.com/v1tbrah/api-gateway/internal/feedcli"
	"github.com/v1tbrah/api-gateway/internal/mediacli"
	"github.com/v1tbrah/api-gateway/internal/postcli"
	"github.com/v1tbrah/api-gateway/internal/relationcli"
	"github.com/v1tbrah/api-gateway/internal/usercli"

	"github.com/v1tbrah/feed-service/fpbapi"
	"github.com/v1tbrah/media-service/mpbapi"
	"github.com/v1tbrah/post-service/ppbapi"
	"github.com/v1tbrah/relation-service/rpbapi"
	"github.com/v1tbrah/user-service/upbapi"
)

func main() {
	newConfig := config.NewDefaultConfig()
	zerolog.SetGlobalLevel(newConfig.LogLvl)

	if err := newConfig.ParseEnv(); err != nil {
		log.Fatal().Err(err).Msg("config.ParseEnv")
	}
	zerolog.SetGlobalLevel(newConfig.LogLvl)

	userServiceConn, err := usercli.NewConn(newConfig.UserCli)
	if err != nil {
		log.Fatal().Err(err).Msg("usclient.NewConn")
	}
	newUserServiceClient := upbapi.NewUserServiceClient(userServiceConn)

	postServiceConn, err := postcli.NewConn(newConfig.PostCli)
	if err != nil {
		log.Fatal().Err(err).Msg("usclient.NewConn")
	}
	newPostServiceClient := ppbapi.NewPostServiceClient(postServiceConn)

	relationServiceConn, err := relationcli.NewConn(newConfig.RelationCli)
	if err != nil {
		log.Fatal().Err(err).Msg("rsclient.NewConn")
	}
	newRelationServiceClient := rpbapi.NewRelationServiceClient(relationServiceConn)

	feedServiceConn, err := feedcli.NewConn(newConfig.FeedCli)
	if err != nil {
		log.Fatal().Err(err).Msg("fsclient.NewConn")
	}
	newFeedServiceClient := fpbapi.NewFeedServiceClient(feedServiceConn)

	mediaServiceConn, err := mediacli.NewConn(newConfig.MediaCli)
	if err != nil {
		log.Fatal().Err(err).Msg("msclient.NewConn")
	}
	newMediaServiceClient := mpbapi.NewMediaServiceClient(mediaServiceConn)

	newAPI := api.New(newConfig,
		newUserServiceClient,
		newPostServiceClient,
		newRelationServiceClient,
		newFeedServiceClient,
		newMediaServiceClient,
	)

	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	errServingCh := make(chan error)
	go func() {
		errServing := newAPI.StartServing(shutdownSig)
		errServingCh <- errServing
	}()

	select {
	case shutdownSigValue := <-shutdownSig:
		close(shutdownSig)
		log.Info().Msgf("Shutdown signal received: %s", strings.ToUpper(shutdownSigValue.String()))
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
	if err = feedServiceConn.Close(); err != nil {
		log.Error().Err(err).Msg("feedServiceConn.Close")
	}
	if err = mediaServiceConn.Close(); err != nil {
		log.Error().Err(err).Msg("mediaServiceConn.Close")
	}
}
