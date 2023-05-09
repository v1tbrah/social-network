package rapi

import (
	"gitlab.com/pet-pr-social-network/relation-service/rpbapi"
)

type RelationAPI struct {
	relationServiceClient rpbapi.RelationServiceClient
}

func New(relationServiceClient rpbapi.RelationServiceClient) *RelationAPI {
	return &RelationAPI{
		relationServiceClient: relationServiceClient,
	}
}
