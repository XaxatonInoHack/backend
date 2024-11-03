package parser

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseReview(text string) (map[string]float64, map[string]string) {
	const op = "parser.ParseReview"

	sections := strings.Split(text, "**")

	scores := make(map[string]float64)
	texts := make(map[string]string)

	for i := 1; i < len(sections); i += 2 {
		section := strings.TrimSpace(sections[i])

		colonIndex := strings.Index(section, ":")
		if colonIndex == -1 {
			continue
		}
		dotIndex := strings.Index(section, ".")
		category := strings.TrimSpace(section[dotIndex+1 : colonIndex])
		category = strings.TrimSpace(category)

		scoreText := strings.TrimSpace(section[colonIndex+1:])

		// Проверяем, есть ли "/"
		var score string
		if strings.Contains(scoreText, "/") {
			scoreEnd := strings.Index(scoreText, "/")
			score = scoreText[:scoreEnd] // Получаем только часть до "/"
		} else {
			score = scoreText // Если нет "/", берем всю строку
		}
		// Преобразуем оценку в float64
		score = strings.TrimSpace(score) // Убираем лишние пробелы
		scoreFloat, err := strconv.ParseFloat(score, 64)
		if err != nil {
			fmt.Printf("Ошибка при преобразовании оценки для категории %s: %v\n", category, err)
			continue
		}

		// Сохраняем категорию и её оценку
		scores[category] = scoreFloat

		if i+1 < len(sections) {
			texts[category] = strings.TrimSpace(sections[i+1])
		}
	}

	return scores, texts
}
