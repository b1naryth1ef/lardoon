package lardoon

import (
	"fmt"
	"net/http"

	"github.com/alioygur/gores"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type HTTPServer struct {
}

func (h *HTTPServer) getReplay(w http.ResponseWriter, r *http.Request) {
	replayId := chi.URLParam(r, "id")

	row, err := db.Query(`SELECT * FROM replays WHERE id=?`, replayId)
	if err != nil {
		gores.Error(w, 500, fmt.Sprintf("error: %v", err))
	}
	defer row.Close()

	if !row.Next() {
		gores.Error(w, 404, "replay not found")
		return
	}

	var replay ReplayWithObjects
	err = row.Scan(
		&replay.Id,
		&replay.Name,
		&replay.ReferenceTime,
		&replay.RecordingTime,
		&replay.Title,
		&replay.DataSource,
		&replay.DataRecorder,
	)
	if err != nil {
		gores.Error(w, 500, fmt.Sprintf("error: %v", err))
		return
	}

	replay.Objects = make([]*ReplayObject, 0)
	rows, err := db.Query(`SELECT * FROM replay_objects WHERE replay_id = ?`, replay.Id)
	if err != nil {
		gores.Error(w, 500, fmt.Sprintf("error: %v", err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var object ReplayObject
		err := rows.Scan(
			&object.ReplayId,
			&object.Id,
			&object.Name,
			&object.Pilot,
			&object.CreatedOffset,
			&object.DeletedOffset,
		)
		if err != nil {
			gores.Error(w, 500, fmt.Sprintf("error: %v", err))
			return
		}
		replay.Objects = append(replay.Objects, &object)
	}

	gores.JSON(w, 200, replay)
}

func (h *HTTPServer) listReplays(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT * FROM replays`)
	if err != nil {
		gores.Error(w, 500, fmt.Sprintf("error: %v", err))
		return
	}
	defer rows.Close()

	replays := make([]*Replay, 0)
	for rows.Next() {
		var replay Replay
		err := rows.Scan(
			&replay.Id,
			&replay.Name,
			&replay.ReferenceTime,
			&replay.RecordingTime,
			&replay.Title,
			&replay.DataSource,
			&replay.DataRecorder,
		)
		if err != nil {
			gores.Error(w, 500, fmt.Sprintf("error: %v", err))
			return
		}
		replays = append(replays, &replay)
	}

	gores.JSON(w, 200, replays)
}

func (h *HTTPServer) Run(bind string) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/api/replay", h.listReplays)
	r.Get("/api/replay/{id}", h.getReplay)
	// r.Get("/api/replay/download/{name}", h.downloadReplay)
	// r.Post("/api/replay/download/{name}", h.downloadReplay)

	return http.ListenAndServe(bind, r)
}
