function setUpMarkers(map, locations) {
  var markers = {};
  Object.keys(locations).forEach(function (locationKey) {
    var locationDetails = locations[locationKey];
    locationDetails.forEach(function (detail) {
      if (detail.coords) {
        var pos = new google.maps.LatLng(detail.coords.lat, detail.coords.lng);
        var marker = new google.maps.Marker({
          map: map,
          position: pos,
          title: locationKey,
        });
        markers[locationKey] = marker;
        console.log(`Marker set for ${locationKey}`);
      } else {
        console.log(`Invalid coordinates for ${locationKey}:`, detail.coords);
      }
    });
  });

  window.centerMapOnMarker = function (markerKey) {
    if (markers[markerKey]) {
      map.setCenter(markers[markerKey].getPosition());
      map.setZoom(15);
      console.log("Centered map on marker:", markerKey);
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
