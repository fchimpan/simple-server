package store

import (
	"context"
	"testing"

	"github.com/fchimpan/simple-server/clock"
	"github.com/fchimpan/simple-server/entity"
	"github.com/fchimpan/simple-server/testutil"
	"github.com/google/go-cmp/cmp"
)

func TestRepository_ListTasks(t *testing.T) {
	t.Skip()
	ctx := context.Background()

	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}

	wants := prepareTasks(ctx, t, tx)

	r := &Repository{}
	gots, err := r.ListTasks(ctx, tx, 1)
	if err != nil {
		t.Fatal(err)
	}

	if d := cmp.Diff(wants, gots); d != "" {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func prepareTasks(ctx context.Context, t *testing.T, con Execer) entity.Tasks {
	t.Helper()
	// clean up
	if _, err := con.ExecContext(ctx, "DELETE FROM task;"); err != nil {
		t.Logf("failed to initialize task: %v", err)
	}
	c := clock.FakeClocker{}
	wants := entity.Tasks{
		{
			Title: "want task 1", Status: "todo",
			CreatedAt: c.Now(), ModifiedAt: c.Now(),
		},
		{
			Title: "want task 2", Status: "todo",
			CreatedAt: c.Now(), ModifiedAt: c.Now(),
		},
		{
			Title: "want task 3", Status: "done",
			CreatedAt: c.Now(), ModifiedAt: c.Now(),
		},
	}
	result, err := con.ExecContext(ctx,
		`INSERT INTO task (title, status, created_at, modified_at)
			VALUES
			    (?, ?, ?, ?),
			    (?, ?, ?, ?),
			    (?, ?, ?, ?);`,
		wants[0].Title, wants[0].Status, wants[0].CreatedAt, wants[0].ModifiedAt,
		wants[1].Title, wants[1].Status, wants[1].CreatedAt, wants[1].ModifiedAt,
		wants[2].Title, wants[2].Status, wants[2].CreatedAt, wants[2].ModifiedAt,
	)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)
	return wants
}
