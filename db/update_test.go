package sql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdate_SQLValues(t *testing.T) {
	testSQL(t,
		"UPDATE users SET c1=?, c2=?",
		[]interface{}{"1", "2"},
		Update("users").
			Columns("c1", "c2").
			Values("1", "2"),
	)
}

func TestUpdate_SQLValue(t *testing.T) {
	testSQL(t,
		"UPDATE users SET c1=?",
		[]interface{}{"1"},
		Update("users").
			Table("users").
			Value("c1", "1"),
	)
}

func TestUpdate_SQLRecord(t *testing.T) {
	testSQL(t,
		"UPDATE users SET id=?, name=?",
		[]interface{}{1, "name"},
		Update("users").
			Record(struct {
				ID   int    `sql:"id"`
				Name string `sql:"name"`
			}{
				ID:   1,
				Name: "name",
			}),
	)
}

func TestUpdate_InvalidArgs(t *testing.T) {
	_, _, err := Update("users").
		Columns("id", "extra").
		Values(1).
		SQL()
	require.NotNil(t, err)
}
