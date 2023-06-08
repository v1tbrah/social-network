package uapi

import "gitlab.com/pet-pr-social-network/user-service/upbapi"

type UserAPI struct {
	userServiceClient upbapi.UserServiceClient
}

func New(userServiceClient upbapi.UserServiceClient) *UserAPI {
	return &UserAPI{
		userServiceClient: userServiceClient,
	}
}
