package uapi

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"gitlab.com/pet-pr-social-network/api-gateway/internal/send"
	"gitlab.com/pet-pr-social-network/user-service/pbapi"
)

// GetAllInterests returns all interests.
//
//	@Summary		Returns all interests.
//	@Description	Returns all interests.
//	@Tags			interest
//	@Produce		json
//	@Success		200	{object}	[]GetCityResp
//	@Failure		500	{object}	send.Error
//	@Router			/user/interest [get]
func (a *UserAPI) GetAllInterests(w http.ResponseWriter, r *http.Request) {
	pbGetAllInterestsResp, err := a.userServiceClient.GetAllInterests(r.Context(), &pbapi.Empty{})
	if err != nil {
		log.Error().Err(err).Msg("userServiceClient.GetAllInterests")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	resp := make([]GetInterestResp, 0, len(pbGetAllInterestsResp.GetInterests()))
	for _, pbCity := range pbGetAllInterestsResp.GetInterests() {
		resp = append(resp, GetInterestResp{
			ID:   pbCity.GetId(),
			Name: pbCity.GetName(),
		})
	}

	send.Send(w, resp, http.StatusOK)
}
