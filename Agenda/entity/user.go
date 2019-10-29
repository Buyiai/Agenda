package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type User struct {
	Username string
	Password string
	Email    string
	Contact  string
}

var curUser User

type UserArray []User

var userList UserArray

const userListFileName string = "userList.txt"
const curUserFileName string = "curUser.txt"

func WriteUserListToFile() error {
	file, err := os.OpenFile(userListFileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	if err := enc.Encode(userList); err != nil {
		return err
	}
	return nil
}

func ReadUserListFromFile() error {
	file, err := os.Open(userListFileName)
	if err != nil {
		return nil
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	if err := dec.Decode(&userList); err != nil {
		return err
	}
	return nil
}

func WriteCurUserToFile() error {
	file, err := os.OpenFile(curUserFileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	if err := enc.Encode(curUser); err != nil {
		return err
	}
	return nil
}

func ReadCurUserFromFile() error {
	file, err := os.Open(curUserFileName)
	if err != nil {
		return nil
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	if err := dec.Decode(&curUser); err != nil {
		return err
	}
	return nil
}

func CheckLogined() error {
	if err := ReadCurUserFromFile(); err != nil {
		return err
	}
	if curUser != (User{}) {
		return errors.New("a user has been logined, logout and try again")
	}
	return nil
}

func CheckLogouted() error {
	if err := ReadCurUserFromFile(); err != nil {
		return err
	}
	if curUser == (User{}) {
		return errors.New("no user has been logined, login and try again")
	}
	return nil
}

func RegisterUser(username, password, email, contact string) error {
	if username == "" || password == "" || email == "" || contact == "" {
		return errors.New("a username, a password, an email and a phone required")
	}
	if err := CheckLogined(); err != nil {
		return err
	}
	if err := ReadUserListFromFile(); err != nil {
		return err
	}
	for _, v := range userList {
		if v.Username == username {
			return errors.New("username already exists")
		}
	}
	user := User{username, password, email, contact}
	userList = append(userList, user)
	if err := WriteUserListToFile(); err != nil {
		return err
	}
	return nil
}

func LoginUser(username, password string) error {
	if username == "" || password == "" {
		return errors.New("a username and a password required")
	}
	if err := CheckLogined(); err != nil {
		return err
	}
	if err := ReadUserListFromFile(); err != nil {
		return err
	}
	for _, v := range userList {
		if v.Username == username && v.Password == password {
			curUser = v
			if err := WriteCurUserToFile(); err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("username or password error")
}

func LogoutUser() error {
	if err := CheckLogouted(); err != nil {
		return err
	}
	curUser = User{}
	if err := WriteCurUserToFile(); err != nil {
		return err
	}
	return nil
}

func (ul UserArray) String() string {
	out := "\n         username               email          phone\n"
	for i, v := range ul {
		out = fmt.Sprintf("%s%2d%15s%20s%15s\n", out, i+1, v.Username, v.Email, v.Contact)
	}
	return out
}

/*
func GetCurUser() (User, error) {
  	if err := ReadCurUserFromFile(); err != nil {
		return User{}, err
	}
  	return curUser, nil
}*/

// func GetUserList() UserArray {
// 	return userList
// }
