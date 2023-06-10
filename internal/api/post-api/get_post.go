package papi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/v1tbrah/api-gateway/internal/send"

	"github.com/v1tbrah/post-service/ppbapi"
)

type GetPostResp struct {
	UserID      int64     `json:"user_id" example:"1"`
	Description string    `json:"description" example:"My cat"`
	HashtagsID  []int64   `json:"hashtags_id" example:"1,2,3"`
	CreatedAt   time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
}

// GetPost returns post.
//
//	@Summary		Returns post.
//	@Description	Returns post by id.
//	@Tags			post
//	@Produce		json
//	@Param			id	path		int	true	"Post id"
//	@Success		200	{object}	GetPostResp
//	@Failure		400	{object}	send.Error
//	@Failure		404	{object}	send.Error
//	@Failure		500	{object}	send.Error
//	@Router			/post/post/{id} [get]
func (a *PostAPI) GetPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		send.Send(w, send.NewErr(errInvalidID.Error()), http.StatusBadRequest)
		return
	}

	pbGetPostResp, err := a.postServiceClient.GetPost(r.Context(), &ppbapi.GetPostRequest{Id: id})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			send.Send(w, send.NewErr(err.Error()), http.StatusNotFound)
			return
		}
		log.Error().Err(err).Msg("postServiceClient.GetPost")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	var resp GetPostResp
	if pbGetPostResp != nil {
		resp.UserID = pbGetPostResp.GetUserID()
		resp.Description = pbGetPostResp.GetDescription()
		resp.HashtagsID = pbGetPostResp.GetHashtagsID()
		if pbGetPostResp.GetCreatedAt() != nil {
			resp.CreatedAt = pbGetPostResp.GetCreatedAt().AsTime()
		}
	}
	send.Send(w, resp, http.StatusOK)
}
