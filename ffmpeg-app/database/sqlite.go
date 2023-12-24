package database

import (
	"log"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var conn *sqlite.Conn

func init() {
	var err error
	conn, err = sqlite.OpenConn("/app/db/jobs.db", sqlite.OpenReadWrite)
	if err != nil {
		log.Fatal(err)
	}

	err = sqlitex.ExecScript(conn, `
	CREATE TABLE IF NOT EXISTS videos (
		videoId TEXT PRIMARY KEY,
		url TEXT
	);
	CREATE TABLE IF NOT EXISTS video_resolutions (
		videoId TEXT,
		scale TEXT,
		duration TEXT,
		FOREIGN KEY(videoId) REFERENCES videos(videoId)
	);
	`)
	if err != nil {
		log.Fatal(err)
	}
}
func SaveVideo(videoId, scale, duration, url string) error {
	opts := &sqlitex.ExecOptions{
		Args: []any{videoId, url},
	}
	err := sqlitex.Execute(conn, "INSERT OR REPLACE INTO videos(videoId, url) VALUES (?, ?)", opts)
	if err != nil {
		return err
	}

	opts = &sqlitex.ExecOptions{
		Args: []any{videoId, scale, duration},
	}
	err = sqlitex.Execute(conn, "INSERT OR REPLACE INTO video_resolutions(videoId, scale, duration) VALUES (?, ?, ?)", opts)
	return err
}

func GetUnfinishedVideos() ([]string, error) {
	var videoIds []string
	opts := &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			videoId := stmt.ColumnText(0)
			videoIds = append(videoIds, videoId)
			return nil
		},
	}
	err := sqlitex.Execute(conn, "SELECT videoId FROM video_resolutions WHERE duration IS NULL", opts)
	if err != nil {
		return nil, err
	}
	return videoIds, nil
}
