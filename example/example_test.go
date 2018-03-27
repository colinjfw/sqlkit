package example

import (
	"context"
	"fmt"
	"testing"

	"github.com/coldog/sqlkit/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

var users *UserRepo

func init() {
	d, err := db.Open("sqlite3", ":memory:", db.WithLogger(db.StdLogger))
	if err != nil {
		panic(err)
	}
	users = &UserRepo{db: d}
	err = d.Exec(context.Background(), db.Raw(
		`create table users (
			id int primary key,
			email text,
			created_at timestamp,
			updated_at timestamp
		)`,
	)).Err()
	syncFixtures(users)
}

func wrap(t *testing.T) (db.TX, func()) {
	ctx, err := users.db.Begin(context.Background())
	require.Nil(t, err)
	return ctx, func() { ctx.Rollback() }
}

func TestExample(t *testing.T) {
	ctx, done := wrap(t)
	defer done()

	u := &User{Email: "test@test.com"}
	err := users.Save(ctx, u)
	require.Nil(t, err)

	u, err = users.Get(ctx, fixtures.users["one"].ID)
	require.Nil(t, err)

	fmt.Printf("%+v\n", u)
}
