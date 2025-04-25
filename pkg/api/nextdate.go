package api

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "20060102"
)

var (
	errInvalidJson   = "invalid JSON format"
	errInvalidDate   = errors.New("invalid date format")
	errInvalidRepeat = errors.New("invalid repeat format")
	errMissingTitle  = "title is missing"
	errMissingID     = "Task ID is missing"
	errEmptyRepeat   = errors.New("empty repeat")
	errFailedGet     = "failed to get Task"
)

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	// Проверка на не пустой repeat
	if repeat == "" {
		return "", errEmptyRepeat
	}

	// Проверка на корректную дату
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("incorrect date: %w", err)
	}

	// Проверка на корректный формат repeat
	splitRepeat := strings.Split(repeat, " ")
	if len(splitRepeat) == 0 {
		return "", errInvalidRepeat
	}

	// Выбор функции по указанному формату
	switch splitRepeat[0] {
	case "d":
		return dailyRepeat(now, date, splitRepeat)
	case "y":
		return yearlyRepeat(now, date)
	case "w":
		return weeklyRepeat(now, date, splitRepeat)
	case "m":
		return monthlyRepeat(now, date, splitRepeat)
	default:
		return "", fmt.Errorf("invalid repeat format: %v", splitRepeat[0])
	}
}

// Функция dailyRepeat добавляет указанное кол-во дней
func dailyRepeat(now, date time.Time, splitRepeat []string) (string, error) {
	if len(splitRepeat) != 2 {
		return "", errInvalidRepeat
	}

	// Проверки на правильный формат дней
	days, err := strconv.Atoi(splitRepeat[1])
	if err != nil {
		return "", fmt.Errorf("invalid days format: %v", err)
	}
	if days > 400 {
		return "", fmt.Errorf("the maximum allowed interval (%v days) has been exceeded", 400)
	}
	if days <= 0 {
		return "", fmt.Errorf("days value must be positive: %v", days)
	}

	for {
		date = date.AddDate(0, 0, days)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format(dateFormat), nil
}

// Функция yearlyRepeat добавляет 1 год к дате
func yearlyRepeat(now, date time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format(dateFormat), nil
}

// Сложные функции...
// Функция weeklyRepeat переносит задачи на указанные дни недели
func weeklyRepeat(now, date time.Time, splitRepeat []string) (string, error) {
	if len(splitRepeat) != 2 {
		return "", errInvalidRepeat
	}

	weekDays := strings.Split(splitRepeat[1], ",")
	//Создаем map для хранения дней
	days := make(map[int]bool, len(weekDays))

	for _, v := range weekDays {
		v, err := strconv.Atoi(v)
		if err != nil || v < 1 || v > 7 {
			return "", fmt.Errorf("invalid weekday: %v", v)
		}
		days[v] = true
	}

	// Начинаем со следующего дня
	date = date.AddDate(0, 0, 1)

	for {
		weekDay := int(date.Weekday())
		if weekDay == 0 {
			weekDay = 7
		}

		if days[weekDay] && afterNow(date, now) {
			break
		}

		date = date.AddDate(0, 0, 1)
	}

	return date.Format(dateFormat), nil
}

// Функция monthlyRepeat переносит задачи на указанные дни месяца
func monthlyRepeat(now, date time.Time, splitRepeat []string) (string, error) {
	// Проверка качества (:
	if len(splitRepeat) < 2 {
		return "", errInvalidRepeat
	}

	// Добавляем дни в slice
	daySpl := strings.Split(splitRepeat[1], ",")
	days := make([]int, 0, len(daySpl))
	for _, v := range daySpl {
		day, err := strconv.Atoi(v)
		if err != nil || day < -2 || day == 0 || day > 31 {
			return "", fmt.Errorf("invalid day format: %v", day)
		}
		days = append(days, day)
	}

	// Парсим месяцы, если указаны
	months := make(map[int]bool)
	if len(splitRepeat) > 2 {
		monthParts := strings.Split(splitRepeat[2], ",")
		for _, m := range monthParts {
			month, err := strconv.Atoi(m)
			if err != nil || month < 1 || month > 12 {
				return "", fmt.Errorf("invalid month format: %v", month)
			}
			months[month] = true
		}
	}

	// Начинаем со следующего дня
	date = date.AddDate(0, 0, 1)

	for {
		currentMonth := int(date.Month())
		currentYear := date.Year()

		// Проверяем, что указанного месяца нет в списке
		if len(months) > 0 && !months[currentMonth] {
			date = time.Date(currentYear, time.Month(currentMonth)+1, 1, 0, 0, 0, 0, time.UTC)
			continue
		}

		// Получаем все валидные дни в этом месяце
		validDays := getValidDaysForMonth(date, days)
		if len(validDays) == 0 {
			date = time.Date(currentYear, time.Month(currentMonth)+1, 1, 0, 0, 0, 0, time.UTC)
			continue
		}

		// Ищем следующий день в этом месяце
		found := false
		for _, day := range validDays {
			if day >= date.Day() {
				date = time.Date(currentYear, time.Month(currentMonth), day, 0, 0, 0, 0, time.UTC)
				found = true
				break
			}
		}

		if !found {
			// Перенос на 1-ый валидный день следующего месяца
			date = time.Date(currentYear, time.Month(currentMonth)+1, validDays[0], 0, 0, 0, 0, time.UTC)
		}

		if afterNow(date, now) {
			break
		}

		date = date.AddDate(0, 0, 1)
	}

	return date.Format(dateFormat), nil
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

// Функция getValidDaysForMonth получает валидные дни
func getValidDaysForMonth(date time.Time, ruleDays []int) []int {
	currentMonth := int(date.Month())
	currentYear := date.Year()
	lastDay := time.Date(currentYear, time.Month(currentMonth)+1, 0, 0, 0, 0, 0, time.UTC).Day()

	var validDays []int

	for _, d := range ruleDays {
		switch {
		case d > 0:
			if d <= lastDay {
				validDays = append(validDays, d)
			}
		case d == -1:
			validDays = append(validDays, lastDay)
		case d == -2:
			if lastDay > 1 {
				validDays = append(validDays, lastDay-1)
			}
		}
	}

	sort.Ints(validDays)
	return validDays
}
