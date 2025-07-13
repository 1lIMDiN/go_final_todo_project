package api

import (
	"net/http"
	"time"

	"go1f/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	now := time.Now()

	if task.Repeat == "" {
		// Удаляем задачу, если repeat пуст
		if err := db.DeleteTask(id); err != nil {
			sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Если задача периодическая, то необходимо получить следующую дату
		nextdate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			sendErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := db.UpdateDate(id, nextdate); err != nil {
			sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	writeJson(w, map[string]interface{}{}, http.StatusOK)
}
