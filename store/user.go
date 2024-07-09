package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/fchimpan/simple-server/entity"
	"github.com/go-sql-driver/mysql"
)

func (r *Repository) RegisterUser(ctx context.Context, db Execer, u *entity.User) error {
	u.CreatedAt = r.Clocker.Now()
	u.ModifiedAt = r.Clocker.Now()

	sql := `INSERT INTO user (name, password, role, created_at, modified_at) VALUES (?, ?, ?, ?, ?)`
	res, err := db.ExecContext(ctx, sql, u.Name, u.Password, u.Role, u.CreatedAt, u.ModifiedAt)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return fmt.Errorf("user already exists: %w", err)
		}
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("fatal get last insert id: %w", err)
	}
	u.ID = entity.UserID(id)
	return nil
}
