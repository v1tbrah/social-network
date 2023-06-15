package mapi

import (
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/v1tbrah/api-gateway/internal/send"
	"github.com/v1tbrah/media-service/mpbapi"
)

type SavePostResponse struct {
	GUID string `json:"guid"`
}

// SavePost saves post media content.
//
//	@Summary		Saves post media content.
//	@Description	Saves post media content and returns guid.
//	@Tags			media
//	@Produce		mpfd
//
//	@Param			file	formData	file	true	"Body with file"
//
//	@Success		200		{object}	SavePostResponse
//	@Failure		400		{object}	send.Error
//	@Failure		500		{object}	send.Error
//	@Router			/media/post [post]
func (a *MediaAPI) SavePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil { // max size 32 mb
		send.Send(w, send.NewErr(err.Error()), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		send.Send(w, send.NewErr(err.Error()), http.StatusBadRequest)
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Error().Err(err).Msg("file.Close")
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		log.Error().Err(err).Msg("io.ReadAll")
		return
	}

	pbAddPostResp, err := a.mediaServiceClient.AddPost(r.Context(), &mpbapi.AddPostRequest{Data: data})
	if err != nil {
		log.Error().Err(err).Msg("mediaServiceClient.AddPost")
		send.Send(w, send.NewErr(err.Error()), http.StatusInternalServerError)
		return
	}

	send.Send(w, SavePostResponse{GUID: pbAddPostResp.GetGuid()}, http.StatusOK)
}
