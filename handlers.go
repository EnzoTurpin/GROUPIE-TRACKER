package main

import (
	"html/template"
	"net/http"
)

func artistsHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := fetchArtists()
	if err != nil {
		http.Error(w, "Impossible de récupérer les artistes", http.StatusInternalServerError)
		return
	}

	// Utilisation de templates pour afficher les artistes
	t, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, artists)
}
