package mapi

import (
	"github.com/v1tbrah/media-service/mpbapi"
)

type MediaAPI struct {
	mediaServiceClient mpbapi.MediaServiceClient
}

func New(mediaServiceClient mpbapi.MediaServiceClient) *MediaAPI {
	return &MediaAPI{
		mediaServiceClient: mediaServiceClient,
	}
}
