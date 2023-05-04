package rapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"gitlab.com/pet-pr-social-network/api-gateway/internal/send"
	"gitlab.com/pet-pr-social-network/relation-service/rpbapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RemoveFriendReq struct {
	UserID   int64 `json:"user_id" example:"1"`
	FriendID int64 `json:"friend_id" example:"2"`
}

// RemoveFriend removes user2 from friends user1.
//
//	@Summary		Removes from friends.
//	@Description	Removes from friends.
//	@Tags			relation
//	@Produce		json
//
//	@Param			objectBody	body	RemoveFriendReq	true	"RemoveFriendReq body"
//
//	@Success		200
//	@Failure		400	{object}	send.Error
//	@Failure		500	{object}	send.Error
//	@Router			/relation/friend [delete]
func (a *RelationAPI) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	var removeFriendReq RemoveFriendReq
	if err := json.NewDecoder(r.Body).Decode(&removeFriendReq); err != nil {
		send.Send(w, send.NewErr(fmt.Sprintf("deserialize req: %v", err)), http.StatusBadRequest)
		return
	}

	_, err := a.relationServiceClient.RemoveFriend(r.Context(),
		&rpbapi.RemoveFriendRequest{UserID: removeFriendReq.UserID, FriendID: removeFriendReq.FriendID})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			send.Send(w, send.NewErr(err.Error()), http.StatusBadRequest)
			return
		}

		log.Error().Err(err).Msg("relationServiceClient.RemoveFriend")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	send.Send(w, nil, http.StatusOK)
}
