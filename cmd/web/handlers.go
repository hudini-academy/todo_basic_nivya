package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"todo/pkg/models"
	"unicode/utf8"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	// Checking the path is in home or not
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	s, err := app.todo.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}
	// panic("oops! something went wrong") // Deliberate panic
	// Using String Slice to add files which will give the template.
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl"}
	ts, err := template.ParseFiles(files...)

	//  checking for any error
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	message := app.Session.PopString(r, "Flash");

	err = ts.Execute(w, struct {
		Tasks []*models.Todo
		Flash string
	}{
		Tasks: s,
		Flash: message,
	})
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server error", 500)
	}

}

func (app *application) AddTask(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("TaskName")
	details := r.FormValue("Details")
	expires := "7"

	// Check that the title field is not blank and is not more than 100
	if !app.validatetask(r,name) && !app.validatetask(r,details) && !app.validatetask(r,expires){
	app.Session.Put(r, "Flash", "Task added successfully !")
	}
	// Redirecting to the home page by using http.Redirect
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Get the value "id" from r.FormValue
	value, _ := strconv.Atoi(r.FormValue("id"))

	err := app.todo.Delete(value)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	// Redirecting to the home page by using http.Redirect
	http.Redirect(w, r, "/", http.StatusSeeOther)

	app.Session.Put(r, "Flash", "Task successfully deleted!")
}

func (app *application) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Use the Model object's Get method to retrieve the data for aspecific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	s, err := app.todo.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		app.serverError(w, err)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Write the  data as a plain-text HTTP response body.
	fmt.Fprintf(w, "%v", s)
}

func (app *application) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Fetch the latest tasks from the database to synchronize TaskList
	_, err := app.fetchTasksFromDB()
	if err != nil {
		app.errorLog.Println("Error fetching tasks from DB:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	value, _ := strconv.Atoi(r.FormValue("id"))
	var name string
	var details string
	name = r.FormValue("TaskName")
	details = r.FormValue("Details")
	if len(name) != 0 && len(details) != 0 {
		err := app.todo.UpdateList(value, name, details)
		if err != nil {
			app.errorLog.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
	} else if len(name) == 0 && len(details) == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if len(name) != 0 && len(details) == 0 {
		err := app.todo.UpdateList(value, name, details)
		if err != nil {
			app.errorLog.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
	} else if len(name) == 0 && len(details) != 0 {
		err := app.todo.UpdateList(value, name, details)
		if err != nil {
			app.errorLog.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}

	app.Session.Put(r, "Flash", "Task successfully updated!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Fetch tasks from the database to populate TaskList
func (app *application) fetchTasksFromDB() ([]*models.Todo, error) {
	tasks, err := app.todo.GetAll()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (app *application) validatetask(r *http.Request, str string) bool{
	if strings.TrimSpace(str) ==""{
		app.Session.Put(r, "Flash","One or more fields are empty ")
		return true
	}else if utf8.RuneCountInString(str)>100{
		app.Session.Put(r, "Flash","This name field is too long (maximum is 100 characters")
		return true
	}
	return false
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request){
//fmt.Fprintln(w, "Display the user signup form...")

    files := []string{"./ui/html/signup.page.tmpl" , "./ui/html/base.layout.tmpl"}
    ts, err := template.ParseFiles(files...)

	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server error", 500)
		return
	}
	ts.Execute(w,nil)
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
//fmt.Fprintln(w, "Create a new user...")
  username := r.FormValue("name")
  useremail := r.FormValue("email")
  userpassword := r.FormValue("password")

  err := app.users.Insert(username,useremail,userpassword)
  if err != nil {
	fmt.Println(err)
}
  http.Redirect(w, r, "/", http.StatusSeeOther)
}


func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
//fmt.Fprintln(w, "Display the user login form...")
files := []string{"./ui/html/login.page.tmpl" , "./ui/html/base.layout.tmpl"}
ts, err := template.ParseFiles(files...)

if err != nil {
	app.errorLog.Println(err.Error())
	http.Error(w, "Internal Server error", 500)
	return
}
ts.Execute(w,app.Session.Pop(r,"Flash"))
}


func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
//fmt.Fprintln(w, "Authenticate and login the user...")
useremail := r.FormValue("email")
userpassword := r.FormValue("password")

isUser,err := app.users.Authenticate(useremail,userpassword)
if err != nil {
	app.errorLog.Println(err.Error())
	http.Error(w, "Internal Server error", 500)
}
if isUser{
	app.Session.Put(r,"Authenticated",true)
	app.Session.Put(r,"Flash","Login successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}else{
	app.Session.Put(r,"Flash","Login failed")	
    http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	app.Session.Put(r,"Authenticated",false)
}
}


func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
fmt.Fprintln(w, "Logout the user...")
}
