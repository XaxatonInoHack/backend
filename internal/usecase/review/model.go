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
