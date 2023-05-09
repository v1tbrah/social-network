package papi

import "gitlab.com/pet-pr-social-network/post-service/ppbapi"

type PostAPI struct {
	postServiceClient ppbapi.PostServiceClient
}

func New(postServiceClient ppbapi.PostServiceClient) *PostAPI {
	return &PostAPI{
		postServiceClient: postServiceClient,
	}
}
