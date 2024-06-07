package mysql

import(
	"database/sql"
	"todo/pkg/models"
)

type UserModel struct{
	DB *sql.DB
}

// We'll use the Insert method to add a new record to the users table.
func (m *UserModel) Insert (name, email, password string)error{
	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err := m.DB.Exec(stmt, name, email, password )
	if err != nil {
		return  err
	}
	return  nil
}

// We'll use the Authenticate method to verify whether a user exists with the provided email address and password.

func (m *UserModel) Authenticate(email, password string)(bool,error){
	stmt:= `select id from users where email =? and hashed_password = ?`
	rows, err := m.DB.Query(stmt, email, password )
	if err != nil {
		return  false,err
	}
	defer rows.Close()
	return rows.Next(),nil
}

// We'll use the Get method to fetch details for a specific user based on their user ID.
func (m *UserModel) Get(id int) (*models.User, error){
	return nil, nil
}