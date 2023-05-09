package uapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"gitlab.com/pet-pr-social-network/api-gateway/internal/send"
	"gitlab.com/pet-pr-social-network/user-service/pbapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateCityReq struct {
	Name string `json:"name" example:"Moscow"`
}

type CreateCityResp struct {
	ID int64 `json:"id" example:"1"`
}

// CreateCity creates city.
//
//	@Summary		Creates city.
//	@Description	Creates city.
//	@Tags			city
//	@Produce		json
//
//	@Param			objectBody	body		CreateCityReq	true	"City body"
//
//	@Success		200			{object}	CreateCityResp
//	@Failure		400			{object}	send.Error
//	@Failure		404			{object}	send.Error
//	@Failure		500			{object}	send.Error
//	@Router			/user/city [post]
func (a *UserAPI) CreateCity(w http.ResponseWriter, r *http.Request) {
	var createCityReq CreateCityReq
	if err := json.NewDecoder(r.Body).Decode(&createCityReq); err != nil {
		send.Send(w, send.NewErr(fmt.Sprintf("deserialize req: %v", err)), http.StatusBadRequest)
		return
	}

	createCityReq.Name = strings.TrimSpace(createCityReq.Name)
	if createCityReq.Name == "" {
		send.Send(w, send.NewErr(errEmptyName.Error()), http.StatusBadRequest)
		return
	}

	pbCreateCityResp, err := a.userServiceClient.CreateCity(r.Context(), &pbapi.CreateCityRequest{Name: createCityReq.Name})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			send.Send(w, send.NewErr(err.Error()), http.StatusBadRequest)
			return
		}
		if status.Code(err) == codes.AlreadyExists {
			send.Send(w, send.NewErr(err.Error()), http.StatusConflict)
			return
		}
		log.Error().Err(err).Msg("userServiceClient.CreateCity")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	send.Send(w, CreateCityResp{ID: pbCreateCityResp.GetId()}, http.StatusOK)
}
