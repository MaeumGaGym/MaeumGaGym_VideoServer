package database

import (
	"database/sql"
	"log"
	_ "zombiezen.com/go/sqlite"
)

var db *sql.DB

func init() {
	var err error
	db, err := sql.Open("sqlite3", "/app/db/jobs.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS videos (
		videoId TEXT PRIMARY KEY,
		scale TEXT,
		duration TEXT,
		url TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func SaveVideo(videoId, scale, duration, url string) error {
	stmt, err := db.Prepare("INSERT INTO videos(videoId, scale, duration, url) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(videoId, scale, duration, url)
	if err != nil {
		return err
	}

	return nil
}

func GetUnfinishedVideos() ([]string, error) {
	rows, err := db.Query("SELECT videoId FROM videos WHERE duration IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videoIds []string
	for rows.Next() {
		var videoId string
		if err := rows.Scan(&videoId); err != nil {
			return nil, err
		}
		videoIds = append(videoIds, videoId)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return videoIds, nil
}
