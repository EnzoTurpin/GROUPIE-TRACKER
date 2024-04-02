package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func artistsHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := fetchArtists()
	if err != nil {
		http.Error(w, "Impossible de récupérer les artistes", http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, artists)
}

func artistDetailHandler(w http.ResponseWriter, r *http.Request) {
	artistIDStr := strings.TrimPrefix(r.URL.Path, "/artist/")
	artistID, err := strconv.Atoi(artistIDStr)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	artistDetail, err := fetchArtistDetails(artistID)
	if err != nil {
		http.Error(w, "Failed to fetch artist details", http.StatusInternalServerError)
		return
	}

	// Convertir les emplacements en JSON
	artistLocationsJson, err := json.Marshal(artistDetail.Locations)
	if err != nil {
		http.Error(w, "Failed to serialize locations", http.StatusInternalServerError)
		return
	}

	// Log des valeurs pour débogage
	log.Printf("Locations: %v", artistDetail.Locations)
	log.Printf("MapLinks: %v", artistDetail.MapLinks)

	// Charger le template
	t, err := template.New("artistDetails.html").Funcs(template.FuncMap{"js": func(js string) template.JS {
		return template.JS(js)
	}}).ParseFiles("templates/artistDetails.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	// Créer une map pour les données à passer au template
	data := map[string]interface{}{
		"Artist":              artistDetail.Artist,
		"ArtistLocationsJson": template.JS(artistLocationsJson),
	}

	// Exécuter le template avec les données
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
	}
}
