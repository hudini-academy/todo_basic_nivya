package main

import (
	"database/sql"
	"log"
	"net/http"
	"todo/pkg/models/mysql"
	"flag"
	"time" 
	"github.com/golangcollege/sessions"
	_ "github.com/go-sql-driver/mysql"
)

//  Define an application struct to hold the application-wide dependencies

type application struct {
	errorLog *log.Logger
	InfoLog  *log.Logger
	Session *sessions.Session
	todo     *mysql.TodoModel
	users *mysql.UserModel
	
}


func main() {
	addr := ":4000"
	infoLog, errorLog := initializeLogFiles()
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret")

	// To keep the main() function tidy I've put the code for creating a connec
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.

	db, err := openDB("root:root@/todo?parseTime=true")
	if err != nil {
		errorLog.Fatal(err)
	}
	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour


	// Initialize a new instance of application containing the dependencies.
	
	app := &application{
		errorLog: errorLog,
		InfoLog:  infoLog,
		Session: session,
		todo:     &mysql.TodoModel{DB: db},
		users: &mysql.UserModel{DB: db},
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

