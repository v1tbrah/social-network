package uapi

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/pet-pr-social-network/api-gateway/internal/send"
	"gitlab.com/pet-pr-social-network/user-service/upbapi"
)

type GetUserReq struct {
	ID string `json:"id" example:"1"`
}

type GetUserResp struct {
	Name        string  `json:"name" example:"John"`
	Surname     string  `json:"surname" example:"Doe"`
	InterestsID []int64 `json:"interests_id" example:"1,2,3"`
	CityID      int64   `json:"city_id" example:"1"`
}

// GetUser returns user.
//
//	@Summary		Returns user.
//	@Description	Returns user by id.
//	@Tags			user
//	@Produce		json
//	@Param			id	path		int	true	"User id"
//	@Success		200	{object}	GetUserResp
//	@Failure		400	{object}	send.Error
//	@Failure		404	{object}	send.Error
//	@Failure		500	{object}	send.Error
//	@Router			/user/user/{id} [get]
func (a *UserAPI) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		send.Send(w, send.NewErr(errInvalidID.Error()), http.StatusBadRequest)
		return
	}

	pbGetUserResp, err := a.userServiceClient.GetUser(r.Context(), &upbapi.GetUserRequest{Id: id})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			send.Send(w, send.NewErr(err.Error()), http.StatusNotFound)
			return
		}
		log.Error().Err(err).Msg("userServiceClient.GetUser")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	var resp GetUserResp
	if pbGetUserResp != nil {
		resp.Name = pbGetUserResp.GetName()
		resp.Surname = pbGetUserResp.GetSurname()
		resp.InterestsID = pbGetUserResp.GetInterestsID()
		resp.CityID = pbGetUserResp.GetCityID()
	}
	send.Send(w, resp, http.StatusOK)
}
