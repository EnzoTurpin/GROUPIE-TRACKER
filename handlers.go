package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

// Lors de la définition de votre handler pour les artistes
func artistsHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")

	artists, err := fetchArtists()
	if err != nil {
		http.Error(w, "Impossible de récupérer les artistes", http.StatusInternalServerError)
		return
	}

	// Filtrer les artistes si searchQuery n'est pas vide
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
		"Artists":     artists, // Assurez-vous que ceci est une tranche de Artist
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

func formatDateToFrench(dateStr string) string {
	if dateStr == "" {
		return "Date non spécifiée" // Return this if the date string is empty
	}
	parsedDate, err := time.Parse("02-01-2006", dateStr) // Ensure this format matches your data
	if err != nil {
		log.Printf("Failed to parse date '%s': %v", dateStr, err)
		return "Format de date invalide"
	}
	frenchDate := fmt.Sprintf("%02d %s %d", parsedDate.Day(), monthsFrench[parsedDate.Month().String()], parsedDate.Year())
	return frenchDate
}

func formatCreationYear(year int) string {
	return strconv.Itoa(year) // Convertit une année de type int en string
}

func formatLocationName(location string) string {
	// Replace underscores with spaces to handle city names with spaces
	location = strings.Replace(location, "_", " ", -1)

	// Split the location into parts (city and country) assuming '-' is the delimiter
	parts := strings.Split(location, "-")
	if len(parts) == 2 {
		// Capitalize the first letter of each word in the city name
		parts[0] = strings.Title(parts[0])
		// Convert the country code/name to uppercase
		parts[1] = strings.ToUpper(parts[1])
	}

	// Join the parts back into a single string
	return strings.Join(parts, ", ")
}
