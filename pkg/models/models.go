package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// Add a new ErrInvalidCredentials error. 
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// Add a new ErrDuplicateEmail error. 
	ErrDuplicateEmail= errors.New("models: duplicatw email")
)


type Todo struct {
	ID      int
	Name    string
	Details string
	Created time.Time
	Expires time.Time
}

// Define a new User type. 

type User struct{
	ID int
	Name string
	Email string
	HashedPassword []byte
	Created time.Time
}