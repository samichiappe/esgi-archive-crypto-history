package server

import (
	"log"
	"net/http"
	"path/filepath"
)

func downloadCSVHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "Paramètre 'file' manquant", http.StatusBadRequest)
		return
	}
	path, err := filepath.Abs(filename)
	if err != nil {
		http.Error(w, "Fichier invalide", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	http.ServeFile(w, r, path)
}

func StartServer(addr string) {
	http.HandleFunc("/download", downloadCSVHandler)
	log.Printf("Serveur démarré sur %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
