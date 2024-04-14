function setUpMarkers(map, locations) {
  var markers = {};
  Object.keys(locations).forEach(function (locationKey) {
    var details = locations[locationKey];
    var coords = details[0].split(","); // Les coordonnées sont le premier élément
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

  window.centerMapOnMarker = function (markerKey) {
    if (markers[markerKey]) {
      map.setCenter(markers[markerKey].getPosition());
      map.setZoom(15);
    } else {
      console.error("Marker not found for key:", markerKey);
    }
  };
}

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

window.onload = loadGoogleMapsAPI;
