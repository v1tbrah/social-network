package api

import (
	"io"
	"net/http"
)

func (a *API) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
}
