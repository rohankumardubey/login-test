package biscuit

import (
	"fmt"
	"net/http"
	"time"
)

//CookieSetter creates a new login cookie, and automatically sets it
func SetLogin(w http.ResponseWriter, s *Session) { //this should eventually return an error
	cookie := http.Cookie{
		Name:     sessionCookieName,
		Value:    fmt.Sprint(s.ID),
		MaxAge:   SessionLength, //Each client should decide how session length should be handled
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
