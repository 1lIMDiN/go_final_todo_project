package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT,
    repeat VARCHAR(128) DEFAULT ''
);

CREATE INDEX idx_date ON scheduler(date);
`

func Init(dbFile string) error {
	//Проверяем существование файла БД
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	//Открываем БД
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	//Если файл не существовал, выполняем SQL команды из переменной schema.
	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return fmt.Errorf("failed to create schema: %v", err)
		}
	}

	//Используем глобальную переменнную для хранения идентификатора открытой БД
	DB = db
	return nil
}
