package example

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/colinjfw/sqlkit/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// Some sample fixtures here.
var fixtures = struct {
	users map[string]*User
}{
	users: map[string]*User{
		"one": &User{
			Email: "test@test.com",
		},
		"two": &User{
			Email: "test@test.com",
		},
	},
}

// Very basic way to resync fixtures.
func syncFixtures(repo *UserRepo) {
	for _, u := range fixtures.users {
		err := repo.Save(context.Background(), u)
		if err != nil {
			panic(err)
		}
	}
}

var users *UserRepo

// Wrap will wrap a test in a transaction, rolling back all changes so the DB is
// left in a clean state after the test has completed.
func wrap(t *testing.T) (db.TX, func()) {
	ctx, err := users.db.Begin(context.Background())
	require.Nil(t, err)
	return ctx, func() { ctx.Rollback() }
}

func TestMain(m *testing.M) {
	d, err := db.Open("sqlite3", ":memory:", db.WithLogger(db.StdLogger))
	if err != nil {
		panic(err)
	}
	users = &UserRepo{db: d}
	err = createTable(d)
	if err != nil {
		panic(err)
	}
	syncFixtures(users)

	code := m.Run()

	if err := d.Close(); err != nil {
		fmt.Printf("db shutdown error: %v", err)
	}
	os.Exit(code)
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

	err = users.Save(ctx, u)
	require.Nil(t, err)

	users, err := users.List(ctx)
	require.Nil(t, err)
	fmt.Printf("%+v\n", users)
}
