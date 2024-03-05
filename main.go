package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Configuration des routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artists", artistsHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Démarrage du serveur
	log.Println("Serveur démarré sur :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}
