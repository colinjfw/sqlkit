package sql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsert_SQLValues(t *testing.T) {
	testSQL(t,
		"INSERT INTO users (c1, c2) VALUES (?, ?)",
		[]interface{}{"1", "2"},
		Insert().
			Into("users").
			Columns("c1", "c2").
			Values("1", "2"),
	)
}

func TestInsert_SQLValue(t *testing.T) {
	testSQL(t,
		"INSERT INTO users (c1) VALUES (?)",
		[]interface{}{"1"},
		Insert().
			Into("users").
			Value("c1", "1"),
	)
}

func TestInsert_SQLRecord(t *testing.T) {
	testSQL(t,
		"INSERT INTO users (id, name) VALUES (?, ?)",
		[]interface{}{1, "name"},
		Insert().
			Into("users").
			Record(struct {
				ID   int    `sql:"id"`
				Name string `sql:"name"`
			}{
				ID:   1,
				Name: "name",
			}),
	)
}

func TestInsert_InvalidArgs(t *testing.T) {
	_, _, err := Insert().
		Into("users").
		Columns("id", "extra").
		Values(1).
		SQL()
	require.NotNil(t, err)
}
