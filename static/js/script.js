const el = document.getElementById("texteAnime");
const text = el.innerText;
el.innerHTML = ""; // Vide l'élément pour accueillir les spans
let spanArray = [];

// Crée un span pour chaque lettre, permettant une animation individuelle
for (let i = 0; i < text.length; i++) {
  let span = document.createElement("span");
  span.innerText = text[i];
  el.appendChild(span);
  spanArray.push(span);
}

function animateLetters() {
  // Anime chaque lettre individuellement
  spanArray.forEach((span, index) => {
    // Temporise l'animation de chaque lettre pour créer un effet de vague
    setTimeout(() => {
      span.style.color = "yellow"; // Change la couleur en jaune

      // Change la couleur en blanc après un bref moment
      setTimeout(() => {
        span.style.color = "white";
      }, 150); // Durée avant que la lettre redevienne blanche
    }, index * 100); // Délai pour créer l'effet de vague à travers le texte
  });
}

// Fonction pour répéter l'animation indéfiniment
function repeatAnimation() {
  animateLetters(); // Lance l'animation des lettres
  // Calcule le délai total pour répéter l'animation, basé sur le nombre de lettres
  let totalDelay = spanArray.length * 100 + 150;
  setTimeout(repeatAnimation, totalDelay);
}

// Démarre l'animation répétée
repeatAnimation();
