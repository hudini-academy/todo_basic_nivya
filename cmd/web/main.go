package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"
	"todo/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

//  Define an application struct to hold the application-wide dependencies

type application struct {
	errorLog    *log.Logger
	InfoLog     *log.Logger
	Session     *sessions.Session
	todo        *mysql.TodoModel
	users       *mysql.UserModel
	Specialtask *mysql.SpecialModel
}

func main() {
	addr := ":4000"
	infoLog, errorLog := initializeLogFiles()
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret")

	// To keep the main() function tidy I've put the code for creating a connection

	db, err := openDB("root:root@/todo?parseTime=true")
	if err != nil {
		errorLog.Fatal(err)
	}
	//  defer a call to db.Close(), so that the connection pool is closed.
	defer db.Close()

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	// Initialize a new instance of application containing the dependencies.

	app := &application{
		errorLog: errorLog,
		InfoLog:  infoLog,
		Session:  session,
		todo:     &mysql.TodoModel{DB: db},
		users:    &mysql.UserModel{DB: db},
		Specialtask: &mysql.SpecialModel{DB: db},
	}

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("starting server on %s", addr)

	// Call the ListenAndServe() method on our new http.Server struct.
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// func humanDate(t time.Time) string {
// 	return t.UTC().Format("02 Jan 2006 at 15:04")
// 	}
	
