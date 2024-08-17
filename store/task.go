package store

import (
	"context"
	"errors"

	"github.com/fchimpan/simple-server/entity"
)

func (r *Repository) ListTasks(ctx context.Context, db Queryer, id entity.UserID) (entity.Tasks, error) {
	var tasks entity.Tasks
	sql := `SELECT id, title, status, created_at, modified_at FROM task WHERE user_id = ?;`

	if err := db.SelectContext(ctx, &tasks, sql, id); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) AddTask(ctx context.Context, db Execer, t *entity.Task) error {
	t.CreatedAt = r.Clocker.Now()
	t.ModifiedAt = r.Clocker.Now()

	sql := `INSERT INTO task (user_id, title, status, created_at, modified_at) VALUES (?, ?, ?, ?, ?);`
	res, err := db.ExecContext(ctx, sql, t.UserID, t.Title, t.Status, t.CreatedAt, t.ModifiedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	// 呼び出し元に auto increment された ID を返す
	t.ID = entity.TaskID(id)

	return nil
}

func (r *Repository) UpdateTask(ctx context.Context, db Execer, t *entity.Task) error {
	t.ModifiedAt = r.Clocker.Now()

	sql := `UPDATE task SET title = ?, status = ?, modified_at = ? WHERE id = ? AND user_id = ?;`
	res, err := db.ExecContext(ctx, sql, t.Title, t.Status, t.ModifiedAt, t.ID, t.UserID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("task not found")
	}
	return nil
}
