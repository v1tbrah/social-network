package uapi

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/v1tbrah/api-gateway/internal/send"

	"github.com/v1tbrah/user-service/upbapi"
)

type GetInterestReq struct {
	ID string `json:"id" example:"1"`
}

type GetInterestResp struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"Music"`
}

// GetInterest returns interest.
//
//	@Summary		Returns interest.
//	@Description	Returns interest by id.
//	@Tags			interest
//	@Produce		json
//	@Param			id	path		int	true	"Interest id"
//	@Success		200	{object}	GetInterestReq
//	@Failure		400	{object}	send.Error
//	@Failure		404	{object}	send.Error
//	@Failure		500	{object}	send.Error
//	@Router			/user/interest/{id} [get]
func (a *UserAPI) GetInterest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		send.Send(w, send.NewErr(errInvalidID.Error()), http.StatusBadRequest)
		return
	}

	pbGetInterestResp, err := a.userServiceClient.GetInterest(r.Context(), &upbapi.GetInterestRequest{Id: id})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			send.Send(w, send.NewErr(err.Error()), http.StatusNotFound)
			return
		}
		log.Error().Err(err).Msg("userServiceClient.GetInterest")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	var resp GetInterestResp
	if pbGetInterestResp.GetInterest() != nil {
		resp.ID = pbGetInterestResp.GetInterest().GetId()
		resp.Name = pbGetInterestResp.GetInterest().GetName()
	}
	send.Send(w, resp, http.StatusOK)
}
