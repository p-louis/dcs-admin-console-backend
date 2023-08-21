package models

import (
	"encoding/base64"
	"errors"
	"github.com/p-louis/dcs-admin/utils/token"
	"net/http"
	"os"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func LoginCheck(username string, password string) (string, error) {

	var err error

	if username != os.Getenv("ADMIN_USERNAME") {
		req, err := http.Get("https://x51squadron.com/login?ver=0.0.0")

		if err != nil {
			return "", err
		}

		req.Header.Add("Authorization", "Basic "+basicAuth(username, password))
		if req.StatusCode != 200 {
			err = errors.New("Invalid Username/Password")
		}

	} else if password != os.Getenv("ADMIN_PASSWORD") {
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
