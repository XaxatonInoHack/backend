package review

import (
	"encoding/json"
	"fmt"
	"os"
)

type UseCase struct {
}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (*UseCase) ParseJSON() error {
	plan, err := os.ReadFile("internal/usecase/review/review_dataset.json")
	if err != nil {
		return err
	}

	var data []Review

	err = json.Unmarshal(plan, &data)
	if err != nil {
		return err
	}

	for _, review := range data {
		fmt.Println(review)
		break
	}

	return nil
}
