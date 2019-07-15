package example

import (
	"context"
	"time"

	"github.com/colinjfw/sqlkit/db"
)

func createTable(d db.DB) error {
	return d.Exec(context.Background(), db.Raw(
		`create table users (
			id integer primary key autoincrement,
			email text,
			created_at timestamp default current_timestamp not null,
			updated_at timestamp default current_timestamp not null
		)`,
	)).Err()
}

// Base is a base type for database objects.
type Base struct {
	ID       int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// User is a database object.
type User struct {
	Base
	Email string
}

// UserRepo represents the user repository.
type UserRepo struct {
	db db.DB
}

func (rep *UserRepo) table() string { return "users" }

// Save a user.
func (rep *UserRepo) Save(ctx context.Context, u *User) error {
	if u.ID == 0 {
		r := rep.db.Exec(ctx, rep.db.Insert().
			Into(rep.table()).
			Record(u, "email"))
		u.ID = int(r.LastID)
		return r.Err()
	}

	u.UpdatedAt = time.Now()
	r := rep.db.Exec(ctx, rep.db.Update(rep.table()).
		Record(u, "email", "updated_at").
		Where("id = ?", u.ID))
	return r.Err()
}

// Get will fetch a user.
func (rep *UserRepo) Get(ctx context.Context, id int) (*User, error) {
	u := &User{}
	err := rep.db.Query(
		ctx,
		rep.db.Select("*").From(rep.table()).Where("id = ?", id),
	).Decode(u)
	return u, err
}

// List will list a set of users.
func (rep *UserRepo) List(ctx context.Context) ([]*User, error) {
	u := []*User{}
	err := rep.db.Query(
		ctx,
		rep.db.Select("*").From(rep.table()),
	).Decode(&u)
	return u, err
}
