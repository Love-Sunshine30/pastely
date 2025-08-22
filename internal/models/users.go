package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql" // New import
	"golang.org/x/crypto/bcrypt"     // New import
)

// the user type
type User struct {
	ID             int
	Name           string
	Email          string
	HashedParrword []byte
	Created        time.Time
}

// Define a user model that wraps around a database connection pool
type UserModel struct {
	DB *sql.DB
}

// Insert will add a new user to our users table
func (m *UserModel) Insert(name, email, password string) error {
	//create a hash of the plain-text password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	// mysql query
	stm := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, NOW())`

	_, err = m.DB.Exec(stm, name, email, string(hashedPassword))
	if err != nil {
		// we check what type of error happend.
		// if the error is duplicated email then we return ErrDuplicateEmail error

		var MySQLError *mysql.MySQLError
		if errors.As(err, &MySQLError) {
			if MySQLError.Number == 1062 && strings.Contains(MySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// Authenticate will check whether a user exist with provided email and password in our database
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stm := "SELECT id, hashed_password FROM users WHERE email=?"

	// Retrive the user's id and password
	err := m.DB.QueryRow(stm, email).Scan(&id, &hashedPassword)
	if err != nil {
		// check if the record exist or not
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredential
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredential
		} else {
			return 0, err
		}
	}
	return id, nil
}

// Exist will check if a user exist with a specific id
func (m *UserModel) Exists(id int) (bool, error) {
	return true, nil
}
