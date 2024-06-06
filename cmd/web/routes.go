package main

import (
	"net/http"
	"github.com/justinas/alice"
	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler{

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders,app.logResponse)

	// Create a new middleware chain containing the middleware specific to
	// our dynamic application routes.
	dynamicMiddleware := alice.New(app.Session.Enable)


	mux := pat.New()
	// Handler functions
	mux.Get("/", dynamicMiddleware.ThenFunc(app.Home))
	mux.Post("/addTask",  dynamicMiddleware.ThenFunc(app.AddTask))
	mux.Get("/getTask",  dynamicMiddleware.ThenFunc(app.GetTask))
	mux.Post("/deleteTask",  dynamicMiddleware.ThenFunc(app.DeleteTask))
	mux.Post("/updateTask",  dynamicMiddleware.ThenFunc(app.UpdateTask))

	// Add the five new routes.
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))


	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))


	return standardMiddleware.Then(mux)

}
