package entity

import (
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Roles    []string  `db:"roles"`
}

func NewUser(username string) (*User, error) {
	user := &User{
		Id:       uuid.New(),
		Username: username,
		Roles:    []string{"user"},
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) Validate() error {
	if u.Username == "" {
		return ErrEmptyUsername
	}

	if len(u.Roles) == 0 {
		u.Roles = []string{"user"}
	}

	return nil
}
