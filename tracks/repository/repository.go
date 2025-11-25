package repository

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	DB *sql.DB
}

var repo Repository

func Init() {
	if db, err := sql.Open("sqlite3", "./tmp/test.db"); err == nil {
		repo = Repository{DB: db}
	} else {
		log.Fatal("Database initialisation")
	}
}

func Create() int {
	const sql = "CREATE TABLE IF NOT EXISTS Tracks" +
		"(Id TEXT PRIMARY KEY, Audio TEXT)"
	if _, err := repo.DB.Exec(sql); err == nil {
		return 0
	} else {
		return -1
	}
}

func Clear() int {
	const sql = "DELETE FROM Tracks"
	if _, err := repo.DB.Exec(sql); err == nil {
		return 0
	} else {
		return -1
	}
}

func Update(t Track) int64 {
	const sql = " UPDATE Tracks SET Audio = ? WHERE Id = ?"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(t.Audio, t.Id); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n
			}
		}
	}
	return -1
}

func Insert(t Track) int64 {
	const sql = "INSERT INTO Tracks (Id, Audio) VALUES (? , ?)"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(t.Id, t.Audio); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n
			}
		}
	}

	return -1
}

func List() ([]string, int64) {
	const sql = "SELECT * FROM Tracks"
	if rows, err := repo.DB.Query(sql); err == nil {
		defer rows.Close()
		var trackIDs []string
		for rows.Next() {
			var t Track
			if err := rows.Scan(&t.Id, &t.Audio); err == nil {
				trackIDs = append(trackIDs, t.Id)
			}
		}
		return trackIDs, 1
	}
	return nil, -1
}

func Read(Id string) (Track, int64) {
	const sql = " SELECT * FROM Tracks WHERE Id = ?"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		var t Track
		row := stmt.QueryRow(Id)
		if err := row.Scan(&t.Id, &t.Audio); err == nil {
			return t, 1
		} else {
			return Track{}, 0
		}
	}
	return Track{}, -1
}

func Delete(Id string) int64 {
	const sql = "DELETE FROM Tracks WHERE Id = ?"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(Id); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n
			}
		}
	}
	return -1
}
