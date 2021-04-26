package tools

import (
	"encoding/json"
	"io/ioutil"
)

//Page is anything that can be used to render an html template in our server. In this particular
//case, that will include users as well as pretty much blank pages
type Page interface {
	Save() error
}

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	Data     []byte `json:"data"`
}

func (u *User) Save() error {
	filename := "../internal/storage/user-storage/" + u.Name + ".json"
	file, err := json.MarshalIndent(u, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, []byte(file), 0644)
	if err != nil {
		return err
	}
	return nil
}

func NewUser(name, password string) *User {
	return &User{Name: name, Password: password}
}

func LoadUser(user string) (*User, error) {
	filepath := "../internal/storage/user-storage/" + user + ".json"
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	a := &User{}
	err = json.Unmarshal(file, a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

type EmptyPage struct {
	Title string //for html title purposes. But really it's useless
}

//this only exists as my crappy, hacky workaround to let this type fulfill the Page interface
func (e *EmptyPage) Save() error {
	return nil
}
