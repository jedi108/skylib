package utils

import (
	"skylib/app"
	"log"
	"github.com/nu7hatch/gouuid"
)

func generateToken() string {

	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	return u.String()
}

func StoreSession(userId int) (string, error) {

	var token = generateToken()

	app.GetConnection()
	_, err := app.DB.Exec(
		"INSERT INTO session "+
			"(Session, UserId)"+
			"VALUES (?,?);",
		string(token),
		userId,
	)

	if err != nil {
		log.Println("Session Cretion failed: ", err)
	}
	return token, err
}

func GetSession(token string) (int, error) {

	var userId int
	app.GetConnection()
	err := app.DB.QueryRow(
		"SELECT UserId FROM session WHERE Session=?",
		token).Scan(&userId)

	if err != nil {
		log.Println("Get Session failed: ", err)
	}

	return userId, err
}