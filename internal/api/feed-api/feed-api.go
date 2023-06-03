package fapi

import (
	"gitlab.com/pet-pr-social-network/feed-service/fpbapi"
)

type FeedAPI struct {
	feedServiceClient fpbapi.FeedServiceClient
}

func New(feedServiceClient fpbapi.FeedServiceClient) *FeedAPI {
	return &FeedAPI{
		feedServiceClient: feedServiceClient,
	}
}
