package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ParseScoreOnly(text string) map[string]float64 {
	// Обновленное регулярное выражение
	re := regexp.MustCompile(`(?m)^\s*\*{0,2}\s*\d+\.\s*([^\n:]+?):\s*([\d\.]+)\s*\*{0,2}`)

	// Находим все совпадения в тексте
	text = strings.ReplaceAll(text, `\n`, "\n")
	matches := re.FindAllStringSubmatch(text, -1)

	// Инициализируем карту для хранения результатов
	result := make(map[string]float64)

	// Проходим по всем найденным совпадениям и заполняем карту
	for _, match := range matches {
		category := strings.TrimSpace(match[1])
		scoreStr := match[2]
		score, err := strconv.ParseFloat(scoreStr, 64)
		if err != nil {
			fmt.Println("Ошибка при парсинге оценки:", err)
			continue
		}
		result[category] = score
	}

	return result
}

func ParseCriteriaText(text string) map[string]string {
	// Регулярное выражение для поиска заголовков критериев
	re := regexp.MustCompile(`(?m)^(?:\d+\.\s*)?([A-Za-zА-Яа-яёЁ ]+):\s*(\d+)?\s*$`)

	// Находим все совпадения заголовков и их позиции в тексте
	text = strings.ReplaceAll(text, `\n`, "\n")
	matches := re.FindAllStringSubmatchIndex(text, -1)

	result := make(map[string]string)

	for i, match := range matches {
		// Получаем название критерия
		criterionStart, criterionEnd := match[2], match[3]
		criterion := strings.TrimSpace(text[criterionStart:criterionEnd])

		// Определяем начало и конец содержимого
		contentStart := match[1]
		var contentEnd int
		if i+1 < len(matches) {
			contentEnd = matches[i+1][0]
		} else {
			contentEnd = len(text)
		}

		// Извлекаем содержимое
		content := text[contentStart:contentEnd]

		// Удаляем строку заголовка из содержимого
		lines := strings.SplitN(content, "\n", 2)
		if len(lines) >= 2 {
			content = strings.TrimSpace(lines[1])
		} else {
			content = ""
		}

		result[criterion] = content
	}

	return result
}

// ParseScores Функция для парсинга строки в map[string]float64
func ParseScores(text string) map[string]float64 {
	result := make(map[string]float64)

	// Регулярное выражение для поиска пар "Критерий: Оценка"
	re := regexp.MustCompile(`([A-Za-zА-Яа-яёЁ ]+):\s*([\d\.]+)`)

	// Находим все совпадения в тексте
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		criterion := strings.TrimSpace(match[1])
		scoreStr := match[2]
		score, err := strconv.ParseFloat(scoreStr, 64)
		if err != nil {
			fmt.Println("Ошибка при парсинге оценки:", err)
			continue
		}
		result[criterion] = score
	}

	return result
}
