package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	fields := strings.Split(data, ",")

	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("некорректные данные")
	}

	steps, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования строки в число: %w", err)
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше нуля")
	}

	interval, err := time.ParseDuration(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования продолжительности: %w", err)
	}

	if interval <= 0 {
		return 0, 0, fmt.Errorf("продолжительность прогулки не может быть нолевой или отрицательной")
	}

	return steps, interval, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, interval, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distanceM := float64(steps) * stepLength
	distanceKm := distanceM / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, interval)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories)
}
