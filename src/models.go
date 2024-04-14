package main

type ArtistDetail struct {
	Artist              Artist
	DatesLocations      map[string][]string
	MapLinks            []string
	FirstLocationCoords struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}
}

// Artist définit la structure de base des données de l'artiste.
type Artist struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Image             string   `json:"image"`
	Members           []string `json:"members"`
	CreationDate      int      `json:"creationDate"`
	FirstAlbum        string   `json:"firstAlbum"`
	FirstAlbumDateStr string
}

// Location représente une structure pour identifier les lieux associés à un événement.
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

// Date représente une structure pour les dates d'événements.
type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// Relation contient des relations entre des dates et des lieux pour un artiste.
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
