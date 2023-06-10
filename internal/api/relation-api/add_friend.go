package rapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/v1tbrah/api-gateway/internal/send"

	"github.com/v1tbrah/relation-service/rpbapi"
)

type AddFriendReq struct {
	UserID   int64 `json:"user_id" example:"1"`
	FriendID int64 `json:"friend_id" example:"2"`
}

// AddFriend adds user2 to friends user1.
//
//	@Summary		Adds to friends.
//	@Description	Adds to friends.
//	@Tags			relation
//	@Produce		json
//
//	@Param			objectBody	body	AddFriendReq	true	"AddFriendReq body"
//
//	@Success		200
//	@Failure		400	{object}	send.Error
//	@Failure		500	{object}	send.Error
//	@Router			/relation/friend [post]
func (a *RelationAPI) AddFriend(w http.ResponseWriter, r *http.Request) {
	var addFriendReq AddFriendReq
	if err := json.NewDecoder(r.Body).Decode(&addFriendReq); err != nil {
		send.Send(w, send.NewErr(fmt.Sprintf("deserialize req: %v", err)), http.StatusBadRequest)
		return
	}

	_, err := a.relationServiceClient.AddFriend(r.Context(),
		&rpbapi.AddFriendRequest{UserID: addFriendReq.UserID, FriendID: addFriendReq.FriendID})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			send.Send(w, send.NewErr(err.Error()), http.StatusBadRequest)
			return
		}

		log.Error().Err(err).Msg("relationServiceClient.AddFriend")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	send.Send(w, nil, http.StatusOK)
}
