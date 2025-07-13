package api

import (
	"net/http"

	"go1f/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	tasks, err := db.GetTasks(50, search)
	if err != nil {
		sendErrorResponse(w, errFailedGet, http.StatusBadRequest)
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	}, http.StatusOK)
}
