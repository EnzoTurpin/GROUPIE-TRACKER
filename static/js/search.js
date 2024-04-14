document.addEventListener("DOMContentLoaded", function () {
  var searchBox = document.getElementById("searchBox");
  if (searchBox) {
    searchBox.addEventListener("keyup", function () {
      var searchQuery = this.value.toLowerCase();
      var artistCards = document.querySelectorAll(".card");
      var found = false; // Flag pour vérifier si un artiste a été trouvé

      artistCards.forEach(function (card) {
        var name = card.getAttribute("data-name").toLowerCase();
        if (name.includes(searchQuery)) {
          card.style.display = "";
          found = true;
        } else {
          card.style.display = "none";
        }
      });

      // Gérer l'affichage du message "Aucun artiste trouvé"
      var noResults = document.getElementById("noResults");
      if (found) {
        noResults.style.display = "none";
      } else {
        noResults.style.display = "block";
      }
    });
  }
});
