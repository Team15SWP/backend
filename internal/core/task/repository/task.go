package repository

import (
	"context"
	"fmt"
	"strings"

	"study_buddy/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ TaskRepository = (*TaskRepo)(nil)

type TaskRepo struct {
	db *pgxpool.Pool
}

func NewTaskRepo(db *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{
		db: db,
	}
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task *model.GeneratedTask) error
}

func (t *TaskRepo) CreateTask(ctx context.Context, task *model.GeneratedTask) error {
	pool, err := t.db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("t.db.Acquire: %w", err)
	}
	defer pool.Release()

	columns := []string{"user_id", "task", "description", "solution", "hint1", "hint2", "hint3", "difficulty", "solved"}

	query, args, err := sq.Insert("tasks").
		Columns(columns...).
		Values(task.UserID, task.TaskName, task.TaskDescription, task.Solution, task.Hints.Hint1, task.Hints.Hint2,
			task.Hints.Hint3, task.Difficulty, task.Solved).
		Suffix("RETURNING id, " + strings.Join(columns, ", ")).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("format user insert SQL: %w", err)
	}

	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}
	return nil
}
