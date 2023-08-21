package models

import (
	"bytes"
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
		err = errors.New("Invalid Username/Password")
		client := &http.Client{}
		req, err := http.NewRequest("POST", "https://uploader.x51squadron.com/login?ver=0.0.0", bytes.NewBuffer([]byte("")))

		if err != nil {
			return "", err
		}

		req.Header.Add("Authorization", "Basic "+basicAuth(username, password))

		res, err := client.Do(req)
		if res.StatusCode != 200 || err != nil {
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
