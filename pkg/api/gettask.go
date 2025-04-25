package api

import (
	"net/http"

	"go1f/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		sendErrorResponse(w, errMissingID, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJson(w, task, http.StatusOK)
}
