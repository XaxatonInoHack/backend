package feedback

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"xaxaton/internal/model"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: pool,
	}
}

func (s *Storage) CreateFeedback(ctx context.Context,
	feedbacks []model.Feedback,
) error {
	const op = "repo.CreateFeedback"

	batch := &pgx.Batch{}
	query := `INSERT INTO feedback (user_id, score, result, resume) VALUES ($1, $2, $3, $4)`

	for _, feedback := range feedbacks {
		batch.Queue(query, feedback.UserID, feedback.Score, feedback.Result, feedback.Resume)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	defer tx.Rollback(ctx)

	br := tx.SendBatch(ctx, batch)
	if err = br.Close(); err != nil {
		return fmt.Errorf("%s: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (s *Storage) GetFeedback(ctx context.Context,
	userId int,
) ([]model.Feedback, error) {
	const op = "repo.GetFeedback"

	var (
		query = `
		SELECT user_id, score, result, resume
		FROM feedback
		WHERE user_id = $1
	`
		feedbacks []model.Feedback
	)

	rows, err := s.pool.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var feedback model.Feedback
		if err = rows.Scan(&feedback.UserID, &feedback.Score, &feedback.Result, &feedback.Resume); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		feedbacks = append(feedbacks, feedback)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return feedbacks, nil
}
