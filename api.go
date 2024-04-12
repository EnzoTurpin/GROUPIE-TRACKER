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
	Artist              Artist
	Locations           []string
	MapLinks            []string // Liens Google Maps pour chaque location
	Dates               []string
	Relation            Relation
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

	locationURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%d", artistID)
	var location Location
	if err := fetchAPI(locationURL, &location); err != nil {
		log.Printf("Error fetching location for artist ID %d: %v", artistID, err)
		return detail, err
	}
	detail.Locations = location.Locations

	datesURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%d", artistID)
	var dates Date
	if err := fetchAPI(datesURL, &dates); err != nil {
		log.Printf("Error fetching dates for artist ID %d: %v", artistID, err)
		return detail, err
	}
	detail.Dates = dates.Dates

	detail.MapLinks = make([]string, len(detail.Locations))
	for i, locationName := range detail.Locations {
		formattedLocationName := formatLocationName(locationName)
		detail.MapLinks[i] = generateGoogleMapsLink(formattedLocationName)
		detail.Locations[i] = formattedLocationName
	}

	if len(detail.Locations) > 0 {
		firstLocationName := detail.Locations[0]
		lat, lng, err := geocodeLocation(firstLocationName)
		if err != nil {
			log.Printf("Error geocoding location %s: %v", firstLocationName, err)
			return detail, err
		}
		detail.FirstLocationCoords.Lat = lat
		detail.FirstLocationCoords.Lng = lng
	} else {
		log.Printf("No locations found for artist ID %d", artistID)
	}

	if len(detail.Dates) > 0 {
		for i, dateStr := range detail.Dates {
			detail.Dates[i] = formatDate(dateStr)
		}
	} else {
		log.Printf("No dates found for artist ID %d", artistID)
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

func formatDate(dateStr string) string {
	cleanDateStr := strings.TrimPrefix(dateStr, "*")
	date, err := time.Parse("02-01-2006", cleanDateStr)
	if err != nil {
		log.Printf("Failed to parse date: %v", err)
		return dateStr
	}
	month := monthsFrench[date.Month().String()]
	return fmt.Sprintf("%02d %s %d", date.Day(), month, date.Year())
}
