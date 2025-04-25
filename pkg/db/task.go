package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Функция AddTask добавляет данные в БД и возвращает id
func AddTask(task *Task) (int, error) {
	res, err := DB.Exec(
		"INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Функция GetTasks получает данные с БД
func GetTasks(limit int, search string) ([]*Task, error) {
	var rows *sql.Rows
	var err error

	if search == "" {
		query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		rows, err = DB.Query(query, limit)
	} else {
		// Проверяем, является ли search датой формата 02.01.2006
		if date, err := time.Parse("02.01.2006", search); err == nil {
			formattedDate := date.Format("20060102")
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
			rows, err = DB.Query(query, formattedDate, limit)
			if err != nil {
				return nil, fmt.Errorf("failed to query tasks: %v", err)
			}
		} else {
			// Поиск по подстроке
			searchTerm := "%" + search + "%"
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
			rows, err = DB.Query(query, searchTerm, searchTerm, limit)
			if err != nil {
				return nil, fmt.Errorf("failed to query tasks: %v", err)
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %v", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("failed to skan task: %v", err)
		}
		tasks = append(tasks, &task)
	}

	if tasks == nil {
		tasks = make([]*Task, 0)
	}

	return tasks, nil
}

// Функция GetTask получает task по полученному id
func GetTask(id string) (*Task, error) {
	var task Task

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %v", err)
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check updated rows: %v", err)
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

// Функция DeleteTask удаляет задачу по полученному ID
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check: %v", err)
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

// Функция UpdateDate обновляет дату
func UpdateDate(id, date string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, date, id)
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check updated rows: %v", err)
	}
	if count == 0 {
		return fmt.Errorf(`incorrect ID for updating task`)
	}
	return nil
}
