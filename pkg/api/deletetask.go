package api

import (
	"net/http"

	"go1f/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		sendErrorResponse(w, errMissingID, http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJson(w, map[string]interface{}{}, http.StatusOK)
}
