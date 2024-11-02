package repo

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"xaxaton/internal/domain/model"
)

func (s *Storage) CreateReview(ctx context.Context,
	userId int, reviewerId int,
	feedback string, period string,
) error {
	const op = "repo.CreateReview"

	var (
		query = `
		INSERT INTO review (user_id, review_id, feedback, period)
		VALUES ($1, $2, $3, $4)
	`
		insertValues = []any{
			userId, reviewerId, feedback, period,
		}
	)

	//TODO: переделать на адекватный логгер мб
	fmt.Println("beginning transaction")
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", err)
		return fmt.Errorf("%s:%w", op, err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query, insertValues...)
	if err != nil {
		log.Error("failed to insert review", err)
		return fmt.Errorf("%s:%w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error("failed to commit transaction", err)
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (s *Storage) GetReview(ctx context.Context,
	userId int,
) (model.Review, error) {
	const op = "repo.GetReview"

	var (
		query = `
		SELECT user_id, review_id, feedback, period
		FROM review
		WHERE user_id = $1
	`
		review model.Review
	)

	err := s.pool.QueryRow(ctx, query, userId).Scan(&review.UserID, &review.ReviewID, &review.Feedback, &review.Period)
	if err != nil {
		return model.Review{}, fmt.Errorf("%s:%w", op, err)
	}
	return review, nil
}
