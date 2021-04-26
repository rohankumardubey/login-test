package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Counter struct {
	Count int
	Max   int
}

//Counter counts the number of login attempts a person has made
func (c *Counter) CountUp() error {
	c.Count++
	if c.Count == c.Max {
		return errors.New("Maximum number of attempts reached. Account locked for 1 minute")
	}
	return nil
}

//Resest sets the counter back to zero
func (c *Counter) Reset() {
	c.Count = 0
}

//CookieSetter creates a new login cookie, and automatically sets it
func Login(w http.ResponseWriter, s string) { //this should eventually return an error
	cookie := http.Cookie{
		Name:     "session",
		Value:    s,
		MaxAge:   60 * 5, //5 minutes for our testing purposes
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

//LogOut deletes the cookie set by Login
func LogOut(w http.ResponseWriter, c *http.Cookie) {
	c.Expires = time.Now()
	c.MaxAge = -1
	http.SetCookie(w, c)
}

func CheckPassword(r *http.Request, c *Counter) error {
	username := r.FormValue("username")
	password := r.FormValue("password")
	userpath := "../internal/storage/user-storage/" + username + ".json"
	file, err := ioutil.ReadFile(userpath)
	if err != nil {
		return err
	}
	user := &User{}
	err = json.Unmarshal([]byte(file), user)
	if err != nil {
		return err
	}
	if password != user.Password {
		err = c.CountUp()
		if err != nil {
			return err
		}
		return fmt.Errorf("Incorrect password. %v more attempts until account is locked.", c.Count)
	}
	return nil
}
