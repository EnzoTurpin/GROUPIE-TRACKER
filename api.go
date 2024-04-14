package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ArtistDetail struct {
	Artist              Artist
	DatesLocations      map[string][]string // Map de localisations vers listes de dates
	MapLinks            []string            // Liens Google Maps pour chaque localisation
	FirstLocationCoords struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"firstLocationCoords"`
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
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
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

	for i := range artists {
		artists[i].FirstAlbumDateStr = formatDateToFrench(artists[i].FirstAlbum)
	}

	return artists, nil
}

func geocodeLocation(locationName string) (float64, float64, error) {
	apiKey := "AIzaSyAh2P2kno4spZ-ERly8TUG4avTK90Z9zrU" // Remplacez ceci par votre clé API Google Maps
	baseUrl := "https://maps.googleapis.com/maps/api/geocode/json"
	requestUrl := fmt.Sprintf("%s?address=%s&key=%s", baseUrl, url.QueryEscape(locationName), apiKey)

	resp, err := http.Get(requestUrl)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Results []struct {
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
		Status string `json:"status"`
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, 0, err
	}

	if result.Status != "OK" || len(result.Results) == 0 {
		return 0, 0, fmt.Errorf("no results found for location: %s", locationName)
	}

	firstResult := result.Results[0]
	return firstResult.Geometry.Location.Lat, firstResult.Geometry.Location.Lng, nil
}

func fetchArtistDetails(artistID int) (ArtistDetail, error) {
	var detail ArtistDetail

	// Fetch artist base details
	artists, err := fetchArtists()
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		return detail, err
	}

	found := false
	for _, artist := range artists {
		if artist.ID == artistID {
			detail.Artist = artist
			found = true
			break
		}
	}
	if !found {
		err := fmt.Errorf("artist with ID %d not found", artistID)
		log.Println(err)
		return detail, err
	}

	// Fetch locations and dates
	locationURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%d", artistID)
	var relations Relation
	if err := fetchAPI(locationURL, &relations); err != nil {
		log.Printf("Error fetching relations for artist ID %d: %v", artistID, err)
		return detail, err
	}

	// Process and format location names and dates
	formattedDatesLocations := make(map[string][]string)
	for location, dates := range relations.DatesLocations {
		formattedLocationName := formatLocationName(location)
		formattedDates := []string{}
		for _, date := range dates {
			formattedDate := formatDate(date) // Format each date
			formattedDates = append(formattedDates, formattedDate)
		}
		formattedDatesLocations[formattedLocationName] = formattedDates
	}
	detail.DatesLocations = formattedDatesLocations

	// Generate map links and ensure names are correctly formatted
	detail.MapLinks = make([]string, len(detail.DatesLocations))
	i := 0
	for location := range detail.DatesLocations {
		detail.MapLinks[i] = generateGoogleMapsLink(location)
		i++
	}

	// Optional: Handle geocoding for the first location if necessary
	if len(detail.DatesLocations) > 0 {
		firstLocation := getFirstKey(detail.DatesLocations)
		lat, lng, err := geocodeLocation(firstLocation)
		if err != nil {
			log.Printf("Error geocoding location %s: %v", firstLocation, err)
			return detail, err
		}
		detail.FirstLocationCoords.Lat = lat
		detail.FirstLocationCoords.Lng = lng
	} else {
		log.Printf("No locations found for artist ID %d", artistID)
	}

	return detail, nil
}

// Helper function to get the first key from a map
func getFirstKey(m map[string][]string) string {
	for k := range m {
		return k
	}
	return ""
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

func formatDate(dateStr string) string {
	// Parse the date from DD-MM-YYYY
	date, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		log.Printf("Failed to parse date '%s': %v", dateStr, err)
		return dateStr // Return the original string if parsing fails
	}
	// Format the date to French format
	day := date.Format("02")
	month := monthsFrench[date.Month().String()]
	year := date.Format("2006")
	return fmt.Sprintf("%s %s %s", day, month, year)
}

func fetchArtistRelations(artistID int) (Relation, error) {
	var relations Relation
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%d", artistID)
	resp, err := http.Get(url)
	if err != nil {
		return relations, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return relations, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return relations, fmt.Errorf("failed to read response body: %v", err)
	}

	err = json.Unmarshal(body, &relations)
	if err != nil {
		return relations, fmt.Errorf("error decoding JSON from API: %v", err)
	}

	// Debugging output to check what data has been fetched
	fmt.Println("Dates and Locations for the artist:")
	for location, dates := range relations.DatesLocations {
		fmt.Println("Location:", location)
		for _, date := range dates {
			fmt.Println(" - Date:", date)
		}
	}

	return relations, nil
}
