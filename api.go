package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type ArtistDetail struct {
	Artist    Artist
	Locations []string
	MapLinks  []string // Liens Google Maps pour chaque location
	Dates     []string
	Relation  Relation
}

type Artist struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Image             string   `json:"image"`
	Members           []string `json:"members"`
	CreationDate      int      `json:"creationDate"`
	FirstAlbum        string   `json:"firstAlbum"`
	FirstAlbumDateStr string   // Assurez-vous que ce champ est là
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

const googleAPIKey = "AIzaSyAh2P2kno4spZ-ERly8TUG4avTK90Z9zrU"

var translateClient *translate.Client

func init() {
	// Initialisation du client de traduction
	var err error
	ctx := context.Background()
	translateClient, err = translate.NewClient(ctx, option.WithAPIKey(googleAPIKey))
	if err != nil {
		log.Fatalf("Failed to create translate client: %v", err)
	}
}

func translateText(text, targetLang string) (string, error) {
	ctx := context.Background()
	translations, err := translateClient.Translate(ctx, []string{text}, language.Make(targetLang), nil)
	if err != nil {
		return "", err
	}
	return translations[0].Text, nil
}

func fetchArtists() ([]Artist, error) {
	var artists []Artist
	if err := fetchAPI("https://groupietrackers.herokuapp.com/api/artists", &artists); err != nil {
		return nil, err
	}

	// Format the FirstAlbum date for each artist
	for i := range artists {
		artists[i].FirstAlbumDateStr = formatDateToFrench(artists[i].FirstAlbum)
	}

	return artists, nil
}

func fetchArtistDetails(artistID int) (ArtistDetail, error) {
	var detail ArtistDetail

	// Récupérer les informations de base de l'artiste
	artists, err := fetchArtists()
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		return detail, err
	}

	found := false
	for _, a := range artists {
		if a.ID == artistID {
			detail.Artist = a
			found = true
			break
		}
	}
	if !found {
		err := fmt.Errorf("artist with ID %d not found", artistID)
		log.Println(err)
		return detail, err
	}

	// Récupérer les localisations
	locationURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%d", artistID)
	var location Location
	if err := fetchAPI(locationURL, &location); err != nil {
		log.Printf("Error fetching location for artist ID %d: %v", artistID, err)
		return detail, err
	}
	detail.Locations = location.Locations

	// Traduction des localisations
	for i, location := range detail.Locations {
		translatedLocation, err := translateText(location, "fr")
		if err != nil {
			log.Printf("Failed to translate location: %v", err)
			continue
		}
		detail.Locations[i] = translatedLocation
	}

	// Récupérer les dates
	datesURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%d", artistID)
	var date Date
	if err := fetchAPI(datesURL, &date); err != nil {
		log.Printf("Error fetching dates for artist ID %d: %v", artistID, err)
		return detail, err
	}
	detail.Dates = date.Dates

	// Formater les dates
	if len(detail.Dates) > 0 {
		for i, dateStr := range detail.Dates {
			detail.Dates[i] = formatDate(dateStr)
		}
	} else {
		log.Printf("No dates found for artist ID %d", artistID)
	}

	// Générer les liens Google Maps pour chaque localisation
	detail.MapLinks = make([]string, len(detail.Locations))
	for i, locationName := range detail.Locations {
		detail.MapLinks[i] = generateGoogleMapsLink(locationName)
	}

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

func enrichArtistData(artists []Artist) []Artist {
	for i, artist := range artists {
		// Formatez la date du premier album ici et assignez à FirstAlbumDateStr
		artists[i].FirstAlbumDateStr = formatDateToFrench(artist.FirstAlbum)
	}
	return artists
}
