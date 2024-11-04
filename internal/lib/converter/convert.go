package converter

import "strconv"

func ScoreToResult(score map[string]float64, prec int) string {
	result := 0.0
	for _, value := range score {
		result += value
	}

	result = result / float64(len(score)) / float64(5)
	return strconv.FormatFloat(result, 'f', prec, 64)
}

func OverAllScore(score map[string]float64) float64 {
	result := 0.0
	cnt := 0
	for _, value := range score {
		result += value
		cnt++
	}

	return result / float64(cnt)
}
