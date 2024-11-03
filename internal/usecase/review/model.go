package review

import "strconv"

type Review struct {
	UserID   int64  `json:"ID_under_review"`
	ReviewID int64  `json:"ID_reviewer"`
	Feedback string `json:"review"`
}

func employeeScoreToDB(score map[string]float64) string {
	result := ""
	for key, value := range score {
		result += key + ": " + strconv.FormatFloat(value, 'f', 1, 64) + " "
	}

	return result
}

func scoreToResult(score map[string]float64) string {
	result := 0.0
	count := 0
	for _, value := range score {
		result += value
		count += 1
	}

	result = result / float64(count)
	return strconv.FormatFloat(result, 'f', 1, 64)
}
