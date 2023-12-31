package uapi

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/api-gateway/internal/send"

	"github.com/v1tbrah/user-service/upbapi"
)

// GetAllCities returns all cities.
//
//	@Summary		Returns all cities.
//	@Description	Returns all cities.
//	@Tags			city
//	@Produce		json
//	@Success		200	{object}	[]GetCityResp
//	@Failure		500	{object}	send.Error
//	@Router			/user/city [get]
func (a *UserAPI) GetAllCities(w http.ResponseWriter, r *http.Request) {
	pbGetAllCitiesResp, err := a.userServiceClient.GetAllCities(r.Context(), &upbapi.Empty{})
	if err != nil {
		log.Error().Err(err).Msg("userServiceClient.GetAllCities")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	resp := make([]GetCityResp, 0, len(pbGetAllCitiesResp.GetCities()))
	for _, pbCity := range pbGetAllCitiesResp.GetCities() {
		resp = append(resp, GetCityResp{
			ID:   pbCity.GetId(),
			Name: pbCity.GetName(),
		})
	}

	send.Send(w, resp, http.StatusOK)
}
