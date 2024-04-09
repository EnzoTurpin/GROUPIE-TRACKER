function initMap() {
  var mapOptions = {
    zoom: 4,
    center: { lat: -25.344, lng: 131.036 }, // Centrez la carte selon vos besoins
  };
  var map = new google.maps.Map(document.getElementById("map"), mapOptions);
  var geocoder = new google.maps.Geocoder();

  var markers = {}; // Un objet pour stocker les marqueurs avec une clé unique pour chaque location

  artistLocations.forEach(function (location, index) {
    geocoder.geocode({ address: location }, function (results, status) {
      if (status === "OK") {
        var marker = new google.maps.Marker({
          map: map,
          position: results[0].geometry.location,
        });
        // Utilisez l'index comme clé unique pour le marqueur
        markers[index] = marker;
      } else {
        console.error(
          "Geocode was not successful for the following reason: " + status
        );
      }
    });
  });

  // Ajouter ici la fonction centerMapOnMarker
  window.centerMapOnMarker = function (markerKey) {
    var marker = markers[markerKey];
    if (marker) {
      map.setCenter(marker.getPosition());
      map.setZoom(15); // Ajustez le zoom selon vos besoins
    }
  };
}
