package mysql

import (
	"alexedwards.net/snippetbox/pkg/models"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

type DepartmentsModel struct {
	DB *sql.DB
}

func (m *DepartmentsModel) Insert(title, content, category string) (int, error) {
	stmt := `INSERT INTO departments (title, content, date, category) VALUES($1, $2, CURRENT_TIMESTAMP, $3) RETURNING id`
	var id int
	err := m.DB.QueryRow(stmt, title, content, category).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *DepartmentsModel) Get(id int) (*models.Departments, error) {
	stmt := `SELECT id, title, content, date, category FROM departments WHERE id = $1`
	row := m.DB.QueryRow(stmt, id)
	n := &models.Departments{}
	err := row.Scan(&n.ID, &n.Title, &n.Content, &n.Date, &n.Category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return n, nil
}

func (m *DepartmentsModel) Latest() ([]*models.Departments, error) {
	stmt := `SELECT id, title, content, date, category FROM departments ORDER BY date DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	departmentsList := []*models.Departments{}
	for rows.Next() {
		n := &models.Departments{}
		err = rows.Scan(&n.ID, &n.Title, &n.Content, &n.Date, &n.Category)
		if err != nil {
			return nil, err
		}
		departmentsList = append(departmentsList, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return departmentsList, nil
}

func (m *DepartmentsModel) GetByCategory(category string) ([]*models.Departments, error) {
	stmt := `SELECT id, title, content, category, date FROM departments WHERE category = $1 ORDER BY date DESC`
	rows, err := m.DB.Query(stmt, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	DepartmentsArray := make([]*models.Departments, 0)
	for rows.Next() {
		n := &models.Departments{}
		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Category, &n.Date)
		if err != nil {
			return nil, err
		}
		DepartmentsArray = append(DepartmentsArray, n)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return DepartmentsArray, nil
}
