# Groupie Tracker

## Introduction

**Groupie Tracker** est une application web permettant de suivre vos groupes de musique préférés. Découvrez les membres des groupes, leurs dates de création, le premier album, ainsi que les dates et lieux de leurs prochains concerts. Explorez les événements via une carte interactive.

## Technologies Utilisées

- **Backend:** Go (Golang)
- **Frontend:** HTML, CSS, JavaScript
- **APIs:** Google Maps JavaScript API, Groupie Trackers API pour les informations sur les artistes et leurs concerts.

## Dépendances

- Golang 1.15 ou plus.
- Un navigateur récent supportant HTML5 et JavaScript.

## Configuration de l'environnement

Assurez-vous d'avoir Go installé sur votre machine et configurez les clés API nécessaires dans vos variables d'environnement ou directement dans votre configuration de projet.

## Architecture du projet

Le projet est structuré en plusieurs composants principaux : Handlers Go pour la logique serveur, Templates HTML pour la vue utilisateur, et Scripts JavaScript pour l'interaction client.

## Installation et utilisation avec Makefile

1. **Clonez le dépôt**:

   ```bash
   git clone https://ytrack.learn.ynov.com/git/tenzo/groupie-tracker
   cd chemin_vers_le_projet
   ```

2. **Construisez le projet**:
   Utilisez le `Makefile` pour construire le projet :

   ```bash
   make build
   ```

   Cela générera un exécutable `GroupieTracker` dans votre répertoire.

3. **Exécutez l'application**:
   Pour démarrer l'application, utilisez :

   ```bash
   make run
   ```

   Accédez ensuite à `http://localhost:8080` dans votre navigateur pour utiliser l'application.

4. **Nettoyer le projet**:
   Pour supprimer l'exécutable et nettoyer votre projet :
   ```bash
   make clean
   ```

## Contribution

Les contributions à ce projet sont les bienvenues. Pour contribuer, veuillez suivre ces étapes :

1. Forkez le projet.
2. Créez votre branche de fonctionnalité (`git checkout -b feature/AmazingFeature`).
3. Committez vos changements (`git commit -m 'Add some AmazingFeature'`).
4. Poussez vers la branche (`git push origin feature/AmazingFeature`).
5. Ouvrez une Pull Request.
