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
	searchQuery := r.URL.Query().Get("search")

	artists, err := fetchArtists()
	if err != nil {
		http.Error(w, "Impossible de récupérer les artistes", http.StatusInternalServerError)
		return
	}

	if searchQuery != "" {
		var filteredArtists []Artist
		for _, artist := range artists {
			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchQuery)) {
				filteredArtists = append(filteredArtists, artist)
			}
		}
		artists = filteredArtists
	}

	funcMap := template.FuncMap{
		"formatCreationYear": formatCreationYear,
		"formatDateToFrench": formatDateToFrench,
	}

	t, err := template.New("artists.html").Funcs(funcMap).ParseFiles("templates/artists.html")
	if err != nil {
		log.Printf("Erreur lors du chargement du template: %v", err)
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Artists":     artists,
		"SearchQuery": searchQuery,
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template: %v", err)
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
	}
}

func artistDetailHandler(w http.ResponseWriter, r *http.Request) {
	artistIDStr := strings.TrimPrefix(r.URL.Path, "/artist/")
	artistID, err := strconv.Atoi(artistIDStr)
	if err != nil {
		log.Printf("Error converting artistID to int: %v", err)
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	artistDetail, err := fetchArtistDetails(artistID)
	if err != nil {
		log.Printf("Error fetching artist details for ID %d: %v", artistID, err)
		http.Error(w, "Failed to fetch artist details", http.StatusInternalServerError)
		return
	}

	artistLocationsDatesJson, err := json.Marshal(artistDetail.DatesLocations)
	if err != nil {
		log.Printf("Error serializing artist locations: %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	firstLocationCoordsJson, err := json.Marshal(artistDetail.FirstLocationCoords)
	if err != nil {
		log.Printf("Error serializing coordinates: %v", err)
		http.Error(w, "Failed to serialize coordinates", http.StatusInternalServerError)
		return
	}

	funcMap := template.FuncMap{
		"formatCreationYear": formatCreationYear,
		"formatDateToFrench": formatDateToFrench,
	}

	t, err := template.New("artistDetails.html").Funcs(funcMap).ParseFiles("templates/artistDetails.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Artist":                   artistDetail.Artist,
		"FormattedFirstAlbum":      artistDetail.Artist.FirstAlbumDateStr,
		"DatesLocations":           artistDetail.DatesLocations,
		"MapLinks":                 artistDetail.MapLinks,
		"ArtistLocationsDatesJson": template.JS(artistLocationsDatesJson),
		"FirstLocationCoords":      template.JS(firstLocationCoordsJson),
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template for artist ID %d: %v", artistID, err)
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
	}
}
