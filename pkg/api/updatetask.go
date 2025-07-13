package api

import (
	"encoding/json"
	"net/http"

	"go1f/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendErrorResponse(w, errInvalidJson, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		sendErrorResponse(w, errMissingID, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		sendErrorResponse(w, errMissingTitle, http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJson(w, map[string]interface{}{}, http.StatusOK)
}
