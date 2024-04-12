document.addEventListener("DOMContentLoaded", function () {
  document.querySelectorAll(".clickable-item").forEach((item) => {
    item.addEventListener("click", function () {
      document.getElementById("map").scrollIntoView({
        behavior: "smooth",
      });
    });
  });
});
