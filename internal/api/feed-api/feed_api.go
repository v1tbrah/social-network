package fapi

import (
	"github.com/v1tbrah/feed-service/fpbapi"
)

type FeedAPI struct {
	feedServiceClient fpbapi.FeedServiceClient
}

func New(feedServiceClient fpbapi.FeedServiceClient) *FeedAPI {
	return &FeedAPI{
		feedServiceClient: feedServiceClient,
	}
}
