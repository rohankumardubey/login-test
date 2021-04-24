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

var templates = template.Must(template.New("*").Funcs(funcMap).ParseGlob("./web/html/*.html"))

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
	//check cookies to see if user is already logged in
	//if so, redirect them to home/user
	//if not, execute normally
	title := r.URL.Path[len("/"):]
	page := &tools.EmptyPage{Title: title}
	renderTemplate(w, "login", page)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	//somewhere in this process, check the cookie to see if the user is logged in
	//if not, either serve 403 forbidden, or just redirect to login
	//I mean honestly, I never see 403 anymore, I always just get redirected. But hey, this is my
	//app, and what fun is it if I don't get to serve an error code every now and then, huh?
	title := r.URL.Path[len("/home/"):]
	user := tools.LoadUser(title)
	renderTemplate(w, "home", user)
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

	log.Fatal(http.ListenAndServe(":8080", nil))
}
