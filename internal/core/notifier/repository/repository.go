package repository

import (
	"context"

	"study_buddy/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ Repository = (*NotifyRepo)(nil)

type NotifyRepo struct {
	db *pgxpool.Pool
}

func NewNotifyRepo(db *pgxpool.Pool) *NotifyRepo {
	return &NotifyRepo{
		db: db,
	}
}

type Repository interface {
	GetAllUsersEmail(ctx context.Context) ([]*model.User, error)
}

func (n *NotifyRepo) GetAllUsersEmail(ctx context.Context) ([]*model.User, error) {
	pool, err := n.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("name", "email").
		From("users").
		OrderBy("id").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		user := &model.User{}
		if err = rows.Scan(&user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
