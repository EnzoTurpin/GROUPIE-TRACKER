package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// fetchAPI réalise une requête GET à l'URL spécifiée et déserialise la réponse JSON dans la cible.
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

// fetchArtists récupère la liste des artistes depuis l'API et convertit les dates du premier album en français.
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

// fetchArtistDetails récupère les détails d'un artiste spécifique par son ID,
// incluant les informations de base, les relations de dates et de lieux, et les coordonnées du premier lieu.
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

	locationURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%d", artistID)
	var relations Relation
	if err := fetchAPI(locationURL, &relations); err != nil {
		log.Printf("Error fetching relations for artist ID %d: %v", artistID, err)
		return detail, err
	}

	formattedDatesLocations := make(map[string][]string)
	for location, dates := range relations.DatesLocations {
		formattedLocationName := formatLocationName(location)
		formattedDates := []string{}
		for _, date := range dates {
			formattedDate := formatDate(date)
			formattedDates = append(formattedDates, formattedDate)
		}
		formattedDatesLocations[formattedLocationName] = formattedDates
	}
	detail.DatesLocations = formattedDatesLocations

	detail.MapLinks = make([]string, len(detail.DatesLocations))
	i := 0
	for location := range detail.DatesLocations {
		detail.MapLinks[i] = generateGoogleMapsLink(location)
		i++
	}

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
