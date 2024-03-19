package main

import (
	"log"
	"net/http"
)

func main() {
	// Configuration des routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artists", artistsHandler)
	http.HandleFunc("/artist/", artistDetailHandler) // Ajoutez cette ligne
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Démarrage du serveur
	log.Println("Serveur démarré sur :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
