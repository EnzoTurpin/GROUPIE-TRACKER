package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func formatDateToFrench(dateStr string) string {
	if dateStr == "" {
		return "Date non spécifiée"
	}
	parsedDate, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		log.Printf("Failed to parse date '%s': %v", dateStr, err)
		return "Format de date invalide"
	}
	frenchDate := fmt.Sprintf("%02d %s %d", parsedDate.Day(), monthsFrench[parsedDate.Month().String()], parsedDate.Year())
	return frenchDate
}

func formatCreationYear(year int) string {
	return strconv.Itoa(year)
}

func formatDate(dateStr string) string {
	date, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		log.Printf("Failed to parse date '%s': %v", dateStr, err)
		return dateStr
	}

	day := date.Format("02")
	month := monthsFrench[date.Month().String()]
	year := date.Format("2006")
	return fmt.Sprintf("%s %s %s", day, month, year)
}

func formatLocationName(location string) string {
	location = strings.Replace(location, "_", " ", -1)

	parts := strings.Split(location, "-")
	if len(parts) == 2 {
		parts[0] = strings.Title(parts[0])
		parts[1] = strings.ToUpper(parts[1])
	}

	return strings.Join(parts, ", ")
}

func generateGoogleMapsLink(locationName string) string {
	baseUrl := "https://www.google.com/maps/search/?api=1&query="
	return baseUrl + url.QueryEscape(locationName)
}

func getFirstKey(m map[string][]string) string {
	for k := range m {
		return k
	}
	return ""
}
