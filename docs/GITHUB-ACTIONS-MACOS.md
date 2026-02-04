# Compilation macOS via GitHub Actions

> Guide pas à pas pour compiler l'application macOS automatiquement sur les serveurs GitHub.

---

## Pourquoi GitHub Actions ?

La librairie `systray` (icône dans la barre des tâches) utilise du code natif Objective-C.
Pour compiler une application macOS avec cette librairie, il faut :
- Un Mac avec Xcode installé, OU
- Les serveurs macOS de GitHub Actions (gratuit pour les repos publics)

**GitHub Actions** permet de compiler automatiquement sur un vrai Mac hébergé par GitHub.

---

## Prérequis

1. Un compte GitHub (gratuit)
2. Le code source pushé sur un repository GitHub

---

## Étape 1 : Créer le repository GitHub

### Option A : Nouveau repository

1. Aller sur https://github.com/new
2. Nom du repository : `billettera-print-agent`
3. Visibilité : **Public** (gratuit) ou **Private** (2000 minutes/mois gratuites)
4. Cliquer sur "Create repository"

### Option B : Repository existant

Si vous avez déjà un repository, passez à l'étape 2.

---

## Étape 2 : Pusher le code

Dans le terminal PowerShell, depuis le dossier du projet :

```powershell
# Initialiser git (si pas déjà fait)
git init

# Ajouter tous les fichiers
git add .

# Premier commit
git commit -m "Initial commit - Billettera Print Agent v1.0.0"

# Ajouter le remote (remplacer VOTRE_USERNAME)
git remote add origin https://github.com/VOTRE_USERNAME/billettera-print-agent.git

# Pusher
git push -u origin main
```

---

## Étape 3 : Créer le workflow GitHub Actions

Le fichier `.github/workflows/build.yml` a été créé dans le projet.

Ce fichier dit à GitHub :
1. À chaque push ou manuellement, lancer le build
2. Compiler sur Windows ET macOS
3. Créer des artifacts téléchargeables

---

## Étape 4 : Lancer le build

### Automatiquement
Le build se lance automatiquement à chaque `git push`.

### Manuellement
1. Aller sur votre repository GitHub
2. Cliquer sur l'onglet **"Actions"**
3. Dans la liste à gauche, cliquer sur **"Build Billettera Print Agent"**
4. Cliquer sur le bouton **"Run workflow"** (à droite)
5. Sélectionner la branche `main`
6. Cliquer sur **"Run workflow"**

---

## Étape 5 : Récupérer les exécutables

1. Aller sur l'onglet **"Actions"**
2. Cliquer sur le workflow qui vient de terminer (icône ✓ verte)
3. En bas de la page, section **"Artifacts"**
4. Télécharger :
   - `billettera-print-agent-windows` (contient le .exe)
   - `billettera-print-agent-macos` (contient l'exécutable macOS)

---

## Structure du workflow

```
.github/
└── workflows/
    └── build.yml    ← Fichier de configuration GitHub Actions
```

Le workflow fait :

```
┌─────────────────────────────────────────────────────────────┐
│  GITHUB ACTIONS                                             │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Job 1: build-windows                                       │
│  ─────────────────────                                      │
│  • Machine: windows-latest                                  │
│  • Installe Go 1.21                                         │
│  • Compile BilletteraPrintAgent.exe                         │
│  • Upload artifact                                          │
│                                                             │
│  Job 2: build-macos (en parallèle)                          │
│  ─────────────────────────────────                          │
│  • Machine: macos-latest (vrai Mac chez GitHub)             │
│  • Installe Go 1.21                                         │
│  • Compile BilletteraPrintAgent                             │
│  • Upload artifact                                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Quotas GitHub Actions

| Type de compte | Minutes gratuites/mois | macOS (x10) |
|----------------|------------------------|-------------|
| Free (public)  | Illimité               | Illimité    |
| Free (private) | 2000 min               | 200 min     |
| Pro            | 3000 min               | 300 min     |

**Note** : Les minutes macOS comptent x10. Un build de 2 minutes = 20 minutes décomptées.

**Recommandation** : Garder le repository **public** pour avoir des builds illimités.

---

## Créer une Release (optionnel)

Pour distribuer les exécutables proprement :

1. Aller sur votre repository
2. Cliquer sur **"Releases"** (à droite)
3. Cliquer sur **"Create a new release"**
4. Tag : `v1.0.0`
5. Titre : `Billettera Print Agent v1.0.0`
6. Description : Notes de version
7. Glisser-déposer les ZIP téléchargés depuis les artifacts
8. Cliquer sur **"Publish release"**

Les clients pourront ensuite télécharger depuis :
```
https://github.com/VOTRE_USERNAME/billettera-print-agent/releases/latest
```

---

## Dépannage

### Le build échoue
1. Aller dans Actions → cliquer sur le build échoué
2. Cliquer sur le job en erreur
3. Lire les logs pour comprendre l'erreur

### "No space left on device"
Les runners GitHub ont ~14 GB. Si le projet devient trop gros, nettoyer les dépendances.

### Build macOS échoue avec "Code signing"
Normal si vous n'avez pas de certificat Apple. L'exécutable fonctionnera quand même,
l'utilisateur devra juste faire "Clic droit → Ouvrir" au premier lancement.

---

## Résumé des commandes

```bash
# Première fois
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/VOTRE_USERNAME/billettera-print-agent.git
git push -u origin main

# Mises à jour suivantes
git add .
git commit -m "Description des changements"
git push
```

À chaque push, GitHub compile automatiquement pour Windows et macOS.
