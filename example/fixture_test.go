package example

import (
	"context"
)

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

func syncFixtures(repo *UserRepo) {
	for _, u := range fixtures.users {
		err := repo.Save(context.Background(), u)
		if err != nil {
			panic(err)
		}
	}
}
