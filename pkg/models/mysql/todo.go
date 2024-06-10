package mysql

import (
	"database/sql"
	"fmt"
	"todo/pkg/models"
)

// Define a TodoModel type which wraps a sql.DB connection pool.
type TodoModel struct {
	DB *sql.DB
}

type SpecialModel struct {
	DB *sql.DB
}

// This will insert a new task into the database.
func (m *TodoModel) Insert(name, details, expires string) (int, error) {

	stmt := `INSERT INTO todo (name, details, created, expires)
	VALUES(?, ?,  UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, name, details, expires)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil

}

func (m *SpecialModel) Insertspecial(name, details, expires string) (int, error) {

	stmt := `INSERT INTO specialtask (name, details, created, expires)
	VALUES(?, ?,  UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, name, details, expires)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
	
}

// This will return a specific task based on its id.
func (m *TodoModel) Get(id int) (*models.Todo, error) {

	stmt := `SELECT id, name, details, created, expires FROM todo
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Todo{}

	err := row.Scan(&s.ID, &s.Name, &s.Details, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

func (m *SpecialModel) Getspecial() ([]*models.Todo, error) {
	stmt := ` SELECT id, name, details, created, expires FROM specialtask  WHERE expires > UTC_TIMESTAMP() and
	  name LIKE 'special%' ORDER BY created DESC LIMIt 10; `

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	specials := []*models.Todo{}

	for rows.Next() {
		s := &models.Todo{}
		err = rows.Scan(&s.ID, &s.Name, &s.Details, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		specials = append(specials, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return specials, nil

}

func (m *TodoModel) GetAll() ([]*models.Todo, error) {
	stmt := `SELECT id, name, details, created, expires FROM todo
    WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []*models.Todo{}

	for rows.Next() {
		s := &models.Todo{}
		err = rows.Scan(&s.ID, &s.Name, &s.Details, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		todos = append(todos, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil

}

func (m *TodoModel) Delete(name string) error {
	stmt := `DELETE FROM todo where name = ?`

	_, err := m.DB.Exec(stmt, name)
	if err != nil {
		return err
	}
	return nil
}

func (m *SpecialModel) Deletespecial(name string) error {
	stmt := `DELETE FROM specialtask where name = ?`

	_, err := m.DB.Exec(stmt, name)
	if err != nil {
		return err
	}
	return nil
}

func (m *TodoModel) UpdateList(id int, name string, details string) error {
	stmt := `update todo set name = ? , details = ?  where id = ?`

	_, err := m.DB.Exec(stmt, name, details, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *SpecialModel) UpdatespecialList(id int, name string, details string) error {
	stmt := `update specialtask set name = ? , details = ?  where id = ?`

	_, err := m.DB.Exec(stmt, name, details, id)
	if err != nil {
		return err
	}
	return nil
}

