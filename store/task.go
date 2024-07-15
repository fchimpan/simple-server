package store

import (
	"context"

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
