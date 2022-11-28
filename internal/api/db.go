package api

import (
	"errors"
	"wishmill/internal/db"
	"wishmill/internal/logger"
)

// Will get user from db
func GetUser(sub string, authProvider string) (*User, error) { //TODO needed?
	rows, err := db.Db.Query("SELECT id,name,auth_provider,sub,email FROM users WHERE sub=$1 AND auth_provider=$2", sub, authProvider)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	defer rows.Close()
	u := User{}
	next := rows.Next()
	if !next {
		err = errors.New("user does not exist")
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	err = rows.Scan(&u.Id, &u.Name, &u.AuthProvider, &u.Sub, &u.Email)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	return &u, nil
}

// Return true if a user is already registered
func getUserRegistration(sub string, authProvider string) (bool, error) {

	rows, err := db.Db.Query("SELECT * FROM users WHERE sub=$1 AND auth_provider=$2", sub, authProvider)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func updateUser(sub string, authProvider string, name string, email string) error {
	_, err := db.Db.Exec("UPDATE users SET name=$1, email=$2 WHERE sub=$3 AND auth_provider=$4", name, email, sub, authProvider)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	return nil
}

// Register a user
func registerUser(sub string, authProvider string, name string, email string) (int64, error) {

	//Check if user is already registered
	reg, err := getUserRegistration(sub, authProvider)
	if err != nil {
		return 0, err
	}

	if reg {
		err = updateUser(sub, authProvider, name, email)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}
		return 0, err
	}

	//Register User
	var id int64
	err = db.Db.QueryRow("INSERT INTO users(name, auth_provider, sub, email) VALUES ($1, $2, $3, $4) RETURNING id", name, authProvider, sub, email).Scan(&id)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return 0, err
	}

	//Return registered user
	return id, nil
}
