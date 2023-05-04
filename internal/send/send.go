package send

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func Send(w http.ResponseWriter, resp any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if resp == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Error().Err(err).Msg("serialize resp")
		http.Error(w, fmt.Sprintf("serialize resp: %v", err), http.StatusInternalServerError)
	}
}
