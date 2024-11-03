package self_review

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

func (s *Storage) CreateSelfReview(ctx context.Context,
	selfReviews []model.SelfReview,
) error {
	const op = "repo.CreateSelfReview"

	batch := &pgx.Batch{}
	query := `INSERT INTO self_review (user_id, score, result, resume) VALUES ($1, $2, $3, $4)`

	for _, self_review := range selfReviews {
		batch.Queue(query, self_review.UserID, self_review.Score, self_review.Result, self_review.Resume)
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

func (s *Storage) GetSelfReviews(ctx context.Context,
	userId int,
) ([]model.SelfReview, error) {
	const op = "repo.GetSelfReviews"

	var (
		query = `
		SELECT user_id, score, result, resume
		FROM self_review
		WHERE user_id = $1
	`
		selfReviews []model.SelfReview
	)

	rows, err := s.pool.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var self_review model.SelfReview
		if err = rows.Scan(&self_review.UserID, &self_review.Score, &self_review.Result, &self_review.Resume); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		selfReviews = append(selfReviews, self_review)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return selfReviews, nil
}

func (s *Storage) InsertSelfScore(ctx context.Context,
	selfReviews []model.SelfReview,
) error {
	const op = "repo.InsertSelfScore"

	batch := &pgx.Batch{}
	query := `INSERT INTO self_review (user_id, score) VALUES ($1, $2)`

	for _, self_review := range selfReviews {
		batch.Queue(query, self_review.UserID, self_review.Score)
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

func (s *Storage) UpdateSelfResume(ctx context.Context,
	selfReviews []model.SelfReview,
) error {
	const op = "repo.UpdateSelfResume"

	var (
		query = `
			UPDATE self_review
			SET result = $2, resume = $3
			WHERE user_id = $1
		`
	)

	batch := &pgx.Batch{}

	for _, self_review := range selfReviews {
		batch.Queue(query, self_review.UserID, self_review.Result, self_review.Resume)
	}
	br := s.pool.SendBatch(ctx, batch)
	defer br.Close()

	for range selfReviews {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
