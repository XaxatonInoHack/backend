package converter

import "strconv"

func ScoreToResult(score map[string]float64) string {
	result := 0.0
	for _, value := range score {
		result += value
	}

	result = result / float64(len(score)) / float64(5)
	return strconv.FormatFloat(result, 'f', -1, 64)
}
