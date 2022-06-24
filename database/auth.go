package database

import (
	"database/sql"
	"errors"
	. "self-hosted-cloud/server/models"
)

func (db *Database) CreateGithubAuthTable() (sql.Result, error) {
	return db.instance.Exec(`
		CREATE TABLE IF NOT EXISTS auth_github (
			username VARCHAR(255) UNIQUE PRIMARY KEY,
			user_id  INTEGER,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)
	`)
}

func (db *Database) GetUserFromGithub(username string) (User, error) {
	statement, err := db.instance.Prepare(`
		SELECT *
		FROM users, auth_github
		WHERE users.id = auth_github.user_id
		  AND auth_github.username = ?;
	`)

	if err != nil {
		return User{}, errors.New("")
	}

	var user User
	err = statement.QueryRow(username).Scan(&user.Id, &user.Username, &user.Name)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *Database) CreateUserFromGithub(githubUser GithubUser) error {
	userId, err := db.CreateUser(User{
		Username: githubUser.Username,
		Name:     githubUser.Name,
	})
	if err != nil {
		return err
	}

	_, err = db.instance.Exec(`INSERT INTO auth_github(username, user_id) VALUES (?, ?)`, githubUser.Username, userId)
	if err != nil {
		return err
	}

	return nil
}
