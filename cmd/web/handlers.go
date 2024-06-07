package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	//"strings"
	"todo/pkg/models"
	//"unicode/utf8"
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

	message := app.Session.PopString(r, "Flash")

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

func (app *application) SpecialTask(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "Display the special task...")
	s, err := app.Specialtask.Getspecial()
	if err != nil {
		app.serverError(w, err)
		return
	}

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

	// Executing the template and checking for any error
	err = ts.Execute(w, struct {
		Tasks []*models.Todo
		Flash string
	}{
		Tasks: s,
		Flash: "",
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

	_, err := app.todo.Insert(name, details, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	isSpec := strings.Contains(name, "special:")

	if isSpec {
		_, err := app.Specialtask.Insertspecial(name, details, expires)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	// Redirecting to the home page by using http.Redirect
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Get the value "id" from r.FormValue
	name := r.FormValue("TaskName")
	
	err := app.todo.Delete(name)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}
 
	{
	err := app.Specialtask.Deletespecial(name)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}
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

	// Use the Model object's Get method to retrieve the data for aspecific record based on its ID. 

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

	// Parse the task ID from the form value
	value, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		app.errorLog.Println("Invalid task ID:", err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get task name and details from the form
	name := r.FormValue("TaskName")
	details := r.FormValue("Details")

	// Log the received values for debugging
	app.InfoLog.Printf("Received update request for Task ID %d with name: %s and details: %s", value, name, details)

	// Check if both name and details are provided, or handle cases where one or both are missing
	if len(name) == 0 && len(details) == 0 {
		app.Session.Put(r, "Flash", "Empty fields cannot be  updated!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Update the task in the todo list
	err = app.todo.UpdateList(value, name, details)
	if err != nil {
		app.errorLog.Println("Error updating task:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
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

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request){

	files := []string{"./ui/html/signup.page.tmpl", "./ui/html/base.layout.tmpl"}
	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server error", 500)
		return
	}
	ts.Execute(w, nil)
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("name")
	useremail := r.FormValue("email")
	userpassword := r.FormValue("password")

	err := app.users.Insert(username, useremail, userpassword)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	files := []string{"./ui/html/login.page.tmpl", "./ui/html/base.layout.tmpl"}
	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server error", 500)
		return
	}
	ts.Execute(w, app.Session.Pop(r, "Flash"))
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	useremail := r.FormValue("email")
	userpassword := r.FormValue("password")

	isUser, err := app.users.Authenticate(useremail, userpassword)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server error", 500)
	}
	if isUser {
		app.Session.Put(r, "Authenticated", true)
		app.Session.Put(r, "Flash", "Login successfully")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		app.Session.Put(r, "Flash", "Login failed")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		app.Session.Put(r, "Authenticated", false)
	}
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.Session.Put(r, "Authenticated", false)
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
