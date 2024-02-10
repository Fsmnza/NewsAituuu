package mysql

import (
	"alexedwards.net/snippetbox/pkg/models"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `
        INSERT INTO users (name, email, hashed_password, created) VALUES($1, $2, $3,CURRENT_TIMESTAMP )`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) {
			if pgError.Code == "23505" && strings.Contains(pgError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = $1 "
	row := m.DB.QueryRow(stmt, email)

	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

//	func (m *UserModel) GetRole(id int) (string, error) {
//		var role string
//		stmt := "SELECT role FROM users WHERE id = $1"
//		row := m.DB.QueryRow(stmt, id)
//		err := row.Scan(&role)
//		if err != nil {
//			if errors.Is(err, sql.ErrNoRows) {
//				return "", models.ErrUserNotFound
//			} else {
//				return "", err
//			}
//		}
//		return role, nil
//	}
//
//	func (m *UserModel) IsAdmin(id int) (bool, error) {
//		var role string
//		stmt := "SELECT role FROM users WHERE id = $1"
//		row := m.DB.QueryRow(stmt, id)
//		err := row.Scan(&role)
//		if err != nil {
//			if errors.Is(err, sql.ErrNoRows) {
//				return false, models.ErrUserNotFound
//			} else {
//				return false, err
//			}
//		}
//		return role == models.RoleAdmin, nil
//	}
//
//	func (m *UserModel) IsTeacher(id int) (bool, error) {
//		var role string
//		stmt := "SELECT role FROM users WHERE id = $1"
//		row := m.DB.QueryRow(stmt, id)
//		err := row.Scan(&role)
//		if err != nil {
//			if errors.Is(err, sql.ErrNoRows) {
//				return false, models.ErrUserNotFound
//			} else {
//				return false, err
//			}
//		}
//		return role == models.RoleTeacher, nil
//	}
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
