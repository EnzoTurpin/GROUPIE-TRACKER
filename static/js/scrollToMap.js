// Attendez que le DOM soit chargé pour exécuter le code.
document.addEventListener("DOMContentLoaded", function () {
  // Pour chaque élément avec la classe "clickable-item", ajoutez un écouteur d'événements "click".
  document.querySelectorAll(".clickable-item").forEach((item) => {
    item.addEventListener("click", function () {
      // Faites défiler la vue vers l'élément avec l'ID "map" en douceur.
      document.getElementById("map").scrollIntoView({
        behavior: "smooth",
      });
    });
  });
});
