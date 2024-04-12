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
	artists, err := fetchArtists()
	if err != nil {
		http.Error(w, "Impossible de récupérer les artistes", http.StatusInternalServerError)
		return
	}

	// Création de FuncMap pour inclure vos fonctions personnalisées
	funcMap := template.FuncMap{
		"formatCreationYear": formatCreationYear,
		"formatDateToFrench": formatDateToFrench,
	}

	// Chargement du template avec FuncMap
	t, err := template.New("artists.html").Funcs(funcMap).ParseFiles("templates/artists.html")
	if err != nil {
		log.Printf("Erreur lors du chargement du template: %v", err)
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}

	// Exécution du template avec les données (artists)
	err = t.Execute(w, artists)
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

	// Convertir les emplacements en JSON pour utilisation dans JavaScript
	artistLocationsJson, err := json.Marshal(artistDetail.Locations)
	if err != nil {
		log.Printf("Error serializing locations: %v", err)
		http.Error(w, "Failed to serialize locations", http.StatusInternalServerError)
		return
	}

	// Préparer les fonctions à utiliser dans le template
	funcMap := template.FuncMap{
		"js": func(js string) template.JS {
			return template.JS(js)
		},
		"formatCreationYear": formatCreationYear,
		"formatDateToFrench": formatDateToFrench,
	}

	// Charger le template avec les fonctions enregistrées
	t, err := template.New("artistDetails.html").Funcs(funcMap).ParseFiles("templates/artistDetails.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	// Créer une map pour les données à passer au template
	data := map[string]interface{}{
		"Artist":              artistDetail.Artist,
		"FormattedFirstAlbum": artistDetail.Artist.FirstAlbumDateStr,
		"Locations":           artistDetail.Locations,
		"Dates":               artistDetail.Dates,
		"MapLinks":            artistDetail.MapLinks,
		"ArtistLocationsJson": template.JS(artistLocationsJson),
	}

	// Exécuter le template avec les données
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
	location = strings.Replace(location, "_", " ", -1)
	parts := strings.Split(location, "-")
	for i, part := range parts {
		if i == 0 { // Assuming the first part is the city
			words := strings.Fields(part)
			for j, word := range words {
				words[j] = strings.Title(strings.ToLower(word))
			}
			parts[i] = strings.Join(words, " ")
		} else { // Assuming the second part is the country
			parts[i] = strings.ToUpper(strings.TrimSpace(part))
		}
	}
	return strings.Join(parts, " - ")
}
