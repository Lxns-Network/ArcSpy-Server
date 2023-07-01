package database

import (
	"database/sql"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

type SyncUserData struct {
	Player interface{} `json:"player"`
	Cookie string      `json:"cookie"`
	Scores interface{} `json:"scores"`
}

func InsertPlayerData(data *SyncUserData) error {
	db, err := sql.Open("sqlite", "./arcspy.db")
	if err != nil {
		return err
	}
	defer db.Close()

	userId := data.Player.(map[string]interface{})["user_id"]

	userData, _ := json.Marshal(data.Player)
	cookie := data.Cookie
	scores, _ := json.Marshal(data.Scores)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM player WHERE user_id = ?", userId).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		_, err = db.Exec("UPDATE player SET cookie = ?, user_data = ? WHERE user_id = ?", cookie, userData, userId)
	} else {
		_, err = db.Exec("INSERT INTO player (user_id, cookie, user_data) VALUES (?, ?, ?)", userId, cookie, userData)
	}
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO score (user_id, scores) VALUES (?, ?)", userId, string(scores))
	if err != nil {
		return err
	}

	return nil
}

func SelectPlayerData(userId string) (map[string]interface{}, error) {
	db, err := sql.Open("sqlite", "./arcspy.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var userData string
	err = db.QueryRow("SELECT user_data FROM player WHERE user_id = ?", userId).Scan(&userData)
	if err != nil {
		return nil, err
	}

	var playerData map[string]interface{}
	err = json.Unmarshal([]byte(userData), &playerData)
	if err != nil {
		return nil, err
	}

	return playerData, nil
}

func SelectPlayerCookie(userId int) (string, error) {
	db, err := sql.Open("sqlite", "./arcspy.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	var cookie string
	err = db.QueryRow("SELECT cookie FROM player WHERE user_id = ?", userId).Scan(&cookie)
	if err != nil {
		return "", err
	}

	return cookie, nil
}

func SelectPlayerScores(userId string) ([]interface{}, error) {
	db, err := sql.Open("sqlite", "./arcspy.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var scores string
	err = db.QueryRow("SELECT scores FROM score WHERE user_id = ?", userId).Scan(&scores)
	if err != nil {
		return nil, err
	}

	var playerScores []interface{}
	err = json.Unmarshal([]byte(scores), &playerScores)
	if err != nil {
		return nil, err
	}

	return playerScores, nil
}

func Create() {
	db, err := sql.Open("sqlite", "./arcspy.db")
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS player (
		    user_id INTEGER PRIMARY KEY,
		    cookie TEXT,
		    upload_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		    user_data TEXT
		);
		CREATE TABLE IF NOT EXISTS score (
		    user_id INTEGER,
		    scores TEXT,
		    upload_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (user_id) REFERENCES player(user_id)
		);
	`)
	if err != nil {
		log.Fatal("Failed to create database table: ", err)
	}
}
