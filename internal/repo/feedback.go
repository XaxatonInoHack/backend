package repo

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"xaxaton/internal/domain/model"
)

func (s *Storage) CreateFeedback(ctx context.Context,
	userId int, score string,
	result string, resume string,
) error {
	const op = "repo.CreateReview"

	var (
		query = `
		INSERT INTO review (user_id, score, result, resume)
		VALUES ($1, $2, $3, $4)
	`
		insertValues = []any{
			userId, score, result, resume,
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

func (s *Storage) GetFeedback(ctx context.Context,
	userId int,
) (model.Feedback, error) {
	const op = "repo.GetFeedback"

	var (
		query = `
		SELECT user_id, score, result, resume
		FROM feedback
		WHERE user_id = $1
	`
		feedback model.Feedback
	)

	err := s.pool.QueryRow(ctx, query, userId).
		Scan(&feedback.UserID, &feedback.Score, &feedback.Result, &feedback.Resume)
	if err != nil {
		return model.Feedback{}, fmt.Errorf("%s:%w", op, err)
	}
	return feedback, nil
}
