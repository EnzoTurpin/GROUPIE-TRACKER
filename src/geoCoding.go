package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

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
				}
			}
		}
		Status string
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
