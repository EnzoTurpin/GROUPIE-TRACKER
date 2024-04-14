// Fonction pour configurer les marqueurs sur la carte.
function setUpMarkers(map, locations) {
  var markers = {}; // Stocke les marqueurs

  // Boucle sur chaque emplacement pour créer un marqueur.
  Object.keys(locations).forEach(function (locationKey) {
    var details = locations[locationKey];
    var coords = details[0].split(",");
    var lat = parseFloat(coords[0]);
    var lng = parseFloat(coords[1]);
    var pos = new google.maps.LatLng(lat, lng);
    var marker = new google.maps.Marker({
      map: map,
      position: pos,
      title: locationKey,
    });
    markers[locationKey] = marker;
  });

  // Fonction pour centrer la carte sur un marqueur.
  window.centerMapOnMarker = function (markerKey) {
    if (markers[markerKey]) {
      map.setCenter(markers[markerKey].getPosition());
      map.setZoom(15);
    } else {
      console.error("Marker not found for key:", markerKey);
    }
  };
}

// Fonction pour initialiser la carte.
window.initMap = function () {
  var firstLocationCoordsElement = document.getElementById(
    "firstLocationCoords"
  );
  if (!firstLocationCoordsElement) {
    console.error("Element with ID 'firstLocationCoords' not found.");
    return;
  }

  var firstLocation = JSON.parse(firstLocationCoordsElement.textContent);
  var mapOptions = {
    zoom: 8,
    center: new google.maps.LatLng(firstLocation.lat, firstLocation.lng),
  };
  var map = new google.maps.Map(document.getElementById("map"), mapOptions);

  // Configure les marqueurs si `window.artistLocations` est défini et est un objet.
  if (
    typeof window.artistLocations === "object" &&
    window.artistLocations !== null
  ) {
    setUpMarkers(map, window.artistLocations);
  } else {
    console.error(
      "artistLocations is not defined or not an object. Received:",
      window.artistLocations
    );
  }
};

// Charge l'API Google Maps si elle n'est pas déjà chargée.
function loadGoogleMapsAPI() {
  if (!window.google || !window.google.maps) {
    var script = document.createElement("script");
    script.src =
      "https://maps.googleapis.com/maps/api/js?key=AIzaSyAh2P2kno4spZ-ERly8TUG4avTK90Z9zrU&callback=initMap";
    script.defer = true;
    document.head.appendChild(script);
  } else {
    initMap();
  }
}

// Appelée lorsque la fenêtre est entièrement chargée.
window.onload = loadGoogleMapsAPI;
