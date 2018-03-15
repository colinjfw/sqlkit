package sql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshal_CustomType(t *testing.T) {
	db, err := Open("sqlite3", ":memory:", WithLogger(StdLogger))
	require.Nil(t, err)

	obj := &struct {
		ID    int    `sql:"id"`
		Data  int    `sql:"data"`
		Other string `sql:"other"`
	}{}

	ctx := context.Background()

	err = db.Exec(ctx,
		Raw("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY, data INT, other TEXT)"),
	).Err()
	require.Nil(t, err)

	err = db.Exec(ctx, db.
		Insert().
		Into("users").
		Value("id", 1).
		Value("other", "hello"),
	).Err()
	require.Nil(t, err)

	err = db.Query(ctx, db.Select("users.*").From("users")).First(obj)
	require.Nil(t, err)

	require.Equal(t, 1, obj.ID)
	require.Nil(t, db.Close())
}
