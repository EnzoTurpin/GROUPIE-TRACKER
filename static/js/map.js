function initMap() {
  var mapOptions = {
    zoom: 4,
    center: { lat: -25.344, lng: 131.036 }, // Centrez la carte selon vos besoins
  };
  var map = new google.maps.Map(document.getElementById("map"), mapOptions);
  var geocoder = new google.maps.Geocoder();

  artistLocations.forEach(function (locationName) {
    geocoder.geocode({ address: locationName }, function (results, status) {
      if (status === "OK") {
        var marker = new google.maps.Marker({
          map: map,
          position: results[0].geometry.location,
        });
      } else {
        console.error(
          "Geocode was not successful for the following reason: " + status
        );
      }
    });
  });
}
