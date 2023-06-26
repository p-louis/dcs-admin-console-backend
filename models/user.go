package models

import (
	"errors"
	"github.com/p-louis/dcs-admin/utils/token"
	"os"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginCheck(username string, password string) (string, error) {

	var err error

	//u := User{}

	if username != os.Getenv("ADMIN_USERNAME") || password != os.Getenv("ADMIN_PASSWORD") {
		err = errors.New("Invalid Username/Password")
	}

	if err != nil {
		return "", err
	}

	token, err := token.GenerateToken(1)

	if err != nil {
		return "", err
	}

	return token, nil

}

func GetUserByID(uid uint) (User, error) {

	var u User

	if uid != 1 {
		return u, errors.New("User not found!")
	}

	u.Username = os.Getenv("ADMIN_USERNAME")

	return u, nil

}
