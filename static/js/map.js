window.initMap = function () {
  var firstLocationCoordsElement = document.getElementById(
    "firstLocationCoords"
  );
  if (!firstLocationCoordsElement) {
    console.error("Element with ID 'firstLocationCoords' not found.");
    return;
  }

  // Assurez-vous que les coordonnées de la première localisation sont correctement formatées en JSON
  try {
    var firstLocation = JSON.parse(firstLocationCoordsElement.textContent);
    // Si le parsing est réussi, log les coordonnées pour vérifier
    console.log("First location coordinates:", firstLocation);
  } catch (e) {
    console.error("Error parsing first location coordinates:", e);
    return;
  }

  // Initialisez la carte avec ces coordonnées
  var mapOptions = {
    zoom: 8,
    center: new google.maps.LatLng(firstLocation.lat, firstLocation.lng),
  };

  var map = new google.maps.Map(document.getElementById("map"), mapOptions);
  var geocoder = new google.maps.Geocoder();
  var markers = {};

  // Utilisez artistLocations qui devrait être défini globalement
  if (Array.isArray(window.artistLocations)) {
    window.artistLocations.forEach(function (location, index) {
      geocoder.geocode({ address: location }, function (results, status) {
        if (status === "OK") {
          var marker = new google.maps.Marker({
            map: map,
            position: results[0].geometry.location,
          });
          markers[index] = marker;
        } else {
          console.error(
            "Geocode was not successful for the following reason: " + status
          );
        }
      });
    });
  } else {
    console.error("artistLocations is not defined or not an array.");
  }

  // La fonction pour centrer la carte sur un marqueur quand un lieu est cliqué
  window.centerMapOnMarker = function (markerKey) {
    if (markers[markerKey]) {
      map.setCenter(markers[markerKey].getPosition());
      map.setZoom(15);
    }
  };
};

// Vérifiez que l'API Google Maps est chargée avant de définir initMap
function loadGoogleMapsAPI() {
  var existingScript = document.querySelector(
    'script[src^="https://maps.googleapis.com/maps/api/js"]'
  );
  if (!existingScript) {
    var script = document.createElement("script");
    script.src =
      "https://maps.googleapis.com/maps/api/js?key=AIzaSyAh2P2kno4spZ-ERly8TUG4avTK90Z9zrU&callback=initMap";
    script.async = true;
    script.defer = true;
    document.head.appendChild(script);
  }
}

// L'événement load assure que l'API Maps ne se charge pas avant que la page ne soit entièrement chargée
window.onload = loadGoogleMapsAPI;
