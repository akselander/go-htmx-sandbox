package auth

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (u *User) nameTaken(r *userRepository) (string, error) {
	if u.Username == "" {
		return "Username is required.", nil
	}

	id, err := r.userExists(u.Username)

	if id != nil {
		return "Username already taken.", nil
	}

	return "", err
}

func (u *User) save(r *userRepository) (*int, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err == bcrypt.ErrPasswordTooLong {
		return nil, "Password is too long.", nil
	}

	u.hash = hashedPassword
	id, err := r.save(u)

	return id, "", err
}

func (u *User) authenticate(r *userRepository) (*int, string) {
	id, hash, err := r.findHashByUsername(u.Username)

	if err != nil {
		fmt.Errorf("Failed authentication attempt for user: %v, reason: %v", u.Username, err)
		return nil, "Credentials did not match"
	}

	if hash == nil {
		fmt.Errorf("Failed authentication attempt for user: %v, reason: no stored hash", u.Username)
		return nil, "Credentials did not match"
	}
	u.hash = hash
	err = bcrypt.CompareHashAndPassword(u.hash, []byte(u.Password))

	if err == bcrypt.ErrHashTooShort {
		fmt.Errorf("Failed authentication attempt for user: %v, reason: %v", u.Username, err)
		return nil, "Credentials did not match"
	}

	if err != nil {
		fmt.Errorf("Failed authentication attempt for user: %v, reason: %v", u.Username, err)
		return nil, "Credentials did not match"
	}

	return id, ""
}

type User struct {
	id       *int
	Username string
	Password string
	hash     []byte
}

func (r *userRepository) userExists(username string) (*int, error) {
	row := r.db.QueryRow("SELECT id FROM users WHERE username = ?", username)

	var id int
	err := row.Scan(id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &id, err
}

func (r *userRepository) findById(id int) (*User, error) {
	row := r.db.QueryRow("SELECT id, username FROM users WHERE id = ?", id)

	var username string
	err := row.Scan(&id, &username)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &User{id: &id, Username: username, Password: "", hash: nil}, nil
}

func (r *userRepository) findHashByUsername(username string) (*int, []byte, error) {
	row := r.db.QueryRow("SELECT id, password FROM users WHERE username = ?", username)

	var id int
	var hash []byte
	err := row.Scan(&id, &hash)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil
	}

	return &id, hash, err
}

func (r *userRepository) save(user *User) (*int, error) {
	res, err := r.db.Exec("INSERT INTO users (username, password) VALUES (?, ?) RETURNING id", user.Username, user.hash)

	if err != nil {
		return nil, err
	}

	int64id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	id := int(int64id)

	return &id, err
}

type userRepository struct {
	db *sql.DB
}
