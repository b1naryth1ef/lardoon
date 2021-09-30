package lardoon

import (
	"log"
	"os"
)

func PruneReplays(dryRun bool) error {
	rows, err := db.Query(`SELECT id, path FROM replays`)
	if err != nil {
		return err
	}
	defer rows.Close()

	toPrune := []int{}
	for rows.Next() {
		var id int
		var path string
		err = rows.Scan(&id, &path)
		if err != nil {
			return err
		}

		_, err := os.Stat(path)
		if err != nil && os.IsNotExist(err) {
			toPrune = append(toPrune, id)
		}
	}

	if len(toPrune) > 0 {

		if !dryRun {
			log.Printf("Pruning replays: %v", toPrune)
			stmt, err := db.Prepare("DELETE FROM replays WHERE id = ?")
			if err != nil {
				return err
			}
			defer stmt.Close()

			for _, replayId := range toPrune {
				_, err = stmt.Exec(replayId)
				if err != nil {
					return err
				}
			}
		} else {
			log.Printf("[DRY-RUN] Would prune replays: %v", toPrune)
		}
	} else {
		log.Printf("No replays to prune")
	}

	return nil

}
