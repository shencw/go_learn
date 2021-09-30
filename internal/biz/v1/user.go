package v1

import "github.com/pkg/errors"

type User struct {
	UserName string
	Password string
}

func NewUser() *User {
	return &User{}
}

func (u *User) Get(username string) (*User, error) {
	if username != "admin" {
		return nil, errors.New("不存在的用户")
	}

	return &User{
		UserName: username,
		Password: "admin",
	}, nil
}

func (u *User) Compare(pwd string) error {
	if u.Password != pwd {
		return errors.New("密码错误")
	}

	return nil
}
