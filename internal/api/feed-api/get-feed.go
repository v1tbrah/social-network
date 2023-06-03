package fapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/pet-pr-social-network/api-gateway/internal/send"
	"gitlab.com/pet-pr-social-network/feed-service/fpbapi"
)

type Post struct {
	Id          int64     `json:"id" example:"1"`
	UserID      int64     `json:"user_id" example:"1"`
	Description string    `json:"description" example:"description"`
	HashtagsID  []int64   `json:"hashtags_id" example:"1,2,3"`
	CreatedAt   time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
}

// GetFeed returns posts by user feed.
//
//	@Summary		Returns posts by user feed.
//	@Description	Returns posts by user feed.
//	@Tags			feed
//	@Produce		json
//
//	@Param			id	path		int	true	"User id"
//
//	@Success		200	{object}	[]Post
//	@Success		204
//	@Failure		400	{object}	send.Error
//	@Failure		404	{object}	send.Error
//	@Failure		500	{object}	send.Error
//	@Router			/feed/{id} [get]
func (a *FeedAPI) GetFeed(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		send.Send(w, send.NewErr(errInvalidID.Error()), http.StatusBadRequest)
		return
	}

	pbGetFeedResp, err := a.feedServiceClient.GetFeed(r.Context(), &fpbapi.GetFeedRequest{UserID: id})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			send.Send(w, send.NewErr(err.Error()), http.StatusNotFound)
			return
		}
		log.Error().Err(err).Msg("feedServiceClient.GetFeed")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	if len(pbGetFeedResp.GetPosts()) == 0 {
		send.Send(w, nil, http.StatusNoContent)
		return
	}

	resp := make([]Post, 0, len(pbGetFeedResp.GetPosts()))
	for _, post := range pbGetFeedResp.GetPosts() {
		resp = append(resp, Post{
			Id:          post.GetId(),
			UserID:      post.GetUserID(),
			Description: post.GetDescription(),
			HashtagsID:  post.GetHashtagsID(),
			CreatedAt:   post.GetCreatedAt().AsTime(),
		})
	}
	send.Send(w, resp, http.StatusOK)
}
