package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ParseScoreOnly(text string) map[string]float64 {
	// Регулярное выражение для поиска строк с оценками
	re := regexp.MustCompile(`(?m)^\d+\.\s*(\w+):\s*([\d\.]+)`)

	// Находим все совпадения в тексте
	matches := re.FindAllStringSubmatch(text, -1)

	// Инициализируем карту для хранения результатов
	result := make(map[string]float64)

	// Проходим по всем найденным совпадениям и заполняем карту
	for _, match := range matches {
		category := match[1]
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

func ParseCategoryTexts(text string) map[string]string {
	// Регулярное выражение для поиска заголовков категорий
	re := regexp.MustCompile(`(?m)^(\d+)\.\s*([A-Za-zА-Яа-яёЁ ]+):\s*\d+\.?\d*`)

	// Находим все заголовки категорий с их позициями в тексте
	matches := re.FindAllStringSubmatchIndex(text, -1)

	result := make(map[string]string)

	// Поиск позиции раздела "Overall resume"
	overallIndex := strings.Index(strings.ToLower(text), "overall resume:")
	var overallContent string
	if overallIndex != -1 {
		overallContent = strings.TrimSpace(text[overallIndex+len("Overall resume:"):])
		// Обрезаем текст до начала "Overall resume"
		text = text[:overallIndex]
	}

	for i, match := range matches {
		// Получаем название категории
		categoryStart, categoryEnd := match[4], match[5]
		category := strings.TrimSpace(text[categoryStart:categoryEnd])

		// Определяем начало контента после заголовка
		contentStart := match[1]
		var contentEnd int
		if i+1 < len(matches) {
			// Если есть следующая категория, контент заканчивается перед ней
			contentEnd = matches[i+1][0]
		} else {
			// Если это последняя категория, контент идет до конца текста
			contentEnd = len(text)
		}

		// Извлекаем контент между заголовками
		content := text[contentStart:contentEnd]

		// Удаляем первую строку (заголовок) из контента
		lines := strings.SplitN(content, "\n", 2)
		if len(lines) >= 2 {
			content = strings.TrimSpace(lines[1])
		} else {
			content = ""
		}

		result[category] = content
	}

	if overallContent != "" {
		result["Overall resume"] = overallContent
	}

	return result
}
