package mapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/v1tbrah/media-service/mpbapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/v1tbrah/api-gateway/internal/send"
)

// GetPost returns post media content.
//
//	@Summary		Returns post by guid.
//	@Description	Returns post by guid.
//	@Tags			media
//	@Produce		json
//
//	@Param			guid	path		string	true	"Post guid"
//
//	@Success		200		{object}	[]byte
//	@Failure		400		{object}	send.Error
//	@Failure		404		{object}	send.Error
//	@Failure		500		{object}	send.Error
//	@Router			/media/post/{guid} [get]
func (a *MediaAPI) GetPost(w http.ResponseWriter, r *http.Request) {
	guid := chi.URLParam(r, "guid")
	if guid == "" {
		send.Send(w, send.NewErr(errEmptyGUID.Error()), http.StatusBadRequest)
		return
	}

	pbGetPostResp, err := a.mediaServiceClient.GetPost(r.Context(), &mpbapi.GetPostRequest{Guid: guid})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			send.Send(w, send.NewErr(err.Error()), http.StatusNotFound)
			return
		}

		log.Error().Err(err).Msg("mediaServiceClient.GetPost")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	send.Send(w, pbGetPostResp.GetData(), http.StatusOK)
}
