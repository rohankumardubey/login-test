package tools

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
	//save user
	return nil
}

func NewUser() *User {
	return &User{}
}

func LoadUser(user string) *User {
	//json to load user, yay
	//return user, not nil!
	return nil
}

type EmptyPage struct {
	Title string //for html title purposes
}

//this only exists as my crappy, hacky workaround to let this be used as type Page
func (e *EmptyPage) Save() error {
	return nil
}
