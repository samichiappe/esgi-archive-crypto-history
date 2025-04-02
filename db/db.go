package db

import (
	"database/sql"
	"log"
	"time"

	"esgi-archive-crypto-history/kraken"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS asset_pairs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		pair_key TEXT,
		altname TEXT,
		wsname TEXT,
		timestamp DATETIME
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Println("Erreur lors de la création de la table:", err)
	}
	return err
}

func InsertAssetPairs(db *sql.DB, pairs map[string]kraken.AssetPair, timestamp time.Time) error {
	stmt, err := db.Prepare(`INSERT INTO asset_pairs (pair_key, altname, wsname, timestamp) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for key, pair := range pairs {
		_, err := stmt.Exec(key, pair.Altname, pair.Wsname, timestamp)
		if err != nil {
			return err
		}
	}
	return nil
}
