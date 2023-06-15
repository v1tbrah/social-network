package papi

import "github.com/v1tbrah/post-service/ppbapi"

type PostAPI struct {
	postServiceClient ppbapi.PostServiceClient
}

func New(postServiceClient ppbapi.PostServiceClient) *PostAPI {
	return &PostAPI{
		postServiceClient: postServiceClient,
	}
}
