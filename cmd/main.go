package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Jonny-Burkholder/login-test-2/internal/tools"
)

const (
	logPath string = "internal/log/log.txt"
)

var funcMap = template.FuncMap{
	//add template funcs here as needed
}

var templates = template.Must(template.New("*").Funcs(funcMap).ParseGlob("../web/html/*.html"))

func myRecoverFunc() {
	if err := recover(); err != nil {
		log.Printf("Recovered from panic: %v", err)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p tools.Page) { //should eventually use a buffer with this
	defer myRecoverFunc() //I figure renderTemplate is getting called by pretty much every handler, so it's more or less safe to put here
	buffer := tools.GetBuf()
	defer tools.PutBuf(buffer) //deferring this just in case of panic
	err := templates.ExecuteTemplate(buffer, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buffer.WriteTo(w)
	return
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	defer myRecoverFunc()
	log.Println("Redirecting user to login")
	http.Redirect(w, r, "/login", http.StatusFound)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session") //check to see if the user already has an active session
	if err != nil {                    //if there is no session cookie, redirect user to login
		log.Println(err)
		page := &tools.EmptyPage{Title: "Account Login"} //This is totally unnecessary, but it fulfills the template pattern
		renderTemplate(w, "login", page)
		return
	}
	//if there is a login cookie, redirect user to homepage
	path := "/home/" + cookie.Value
	http.Redirect(w, r, path, http.StatusFound)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	defer myRecoverFunc()
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(cookie)
	log.Println("Logging user out")
	tools.LogOut(w, cookie)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	page := &tools.EmptyPage{Title: "Create Account"} //this should be calling a function and EmptyPage shouldn't be exported
	renderTemplate(w, "create-account", page)
}

func handleNewUser(w http.ResponseWriter, r *http.Request) {
	defer myRecoverFunc()
	fmt.Println("Hello, creating a new user here!")
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println(username)
	fmt.Println(password)
	fmt.Println(r.FormValue("password"))
	fmt.Println(r.FormValue("confirm-password"))
	user := tools.NewUser(username, password)
	err := user.Save()
	if err != nil {
		log.Panic(err)
	}
	tools.Login(w, username)
	path := "/home/" + username
	http.Redirect(w, r, path, http.StatusFound)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session")
	if err != nil {
		log.Println(err)
		http.Error(w, "Access denied", 403)
		return
	}
	title := r.URL.Path[len("/home/"):]
	user, err := tools.LoadUser(title)
	if err != nil {
		log.Panic(err)
	}
	renderTemplate(w, "home", user)
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	defer myRecoverFunc()
	//probably need to use a sync pool to keep the counter alive here
	c := &tools.Counter{} //this is currently useless, since I haven't decided yet how to pass it between login and validate
	//c.Reset() //for use with syncpool
	err := tools.CheckPassword(r, c) //is this hiding complexity? Should I establish form values as variables before this point?
	if err != nil {
		http.Error(w, err.Error(), http.StatusFound)
		return
	}
	fmt.Println(r.FormValue("username"))
	tools.Login(w, r.FormValue("username"))
	path := "/home/" + r.FormValue("username")
	fmt.Println(path)
	http.Redirect(w, r, path, http.StatusFound)
}

func main() {
	/*
		//UNCOMMENT TO ENABLE LOGGING TO LOG.TXT RATHER THAN STD.OUT
		f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		check(err)
		log.SetOutput(f)
	*/

	fmt.Println("Now serving on port 8080")
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/")))

	http.Handle("/static/", static)
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/login/", handleLogin)
	http.HandleFunc("/home/", handleHome)
	http.HandleFunc("/validate", handleValidate)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/new-user", handleNewUser)
	http.HandleFunc("/create-account", handleCreateAccount)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
