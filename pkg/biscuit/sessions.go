package biscuit

import (
	"math/rand"
	"sync"
	"time"
)

const (
	sessionCookieName string = "SESS_bsct"
)

//This way each client can decide how long they'd like their sessions to last
var SessionLength int

//Since this is not safe for concurrent use, this should be set only once
func SetSessionLength(i int) {
	SessionLength = i
}

//Defines how many login attempts a user has
var MaxAttempts int

//Not safe for concurrent use. Should be set once per server
func SetMaxAttempts(i int) {
	MaxAttempts = i
}

//This interface allows us to accept users from different user handling packages
type User interface {
	Save() error
}

//Sessions is a map that holds information about active sessions
type Sessions struct {
	Mux      sync.Mutex
	Sessions map[uint64]*Session
}

//Session holds information about a given session
type Session struct {
	Mux      sync.Mutex
	Username string
	ID       uint64
	Counter  *counter
	Alive    bool
}

//NewSession takes a username argument, sets a session cookie in browser
func (s *Sessions) NewSession(username string) uint64 { //my gut tells me to return an error here, but I can't really think what error I would need to return
	s.Mux.Lock()

	counter := getCounter(MaxAttempts)

	session := Session{Username: username, Counter: counter, Alive: false}

	for { //this loops until a unique id is generated. Just in case
		id := newSessionID()
		_, ok := s.Sessions[id]
		if ok != true {
			s.Sessions[id] = &session
			s.Mux.Unlock()
			return id
		}
	}
}

//generates a random uint64 as a session ID
func newSessionID() uint64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Uint64()
}
