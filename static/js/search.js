// Attend que le DOM soit chargé pour exécuter le code.
document.addEventListener("DOMContentLoaded", function () {
  var searchBox = document.getElementById("searchBox");
  if (searchBox) {
    // Ajoute un écouteur d'événements pour le clavier à la zone de recherche.
    searchBox.addEventListener("keyup", function () {
      // Convertit la valeur de recherche en minuscules.
      var searchQuery = this.value.toLowerCase();
      var artistCards = document.querySelectorAll(".card");
      var found = false; // Flag pour vérifier si un artiste a été trouvé

      // Parcours chaque carte d'artiste pour vérifier la correspondance avec la recherche.
      artistCards.forEach(function (card) {
        var name = card.getAttribute("data-name").toLowerCase();
        if (name.includes(searchQuery)) {
          card.style.display = ""; // Affiche la carte si elle correspond à la recherche.
          found = true;
        } else {
          card.style.display = "none"; // Masque la carte si elle ne correspond pas à la recherche.
        }
      });

      // Gère l'affichage du message "Aucun artiste trouvé".
      var noResults = document.getElementById("noResults");
      if (found) {
        noResults.style.display = "none"; // Cache le message si des résultats ont été trouvés.
      } else {
        noResults.style.display = "block"; // Affiche le message si aucun résultat n'a été trouvé.
      }
    });
  }
});
