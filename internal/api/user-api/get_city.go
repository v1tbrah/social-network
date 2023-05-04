package uapi

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"gitlab.com/pet-pr-social-network/api-gateway/internal/send"
	"gitlab.com/pet-pr-social-network/user-service/pbapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetCityReq struct {
	ID string `json:"id" example:"1"`
}

type GetCityResp struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"Moscow"`
}

// GetCity returns city.
//
//	@Summary		Returns city.
//	@Description	Returns city by id.
//	@Tags			city
//	@Produce		json
//	@Param			id	path		int	true	"City id"
//	@Success		200	{object}	GetCityResp
//	@Failure		400	{object}	send.Error
//	@Failure		404	{object}	send.Error
//	@Failure		500	{object}	send.Error
//	@Router			/user/city/{id} [get]
func (a *UserAPI) GetCity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		send.Send(w, send.NewErr(errInvalidID.Error()), http.StatusBadRequest)
		return
	}

	pbGetCityResp, err := a.userServiceClient.GetCity(r.Context(), &pbapi.GetCityRequest{Id: id})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			send.Send(w, send.NewErr(err.Error()), http.StatusNotFound)
			return
		}
		log.Error().Err(err).Msg("userServiceClient.GetCity")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	var resp GetCityResp
	if pbGetCityResp.GetCity() != nil {
		resp.ID = pbGetCityResp.GetCity().GetId()
		resp.Name = pbGetCityResp.GetCity().GetName()
	}
	send.Send(w, resp, http.StatusOK)
}
