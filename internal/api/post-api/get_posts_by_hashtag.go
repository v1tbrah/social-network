package papi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/v1tbrah/api-gateway/internal/send"

	"github.com/v1tbrah/post-service/ppbapi"
)

type GetPostsByHashtagReq struct {
	HashtagID    int64     `json:"hashtagID" example:"1"`
	Direction    Direction `json:"direction" example:"0"`
	PostOffsetID int64     `json:"post_offset_id" example:"0"`
	Limit        int64     `json:"limit" example:"10"`
}

type Direction int32

const (
	First Direction = iota
	Next
	Prev
)

type Post struct {
	Id          int64     `json:"id" example:"1"`
	UserID      int64     `json:"user_id" example:"1"`
	Description string    `json:"description" example:"description"`
	HashtagsID  []int64   `json:"hashtags_id" example:"1,2,3"`
	CreatedAt   time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
}

// GetPostsByHashtag returns posts by hashtag.
//
//	@Summary		Returns posts by hashtag.
//	@Description	Returns posts by hashtag.
//	@Tags			post
//	@Produce		json
//
//	@Param			objectBody	body		GetPostsByHashtagReq	true	"GetPostsByHashtagReq body. Direction: 0 - first, 1 - next, 2 - prev"
//
//	@Success		200			{object}	[]Post
//	@Success		204
//	@Failure		400	{object}	send.Error
//	@Failure		500	{object}	send.Error
//	@Router			/post/post/get_by_hashtag [post]
func (a *PostAPI) GetPostsByHashtag(w http.ResponseWriter, r *http.Request) {
	var getPostsReq GetPostsByHashtagReq
	if err := json.NewDecoder(r.Body).Decode(&getPostsReq); err != nil {
		send.Send(w, send.NewErr(fmt.Sprintf("deserialize req: %v", err)), http.StatusBadRequest)
		return
	}

	pbGetPostsResp, err := a.postServiceClient.GetPostsByHashtag(r.Context(), &ppbapi.GetPostsByHashtagRequest{
		HashtagID: getPostsReq.HashtagID, Direction: ppbapi.GetPostsByHashtagRequest_DIRECTION(getPostsReq.Direction),
		PostOffsetID: getPostsReq.PostOffsetID, Limit: getPostsReq.Limit})
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

	if len(pbGetPostsResp.GetPosts()) == 0 {
		send.Send(w, nil, http.StatusNoContent)
		return
	}

	resp := make([]Post, 0, len(pbGetPostsResp.GetPosts()))
	for _, post := range pbGetPostsResp.GetPosts() {
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
