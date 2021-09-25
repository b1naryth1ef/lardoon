package lardoon

import (
	"database/sql"
	"os"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS replays (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT UNIQUE,
	reference_time TEXT,
	recording_time TEXT,
	title TEXT,
	data_source TEXT,
	data_recorder TEXT
);

CREATE TABLE IF NOT EXISTS replay_objects (
	replay_id INTEGER,
	object_id INTEGER,
	name TEXT,
	pilot TEXT,
	created_offset INTEGER,
	deleted_offset INTEGER,

	UNIQUE(replay_id, object_id),
	FOREIGN KEY(replay_id) REFERENCES replays(id)
);
`

type Replay struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	ReferenceTime string `json:"reference_time"`
	RecordingTime string `json:"recording_time"`
	Title         string `json:"title"`
	DataSource    string `json:"data_source"`
	DataRecorder  string `json:"data_recorder"`
}

type ReplayObject struct {
	Id            int    `json:"id"`
	ReplayId      int    `json:"replay_id"`
	Name          string `json:"name"`
	Pilot         string `json:"pilot"`
	CreatedOffset int    `json:"created_offset"`
	DeletedOffset int    `json:"deleted_offset"`
}

type ReplayWithObjects struct {
	Replay

	Objects []*ReplayObject `json:"objects"`
}

func createReplayObject(replayId int, objectId int, name string, pilot string, createdOffset int, deletedOffset int) error {
	_, err := db.Exec(`
		INSERT INTO replay_objects VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (replay_id, object_id) DO UPDATE
			SET name = EXCLUDED.name, pilot = EXCLUDED.pilot, created_offset = EXCLUDED.created_offset, deleted_offset = EXCLUDED.deleted_offset
	`, replayId, objectId, name, pilot, createdOffset, deletedOffset)
	return err
}

func createReplay(name string, referenceTime string, recordingTime string, title string, dataSource string, dataRecorder string, force bool) (int, error) {
	var row *sql.Rows
	var err error
	if force {
		row, err = db.Query(`
		INSERT INTO replays (name, reference_time, recording_time, title, data_source, data_recorder) VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (replays.name) DO UPDATE
			SET reference_time = EXCLUDED.reference_time, recording_time = EXCLUDED.recording_time, title = EXCLUDED.title,
			data_source = EXCLUDED.data_source, data_recorder = EXCLUDED.data_recorder
		RETURNING id
	`, name, referenceTime, recordingTime, title, dataSource, dataRecorder)
	} else {
		row, err = db.Query(`
		INSERT INTO replays (name, reference_time, recording_time, title, data_source, data_recorder) VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id
	`, name, referenceTime, recordingTime, title, dataSource, dataRecorder)
	}
	if err != nil {
		return -1, err
	}
	defer row.Close()

	row.Next()

	var id int
	err = row.Scan(&id)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return -1, nil
			}
		}
		return -1, err
	}

	return id, nil
}

func init() {
	dbPath := os.Getenv("LARDOON_DB_PATH")
	if dbPath == "" {
		dbPath = "./lardoon.db"
	}

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(schema)
	if err != nil {
		panic(err)
	}
}
