# Projet Go – FileOps / WebOps / ProcOps

Projet individuel réalisé dans le cadre du M1 DevOps.  
Application console en Go permettant la manipulation de fichiers, l’analyse de contenu web (Wikipédia) et la gestion de processus système.

---

## Objectifs du projet

- Manipuler des fichiers texte (lecture, analyse, filtrage, fusion)
- Lire une configuration externe
- Récupérer et analyser du contenu web (Wikipédia)
- Interagir avec les processus système (Windows / macOS)
- Gérer les erreurs proprement
- Organiser les sorties dans un dossier dédié (`out/`)

---

## Structure du projet

projet-go-exam/
│
├── cmd/
│ └── main.go
├── data/
│ └── input.txt
├── out/
│ ├── filtered.txt
│ ├── filtered_not.txt
│ ├── head.txt
│ ├── tail.txt
│ ├── index.txt
│ ├── report.txt
│ ├── merged.txt
│ └── wiki_*.txt
├── config.txt
├── go.mod
└── go.sum


---

## Configuration

Le programme utilise un fichier `config.txt` au démarrage :

default_file=data/input.txt
base_dir=data
out_dir=out
default_ext=.txt


Les lignes vides ou commençant par `#` sont ignorées.  
Des valeurs par défaut sont utilisées si une clé est absente.

---

## Lancement du programme

Prérequis :
- Go installé (version 1.20+ recommandée)
- Git (pour cloner le projet)

Commande :

```bash
go run ./cmd

Fonctionnalités implémentées

Choix A – Analyse d’un fichier
Informations fichier (taille, date, nombre de lignes)
Statistiques sur les mots (hors numériques)
Filtrage par mot-clé
Génération de fichiers filtrés
Head / Tail (N premières / dernières lignes)

Choix B – Analyse multi-fichiers
Listing des fichiers .txt
Rapport global
Indexation des fichiers
Fusion de fichiers texte

Choix C – WebOps (Wikipédia)
Téléchargement d’un article Wikipédia
Extraction du texte des paragraphes
Analyse des mots
Sauvegarde dans out/wiki_<article>.txt
User-Agent personnalisé pour éviter les erreurs 403

Choix D – ProcOps
Listing des processus système (top N)
Recherche par mot-clé
Suppression sécurisée d’un processus
Compatibilité Windows / macOS (détection OS)


Niveau visé
Niveau 16/20 – FileOps + WebOps + ProcOps

Auteur
Projet réalisé par Gwenn Dissanda
M1 DevOps