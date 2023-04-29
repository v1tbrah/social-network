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

type CreateInterestReq struct {
	Name string `json:"name"`
}

type CreateInterestResp struct {
	ID int64 `json:"id"`
}

// CreateInterest creates interest.
//
//	@Summary		Creates interest.
//	@Description	Creates interest.
//	@Tags			interest
//	@Produce		json
//
//	@Param			objectBody	body		CreateInterestReq	true	"Interest body"
//
//	@Success		200			{object}	CreateInterestResp
//	@Failure		400			{object}	send.Error
//	@Failure		404			{object}	send.Error
//	@Failure		500			{object}	send.Error
//	@Router			/user/interest [post]
func (a *UserAPI) CreateInterest(w http.ResponseWriter, r *http.Request) {
	var createInterestReq CreateInterestReq
	if err := json.NewDecoder(r.Body).Decode(&createInterestReq); err != nil {
		send.Send(w, send.NewErr(fmt.Sprintf("deserialize req: %v", err)), http.StatusBadRequest)
		return
	}

	createInterestReq.Name = strings.TrimSpace(createInterestReq.Name)
	if createInterestReq.Name == "" {
		send.Send(w, send.NewErr(errEmptyName.Error()), http.StatusBadRequest)
		return
	}

	pbCreateInterestResp, err := a.userServiceClient.CreateInterest(r.Context(), &pbapi.CreateInterestRequest{Name: createInterestReq.Name})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			send.Send(w, send.NewErr(err.Error()), http.StatusBadRequest)
			return
		}
		if status.Code(err) == codes.AlreadyExists {
			send.Send(w, send.NewErr(err.Error()), http.StatusConflict)
			return
		}
		log.Error().Err(err).Msg("userServiceClient.CreateInterest")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	send.Send(w, CreateInterestResp{ID: pbCreateInterestResp.GetId()}, http.StatusOK)
}
