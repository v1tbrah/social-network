package rapi

import (
	"github.com/v1tbrah/relation-service/rpbapi"
)

type RelationAPI struct {
	relationServiceClient rpbapi.RelationServiceClient
}

func New(relationServiceClient rpbapi.RelationServiceClient) *RelationAPI {
	return &RelationAPI{
		relationServiceClient: relationServiceClient,
	}
}
