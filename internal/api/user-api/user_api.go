package uapi

import "github.com/v1tbrah/user-service/upbapi"

type UserAPI struct {
	userServiceClient upbapi.UserServiceClient
}

func New(userServiceClient upbapi.UserServiceClient) *UserAPI {
	return &UserAPI{
		userServiceClient: userServiceClient,
	}
}
