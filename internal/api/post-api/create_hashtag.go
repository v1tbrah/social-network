package papi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"gitlab.com/pet-pr-social-network/api-gateway/internal/send"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateHashtagReq struct {
	Name string `json:"name"`
}

type CreateHashtagResp struct {
	ID int64 `json:"id"`
}

// CreateHashtag creates hashtag.
//
//	@Summary		Creates hashtag.
//	@Description	Creates hashtag.
//	@Tags			hashtag
//	@Produce		json
//
//	@Param			objectBody	body		CreateHashtagReq	true	"Hashtag body"
//
//	@Success		200			{object}	CreateHashtagResp
//	@Failure		400			{object}	send.Error
//	@Failure		404			{object}	send.Error
//	@Failure		500			{object}	send.Error
//	@Router			/post/hashtag [post]
func (a *PostAPI) CreateHashtag(w http.ResponseWriter, r *http.Request) {
	var createHashtagReq CreateHashtagReq
	if err := json.NewDecoder(r.Body).Decode(&createHashtagReq); err != nil {
		send.Send(w, send.NewErr(fmt.Sprintf("deserialize req: %v", err)), http.StatusBadRequest)
		return
	}

	createHashtagReq.Name = strings.TrimSpace(createHashtagReq.Name)
	if createHashtagReq.Name == "" {
		send.Send(w, send.NewErr(errEmptyName.Error()), http.StatusBadRequest)
		return
	}

	pbCreateHashtagResp, err := a.postServiceClient.CreateHashtag(r.Context(), &ppbapi.CreateHashtagRequest{Name: createHashtagReq.Name})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			send.Send(w, send.NewErr(err.Error()), http.StatusBadRequest)
			return
		}
		if status.Code(err) == codes.AlreadyExists {
			send.Send(w, send.NewErr(err.Error()), http.StatusConflict)
			return
		}
		log.Error().Err(err).Msg("postServiceClient.CreateHashtag")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	send.Send(w, CreateHashtagResp{ID: pbCreateHashtagResp.GetId()}, http.StatusOK)
}
