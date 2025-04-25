package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go1f/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Декодируем полученный json в task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendErrorResponse(w, errInvalidJson, http.StatusBadRequest)
		return
	}

	// Проверка на не пустой Title
	if task.Title == "" {
		sendErrorResponse(w, errMissingTitle, http.StatusBadRequest)
		return
	}

	// Проверяем на корректность Date
	if err := checkDate(&task); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Добавляем task в БД
	id, err := db.AddTask(&task)
	if err != nil {
		log.Printf("Failed to add task: %v, task: %+v", err, task)
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]interface{}{"id": strconv.Itoa(id)}, http.StatusOK)
}

// Функция checkDate проверяет Date на корректность
func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format(dateFormat) // dateFormat указан в nextdate.go

	// Если Date пустой - присваиваем ему текущее время
	if task.Date == "" {
		task.Date = today
		return nil
	}

	// Проверяем формат даты (YYYYMMDD)
	if len(task.Date) != 8 {
		return errInvalidDate
	}

	// Парсим дату
	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return errInvalidDate
	}

	// Проверяем, что дата существует (не 30 февраля и т.д)
	if t.Format(dateFormat) != task.Date {
		return errInvalidDate
	}

	// Проверяем repeat, только если он не пустой
	if task.Repeat != "" {
		parts := strings.Split(task.Repeat, " ")
		if len(parts) == 0 {
			return errInvalidRepeat
		}
		// Проверка repeat форматов
		switch parts[0] {
		case "d", "y", "w", "m":
			break
		default:
			return errInvalidRepeat
		}

	}

	chNow := now.Truncate(24 * time.Hour)
	chDate := t.Truncate(24 * time.Hour)

	// Обрабатываем даты в прошлом
	if chDate.Before(chNow) {
		if task.Repeat == "" {
			task.Date = today
		} else {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = nextDate
		}
	}

	return nil
}

// Функция writeJson сериализует Data в JSON
func writeJson(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Печатаем ошибки
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
