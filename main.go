package main

import (
	"esgi-archive-crypto-history/csvwriter"
	"esgi-archive-crypto-history/db"
	"esgi-archive-crypto-history/kraken"
	"esgi-archive-crypto-history/server"
	"log"
	"time"
)

func main() {
	database, err := db.InitDB("data.db")
	if err != nil {
		log.Fatal("Erreur lors de l'initialisation de la DB:", err)
	}
	defer database.Close()

	if err := db.CreateTables(database); err != nil {
		log.Fatal("Erreur lors de la création des tables:", err)
	}

	go server.StartServer(":8080")

	for {
		now := time.Now()
		log.Printf("Récupération des données à %s\n", now.Format("15:04:05"))

		data, err := kraken.FetchAllDataConcurrently()
		if err != nil {
			log.Println("Erreur lors de la récupération des données:", err)
		} else {
			log.Printf("Statut du serveur Kraken: %s (timestamp: %s)\n", data.Status.Status, data.Status.Timestamp)

			filename := csvwriter.GetCSVFileName(now)
			if err := csvwriter.WriteAssetPairsCSV(filename, data.Pairs); err != nil {
				log.Println("Erreur lors de l'écriture du CSV:", err)
			} else {
				log.Printf("Les données ont été écrites dans %s\n", filename)
			}

			if err := db.InsertAssetPairs(database, data.Pairs, now); err != nil {
				log.Println("Erreur lors de l'insertion dans la DB:", err)
			}
		}

		time.Sleep(5 * time.Minute)
	}
}
