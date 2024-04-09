package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ArtistDetail struct {
	Artist    Artist
	Locations []string
	MapLinks  []string // Liens Google Maps pour chaque location
	Dates     []string
	Relation  Relation
}

type Artist struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Members []string `json:"members"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID        int    `json:"id"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
}

var monthsFrench = map[string]string{
	"January":   "Janvier",
	"February":  "Février",
	"March":     "Mars",
	"April":     "Avril",
	"May":       "Mai",
	"June":      "Juin",
	"July":      "Juillet",
	"August":    "Août",
	"September": "Septembre",
	"October":   "Octobre",
	"November":  "Novembre",
	"December":  "Décembre",
}

func fetchArtists() ([]Artist, error) {
	var artists []Artist
	if err := fetchAPI("https://groupietrackers.herokuapp.com/api/artists", &artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func fetchArtistDetails(artistID int) (ArtistDetail, error) {
	var detail ArtistDetail

	// Fetch artist information
	artist, err := fetchArtists()
	if err != nil {
		return detail, err
	}
	found := false
	for _, a := range artist {
		if a.ID == artistID {
			detail.Artist = a
			found = true
			break
		}
	}
	if !found {
		return detail, fmt.Errorf("artist with ID %d not found", artistID)
	}

	// Fetch location
	locationURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%d", artistID)
	var location Location
	if err := fetchAPI(locationURL, &location); err != nil {
		return detail, err
	}
	detail.Locations = location.Locations

	for i, loc := range detail.Locations {
		detail.Locations[i] = formatLocation(loc)
	}

	// Générez les liens Google Maps pour chaque lieu
	detail.MapLinks = make([]string, len(detail.Locations))
	for i, locationName := range detail.Locations {
		detail.MapLinks[i] = generateGoogleMapsLink(locationName)
	}

	// Fetch dates
	datesURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%d", artistID)
	var date Date
	if err := fetchAPI(datesURL, &date); err != nil {
		return detail, err
	}
	detail.Dates = date.Dates

	for i, date := range detail.Dates {
		detail.Dates[i] = formatDate(date)
	}

	// Fetch relation
	relationURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%d", artistID)
	var relation Relation
	if err := fetchAPI(relationURL, &relation); err != nil {
		return detail, err
	}
	detail.Relation = relation

	return detail, nil
}

func fetchAPI(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &target)
}

func generateGoogleMapsLink(locationName string) string {
	baseUrl := "https://www.google.com/maps/search/?api=1&query="
	return baseUrl + url.QueryEscape(locationName)
}

func formatLocation(location string) string {
	parts := strings.Split(location, "-")
	for i, part := range parts {
		part = strings.ReplaceAll(part, "_", " ")
		parts[i] = strings.Title(strings.ToLower(part)) // Ensure consistent capitalization
	}
	return strings.Join(parts, " - ")
}

func formatDate(dateStr string) string {
	// Supprimer l'astérisque (*) du début de la chaîne de date si présent
	cleanDateStr := strings.TrimPrefix(dateStr, "*")

	// Parser la date sans l'astérisque, en utilisant le format attendu après nettoyage
	date, err := time.Parse("02-01-2006", cleanDateStr)
	if err != nil {
		log.Printf("Failed to parse date: %v", err)
		return dateStr // Retourner la chaîne originale en cas d'erreur
	}

	// Traduire le mois en français
	month := monthsFrench[date.Month().String()]

	// Formater la date en "23 Août 2019"
	return fmt.Sprintf("%02d %s %d", date.Day(), month, date.Year())
}
