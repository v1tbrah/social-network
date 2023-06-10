package papi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/v1tbrah/api-gateway/internal/send"

	"github.com/v1tbrah/post-service/ppbapi"
)

type CreatePostReq struct {
	UserID      int64   `json:"user_id" example:"1"`
	Description string  `json:"description" example:"description"`
	HashtagsID  []int64 `json:"hashtags_id" example:"[1,2,3]"`
}

type CreatePostResp struct {
	ID int64 `json:"id" example:"1"`
}

// CreatePost creates post.
//
//	@Summary		Creates post.
//	@Description	Creates post.
//	@Tags			post
//	@Produce		json
//
//	@Param			objectBody	body		CreatePostReq	true	"Post body"
//
//	@Success		200			{object}	CreatePostResp
//	@Failure		400			{object}	send.Error
//	@Failure		404			{object}	send.Error
//	@Failure		500			{object}	send.Error
//	@Router			/post/post [post]
func (a *PostAPI) CreatePost(w http.ResponseWriter, r *http.Request) {
	var createPostReq CreatePostReq
	if err := json.NewDecoder(r.Body).Decode(&createPostReq); err != nil {
		send.Send(w, send.NewErr(fmt.Sprintf("deserialize req: %v", err)), http.StatusBadRequest)
		return
	}

	pbCreatePostResp, err := a.postServiceClient.CreatePost(r.Context(), &ppbapi.CreatePostRequest{
		UserID: createPostReq.UserID, Description: createPostReq.Description, HashtagsID: createPostReq.HashtagsID})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			send.Send(w, send.NewErr(err.Error()), http.StatusBadRequest)
			return
		}
		if status.Code(err) == codes.AlreadyExists {
			send.Send(w, send.NewErr(err.Error()), http.StatusConflict)
			return
		}
		log.Error().Err(err).Msg("postServiceClient.CreatePost")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	send.Send(w, CreatePostResp{ID: pbCreatePostResp.GetId()}, http.StatusOK)
}
