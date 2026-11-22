package repositories

import (
	"errors"
	"fmt"
	"main/src/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"context"
)

var ErrNotFound = errors.New("Task not found")

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func deserializeTask(row pgx.Row) (models.Task, error) {
	var t models.Task

	err := row.Scan(
		&t.Id,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt)

	return t, err
}

const columns = `id,title,description,status,created_at,updated_at`

func (r *TaskRepository) Create(
	ctx context.Context,
	title,
	description string,
) (models.Task, error) {
	return deserializeTask(
		r.db.QueryRow(
			ctx,
			`INSERT INTO tasks(title,description)
			VALUES($1,$2)
			RETURNING `+columns,
			title,
			description,
		),
	)
}

func (r *TaskRepository) List(
	ctx context.Context,
) ([]models.Task, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT `+columns+` FROM tasks ORDER BY id`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]models.Task, 0)

	for rows.Next() {
		t, e := deserializeTask(rows)

		if e != nil {
			return nil, e
		}

		result = append(result, t)
	}

	return result, rows.Err()
}

func (r *TaskRepository) Get(
	ctx context.Context,
	id int64,
) (models.Task, error) {
	t, err := deserializeTask(
		r.db.QueryRow(
			ctx,
			`SELECT `+columns+` FROM tasks WHERE id=$1`,
			id,
		),
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return t, ErrNotFound
	}

	return t, err
}

func (r *TaskRepository) Update(
	ctx context.Context,
	id int64,
	title,
	description,
	status string,
) (models.Task, error) {
	t, err := deserializeTask(
		r.db.QueryRow(
			ctx,
			`UPDATE tasks
			SET
				title=$2,
				description=$3,
				status=$4,
				updated_at=NOW()
			WHERE id=$1
			RETURNING `+columns,
			id,
			title,
			description,
			status,
		),
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return t, ErrNotFound
	}

	return t, err
}

func (r *TaskRepository) UpdateStatus(
	ctx context.Context,
	id int64,
	status string,
) (models.Task, error) {
	t, err := deserializeTask(
		r.db.QueryRow(
			ctx,
			`UPDATE tasks
			SET
				status=$2,
				updated_at=NOW()
			WHERE id=$1
			RETURNING `+columns,
			id,
			status,
		),
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return t, ErrNotFound
	}

	return t, err
}

func (r *TaskRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	tag, err := r.db.Exec(
		ctx,
		`DELETE FROM tasks WHERE id=$1`,
		id,
	)

	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
