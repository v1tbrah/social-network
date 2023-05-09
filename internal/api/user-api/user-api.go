package uapi

import "gitlab.com/pet-pr-social-network/user-service/pbapi"

type UserAPI struct {
	userServiceClient pbapi.UserServiceClient
}

func New(userServiceClient pbapi.UserServiceClient) *UserAPI {
	return &UserAPI{
		userServiceClient: userServiceClient,
	}
}
