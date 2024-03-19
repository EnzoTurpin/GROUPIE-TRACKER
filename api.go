package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ArtistDetail struct {
	Artist    Artist
	Locations []string
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

// Ajoutez ici vos structures pour les détails de l'artiste comprenant Artist, Location, Date, Relation

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

	// Fetch dates
	datesURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%d", artistID)
	var date Date
	if err := fetchAPI(datesURL, &date); err != nil {
		return detail, err
	}
	detail.Dates = date.Dates

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
