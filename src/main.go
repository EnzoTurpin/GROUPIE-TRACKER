package main

import (
	"log"
	"net/http"
)

// Fonctions utilitaires pour les couleurs
const (
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

func colorize(color string, message string) string {
	return color + message + colorReset
}

func main() {
	// Configuration des routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artists", artistsHandler)
	http.HandleFunc("/artist/", artistDetailHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Démarrage du serveur
	log.Println(colorize(colorGreen, "[SERVER_READY] Serveur démarré sur :8080"))
	log.Println(colorize(colorYellow, "[SERVER_INFO] Appuyez sur Ctrl+C pour arrêter le serveur."))

	// La fonction ListenAndServe démarre le serveur HTTP et bloque jusqu'à ce qu'une erreur se produise ou que le serveur soit arrêté
	log.Fatal(http.ListenAndServe(":8080", nil))
}
