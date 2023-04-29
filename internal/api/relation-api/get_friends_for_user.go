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

type GetFriendsByUserReq struct {
	FriendID     int64     `json:"friend_id" example:"1"`
	Direction    Direction `json:"direction" example:"1"`
	PostOffsetID int64     `json:"postOffsetID" example:"11"`
	Limit        int64     `json:"limit" example:"100"`
}

type Direction int32

const (
	First Direction = iota
	Next
	Prev
)

// GetFriendsByUser returns friends by user.
//
//	@Summary		Returns friends by user.
//	@Description	Returns friends by user.
//	@Tags			relation
//	@Produce		json
//
//	@Param			objectBody	body		GetFriendsByUserReq	true	"GetFriendsByUserReq body"
//
//	@Success		200			{object}	[]int64
//	@Success		204
//	@Failure		400	{object}	send.Error
//	@Failure		500	{object}	send.Error
//	@Router			/relation/friend/get_friends_by_user [post]
func (a *RelationAPI) GetFriendsByUser(w http.ResponseWriter, r *http.Request) {
	var getFriendsReq GetFriendsByUserReq
	if err := json.NewDecoder(r.Body).Decode(&getFriendsReq); err != nil {
		send.Send(w, send.NewErr(fmt.Sprintf("deserialize req: %v", err)), http.StatusBadRequest)
		return
	}

	pbGetFriendsResp, err := a.relationServiceClient.GetFriends(r.Context(), &rpbapi.GetFriendsRequest{
		UserID: getFriendsReq.FriendID, Direction: rpbapi.GetFriendsRequest_DIRECTION(getFriendsReq.Direction),
		UserOffsetID: getFriendsReq.PostOffsetID, Limit: getFriendsReq.Limit})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			send.Send(w, send.NewErr(err.Error()), http.StatusBadRequest)
			return
		}
		if status.Code(err) == codes.NotFound {
			send.Send(w, send.NewErr(err.Error()), http.StatusNoContent)
			return
		}

		log.Error().Err(err).Msg("postServiceClient.GetPostsByHashtag")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	if len(pbGetFriendsResp.GetFriends()) == 0 {
		send.Send(w, nil, http.StatusNoContent)
		return
	}

	send.Send(w, pbGetFriendsResp.GetFriends(), http.StatusOK)
}
