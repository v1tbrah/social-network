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

type CreateUserReq struct {
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	InterestsID []int64 `json:"interests_id"`
	CityID      int64   `json:"city_id"`
}

type CreateUserResp struct {
	ID int64 `json:"id"`
}

// CreateUser creates user.
//
//	@Summary		Creates user.
//	@Description	Creates user.
//	@Tags			user
//	@Produce		json
//
//	@Param			objectBody	body		CreateUserReq	true	"User body"
//
//	@Success		200			{object}	CreateUserResp
//	@Failure		400			{object}	send.Error
//	@Failure		404			{object}	send.Error
//	@Failure		500			{object}	send.Error
//	@Router			/user/user [post]
func (a *UserAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	var createUserReq CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&createUserReq); err != nil {
		send.Send(w, send.NewErr(fmt.Sprintf("deserialize req: %v", err)), http.StatusBadRequest)
		return
	}

	createUserReq.Name = strings.TrimSpace(createUserReq.Name)
	if createUserReq.Name == "" {
		send.Send(w, send.NewErr(errEmptyName.Error()), http.StatusBadRequest)
		return
	}

	createUserReq.Surname = strings.TrimSpace(createUserReq.Surname)
	if createUserReq.Surname == "" {
		send.Send(w, send.NewErr(errEmptySurname.Error()), http.StatusBadRequest)
		return
	}

	pbCreateUserResp, err := a.userServiceClient.CreateUser(r.Context(), createUserReqToPBCreateUserReq(createUserReq))
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			send.Send(w, send.NewErr(err.Error()), http.StatusBadRequest)
			return
		}
		if status.Code(err) == codes.AlreadyExists {
			send.Send(w, send.NewErr(err.Error()), http.StatusConflict)
			return
		}
		log.Error().Err(err).Msg("userServiceClient.CreateUser")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	send.Send(w, CreateUserResp{ID: pbCreateUserResp.GetId()}, http.StatusOK)
}

func createUserReqToPBCreateUserReq(createUserReq CreateUserReq) *pbapi.CreateUserRequest {
	return &pbapi.CreateUserRequest{
		Name:        createUserReq.Name,
		Surname:     createUserReq.Surname,
		InterestsID: createUserReq.InterestsID,
		CityID:      createUserReq.CityID,
	}
}
