package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // Средняя длина шага.
	mInKm                      = 1000 // Количество метров в километре.
	minInH                     = 60   // Количество минут в часе.
	stepLengthCoefficient      = 0.45 // Коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // Коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	fields := strings.Split(data, ",")

	if len(fields) != 3 {
		return 0, "0", 0, fmt.Errorf("некорректные данные")
	}

	steps, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, "0", 0, fmt.Errorf("ошибка преобразования строки в число: %w", err)
	}

	duration, err := time.ParseDuration(fields[2])
	if err != nil {
		return 0, "0", 0, fmt.Errorf("ошибка преобразования продолжительности: %w", err)
	}

	if steps <= 0 || duration <= 0 {
		return 0, "0", 0, fmt.Errorf("количество шагов и продолжительность не могут быть отрицательными")
	}

	return steps, fields[1], duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceKm := (stepLength * float64(steps)) / mInKm

	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distanceKm := distance(steps, height)
	durationInHours := duration.Hours()

	return distanceKm / durationInHours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	durationInHours := duration.Hours()
	avgSpeed := meanSpeed(steps, height, duration)
	distanceKm := distance(steps, height)
	var calories float64

	switch trainingType {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("ошибка расчёта калорий при ходьбе: %w", err)
		}

	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("ошибка расчёта калорий при беге: %w", err)
		}

	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		trainingType, durationInHours, distanceKm, avgSpeed, calories,
	), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("один или несколько параметров равны или меньше нуля")
	}

	avgSpeed := meanSpeed(steps, height, duration)
	if avgSpeed <= 0 {
		return 0, fmt.Errorf("скорость не может быть меньше или равна нулю")
	}

	durationInMinutes := duration.Minutes()
	return weight * avgSpeed * durationInMinutes / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("один или несколько параметров равны или меньше нуля")
	}

	avgSpeed := meanSpeed(steps, height, duration)
	if avgSpeed <= 0 {
		return 0, fmt.Errorf("средняя скорость не может быть меньше или равна нулю")
	}

	durationInMinutes := duration.Minutes()
	calories := weight * avgSpeed * durationInMinutes / minInH

	return calories * walkingCaloriesCoefficient, nil
}
