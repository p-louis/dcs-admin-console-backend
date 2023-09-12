package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/p-louis/dcs-admin/utils/token"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenBody struct {
	Token string `json:"token"`
}

type RoleBody struct {
	Role string `json:"message"`
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func LoginCheck(username string, password string) (string, error) {

	var err error

	if username != os.Getenv("ADMIN_USERNAME") {
		_, err = VerifyExternalUser(username, password)
	} else if password != os.Getenv("ADMIN_PASSWORD") {
		err = errors.New("Invalid Username/Password")
	}

	if err != nil {
		return "", err
	}

	token, err := token.GenerateToken(24)

	if err != nil {
		return "", err
	}

	return token, nil

}

func VerifyExternalUser(username string, password string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://hq.x51squadron.com/fileapi/login?ver=0.0.0", bytes.NewBuffer([]byte("")))

	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))

	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		err := errors.New("Invalid Username/Password")
		if err != nil {
			return "", err
		}
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	var token TokenBody
	json.Unmarshal(body, &token)

	req, err = http.NewRequest("GET", "https://hq.x51squadron.com/fileapi/validate-token", bytes.NewBuffer([]byte("")))

	if err != nil {
		return "", err
	}
	req.Header.Add("x-access-tokens", token.Token)

	res, err = client.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	var role RoleBody
	json.Unmarshal(body, &role)

	if role.Role != "curator" {
		err = errors.New("User has insufficient permissions")
	}

	if err != nil {
		return "", err
	}

	return "ok", nil
}

func GetUserByID(uid uint) (User, error) {

	var u User

	if uid != 1 {
		return u, errors.New("User not found!")
	}

	u.Username = os.Getenv("ADMIN_USERNAME")

	return u, nil

}
