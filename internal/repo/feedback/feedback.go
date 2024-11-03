package feedback

import (
	"context"
	"fmt"
	"strings"

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
	userId int64,
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

func (s *Storage) GetFeedbackAll(ctx context.Context,
	userIDs []int64,
) (map[int64]model.Feedback, error) {
	const op = "repo.GetFeedbackAll"

	// Проверяем, что срез userIDs не пуст
	if len(userIDs) == 0 {
		return nil, fmt.Errorf("%s: no user IDs provided", op)
	}

	// Генерируем плейсхолдеры для запроса
	placeholders := make([]string, len(userIDs))
	args := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	// Формируем запрос с использованием плейсхолдеров
	query := fmt.Sprintf(`
        SELECT user_id, score, result, resume
        FROM feedback
        WHERE user_id IN (%s)
    `, strings.Join(placeholders, ", "))

	// Инициализируем карту для результатов
	feedbacksMap := make(map[int64]model.Feedback)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var feedback model.Feedback
		if err = rows.Scan(&feedback.UserID, &feedback.Score, &feedback.Result, &feedback.Resume); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		// Добавляем feedback в карту, перезаписывая предыдущий, если такой уже есть
		feedbacksMap[feedback.UserID] = feedback
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return feedbacksMap, nil
}

func (s *Storage) InsertScore(ctx context.Context,
	feedbacks []model.Feedback,
) error {
	const op = "repo.InsertScore"

	batch := &pgx.Batch{}
	query := `INSERT INTO feedback (user_id, score) VALUES ($1, $2)`

	for _, feedback := range feedbacks {
		batch.Queue(query, feedback.UserID, feedback.Score)
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

func (s *Storage) UpdateResume(ctx context.Context,
	feedbacks []model.Feedback,
) error {
	const op = "repo.UpdateResume"

	var (
		query = `
			UPDATE feedback
			SET result = $2, resume = $3, score = $4
			WHERE user_id = $1
		`
	)

	batch := &pgx.Batch{}

	for _, feedback := range feedbacks {
		batch.Queue(query, feedback.UserID, feedback.Result, feedback.Resume, feedback.Score)
	}
	br := s.pool.SendBatch(ctx, batch)
	defer br.Close()

	for range feedbacks {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
