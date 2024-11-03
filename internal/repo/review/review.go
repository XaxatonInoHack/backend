package review

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

func (s *Storage) CreateReview(ctx context.Context,
	reviews []model.Review,
) error {
	const op = "repo.CreateReview"

	batch := &pgx.Batch{}
	query := `INSERT INTO review 
    (user_id, review_id, feedback, period) 
	VALUES ($1, $2, $3, $4)
	ON CONFLICT DO NOTHING
`

	for _, review := range reviews {
		batch.Queue(query, review.UserID, review.ReviewID, review.Feedback, review.Period)
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

func (s *Storage) GetReview(ctx context.Context,
	userId int,
) ([]model.Review, error) {
	const op = "repo.GetReview"

	var (
		query = `
		SELECT user_id, review_id, feedback, period
		FROM review
		WHERE user_id = $1
	`
		reviews []model.Review
	)

	rows, err := s.pool.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var review model.Review
		if err = rows.Scan(&review.UserID, &review.ReviewID, &review.Feedback, &review.Period); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return reviews, nil
}
