package review

type Review struct {
	UserID   int64  `json:"ID_under_review"`
	ReviewID int64  `json:"ID_reviewer"`
	Feedback string `json:"review"`
}
