package papi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/pet-pr-social-network/api-gateway/internal/send"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
)

type AddHashtagToPostReq struct {
	PostID    int64 `json:"post_id" example:"1"`
	HashtagID int64 `json:"hashtag_id" example:"1"`
}

// AddHashtagToPost creates hashtag.
//
//	@Summary		Adds hashtag to post.
//	@Description	Adds hashtag to post.
//	@Tags			hashtag
//	@Produce		json
//
//	@Param			objectBody	body	AddHashtagToPostReq	true	"AddHashtagToPostReq body"
//
//	@Success		200
//	@Failure		400	{object}	send.Error
//	@Failure		500	{object}	send.Error
//	@Router			/post/hashtag [post]
func (a *PostAPI) AddHashtagToPost(w http.ResponseWriter, r *http.Request) {
	var addHashtagToPostReq AddHashtagToPostReq
	if err := json.NewDecoder(r.Body).Decode(&addHashtagToPostReq); err != nil {
		send.Send(w, send.NewErr(fmt.Sprintf("deserialize req: %v", err)), http.StatusBadRequest)
		return
	}

	_, err := a.postServiceClient.AddHashtagToPost(r.Context(), &ppbapi.AddHashtagToPostRequest{
		PostID: addHashtagToPostReq.PostID, HashtagID: addHashtagToPostReq.HashtagID})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			send.Send(w, send.NewErr(err.Error()), http.StatusBadRequest)
			return
		}

		log.Error().Err(err).Msg("postServiceClient.AddHashtagToPost")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	send.Send(w, nil, http.StatusOK)
}
