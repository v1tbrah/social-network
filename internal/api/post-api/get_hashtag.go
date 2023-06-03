package papi

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/pet-pr-social-network/api-gateway/internal/send"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
)

type GetHashtagResp struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"#cat"`
}

// GetHashtag returns hashtag.
//
//	@Summary		Returns hashtag.
//	@Description	Returns hashtag by id.
//	@Tags			city
//	@Produce		json
//	@Param			id	path		int	true	"Hashtag id"
//	@Success		200	{object}	GetHashtagResp
//	@Failure		400	{object}	send.Error
//	@Failure		404	{object}	send.Error
//	@Failure		500	{object}	send.Error
//	@Router			/post/hashtag/{id} [get]
func (a *PostAPI) GetHashtag(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		send.Send(w, send.NewErr(errInvalidID.Error()), http.StatusBadRequest)
		return
	}

	pbGetHashtagResp, err := a.postServiceClient.GetHashtag(r.Context(), &ppbapi.GetHashtagRequest{Id: id})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			send.Send(w, send.NewErr(err.Error()), http.StatusNotFound)
			return
		}
		log.Error().Err(err).Msg("postServiceClient.GetHashtag")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	var resp GetHashtagResp
	if pbGetHashtagResp.GetHashtag() != nil {
		resp.ID = pbGetHashtagResp.GetHashtag().GetId()
		resp.Name = pbGetHashtagResp.GetHashtag().GetName()
	}
	send.Send(w, resp, http.StatusOK)
}
