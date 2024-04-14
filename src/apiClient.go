package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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

	// Fetch artist relations
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
