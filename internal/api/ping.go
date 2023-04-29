package api

import (
	"io"
	"net/http"
)

func (a *API) ping(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, http.StatusText(http.StatusOK))
}
