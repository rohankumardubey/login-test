//Authorization handles user login and permission to see pages
package middleware

import "fmt"

//Permission accepts as an interface any type that has a method to add a permission and returns an error
type Permission map[string]bool

//Check first looks to see if the permission exists, then returns the associated bool
func (p Permission) Check(s string) bool {
	ok := p[s]
	if ok != true {
		return ok
	}
	return p[s]
}

//Add adds a permission to the permission map
func (p Permission) Add(s string, b bool) error {
	ok := p[s]
	if ok != true {
		p[s] = b
		return nil
	}
	return fmt.Errorf("Permission %v already exists. State is %v", s, p[s])
}

//Toggle checks the current state of a permission and reverses it
func (p Permission) Toggle(s string) error {
	ok := p[s]
	if ok != true {
		return fmt.Errorf("Error: Permission %v not found", s)
	}
	if p[s] == true {
		p[s] = false
	} else {
		p[s] = true
	}
	return nil
}
