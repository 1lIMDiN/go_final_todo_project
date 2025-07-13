package server

import (
	"fmt"
	"net/http"

	"go1f/pkg/api"
)

func Run(port int) error {
	// Обработчик API
	api.Init()

	// Настраиваем обработчик для статических файлов
	http.Handle("/", http.FileServer(http.Dir("web")))

	// Запускаем сервер
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
