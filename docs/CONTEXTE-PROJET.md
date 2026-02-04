# Contexte du Projet Billettera Print Agent

> Document de référence expliquant pourquoi ce projet existe et les décisions techniques prises.

---

## Origine du besoin

### Date de création : 30 Janvier 2026

### Problématique initiale

La plateforme **Billettera** (https://billettera.com) est une solution SaaS de billetterie qui inclut un **Point de Vente (POS)** pour la vente au guichet.

Le POS nécessite :
1. **Impression de tickets thermiques** (format ESC/POS)
2. **Ouverture automatique du tiroir-caisse**

**Problème** : Les navigateurs web sont sandboxés et n'ont pas accès aux périphériques matériels.

### Solution précédente : QZ Tray

Billettera utilisait **QZ Tray**, un logiciel tiers open-source.

**Inconvénients de QZ Tray :**

| Problème | Impact business |
|----------|-----------------|
| Version gratuite : popups à chaque impression | Expérience utilisateur dégradée |
| Version payante : ~500$/an/site | Coût répercuté sur les clients |
| Dépendance tierce | Pas de contrôle sur les mises à jour |
| Pas de branding | Logo QZ visible |

### Décision : Développer notre propre agent

**Avantages :**
- ✅ Pas de coût de licence
- ✅ Pas de popups
- ✅ Contrôle total
- ✅ Branding Billettera
- ✅ Personnalisation possible

---

## Décisions techniques

### Langage : Go

Après analyse des options (Go, Rust, C#, Electron, Python), **Go** a été choisi pour :

| Critère | Go | Electron | C# |
|---------|-----|----------|-----|
| Exécutable unique | ✅ | ❌ (~150MB) | ❌ (.NET requis) |
| Cross-platform | ✅ | ✅ | ❌ (Windows only) |
| Taille finale | ~5-15MB | ~150MB | ~50MB |
| Facilité | ✅ | ✅ | ✅ |

### Architecture

```
Navigateur (billettera.com)
         │
         │ WebSocket (ws://localhost:19195)
         ▼
Billettera Print Agent (Go)
         │
         ├── Imprimante thermique (ESC/POS)
         └── Tiroir-caisse
```

### Port de communication

**Port 19195** choisi car :
- Non utilisé par d'autres services courants
- Facile à retenir (19-19-5)
- Supérieur à 1024 (pas besoin de droits admin)

### Protocole

Communication via **WebSocket** avec messages **JSON** :

```json
// Requête
{ "action": "open_drawer", "printer": "EPSON TM-T20III" }

// Réponse
{ "action": "drawer_result", "success": true }
```

### Sécurité

| Mesure | Implémentation |
|--------|----------------|
| CORS | Seules les origines `billettera.com` et `localhost:8000` sont acceptées |
| Signature du code | Non (accepté : "Éditeur inconnu" sur Windows, "Clic droit > Ouvrir" sur macOS) |

---

## Plateformes cibles

### Priorité 1 : Windows

Estimation de la répartition des clients Billettera :
- **70% Windows** (majorité des commerces)

### Priorité 2 : macOS

- **30% macOS** (certains professionnels)
- Nécessite une compilation sur Mac ou via GitHub Actions (voir ci-dessous)

### Priorité 3 : Linux (futur)

À considérer si demande client.

---

## État du développement (30 Janvier 2026)

### Ce qui a été fait

| Élément | Statut | Notes |
|---------|--------|-------|
| Structure du projet Go | ✅ Terminé | Architecture modulaire |
| Serveur WebSocket | ✅ Terminé | Port 19195, vérification CORS |
| Gestion des messages JSON | ✅ Terminé | Actions : ping, list_printers, print_raw, open_drawer |
| Détection imprimantes Windows | ✅ Terminé | Via `alexbrainman/printer` |
| Détection imprimantes macOS/Linux | ✅ Terminé | Via commande `lpstat` (CUPS) |
| Impression RAW (ESC/POS) | ✅ Terminé | Windows + macOS + Linux |
| Ouverture tiroir-caisse | ✅ Terminé | Commande ESC/POS standard |
| System Tray (icône) | ✅ Terminé | Via `getlantern/systray` |
| Configuration (config.json) | ✅ Terminé | Port, origines autorisées |
| Build Windows | ✅ Terminé | 5.5 MB, compilé et prêt |
| Build macOS | ⏳ En attente | Nécessite GitHub Actions |
| GitHub Actions workflow | ✅ Terminé | Prêt à être utilisé |

### Problème rencontré : Cross-compilation macOS

**Problème** : La librairie `systray` utilise du code natif Objective-C (CGO). Il est impossible de cross-compiler depuis Linux/Windows vers macOS.

**Erreur obtenue** :
```
# github.com/getlantern/systray
undefined: nativeLoop
undefined: registerSystray
...
```

**Solutions disponibles** :
1. **GitHub Actions** (recommandé) - Compile automatiquement sur les serveurs macOS de GitHub
2. **Compiler sur un Mac** - Si accès à un Mac disponible
3. **Supprimer systray pour macOS** - Version sans icône, ligne de commande uniquement

**Documentation** : Voir `docs/GITHUB-ACTIONS-MACOS.md` pour le guide complet.

---

## Prochaines étapes

### Immédiat
1. ⬜ **Tester l'exécutable Windows** (`build/windows/BilletteraPrintAgent.exe`)
   - Lancer l'application
   - Vérifier l'icône dans la barre des tâches
   - Tester depuis la page de test billettera.com

### Si test Windows concluant
2. ⬜ **Pusher le code sur GitHub**
3. ⬜ **Configurer GitHub Actions** pour compiler macOS automatiquement
4. ⬜ **Tester la version macOS**

### Intégration future
5. ⬜ Intégrer dans le module Encaissement de Billettera
6. ⬜ Créer la page de téléchargement sur billettera.com
7. ⬜ Documentation utilisateur finale

---

## Structure du projet

```
billettera-print-agent/
├── .github/
│   └── workflows/
│       └── build.yml              # GitHub Actions pour build auto
├── build/
│   ├── windows/
│   │   └── BilletteraPrintAgent.exe  # ✅ Prêt (5.5 MB)
│   ├── LISEZ-MOI-Windows.txt
│   └── LISEZ-MOI-macOS.txt
├── cmd/
│   └── agent/
│       └── main.go                # Point d'entrée
├── docs/
│   ├── CONTEXTE-PROJET.md         # Ce fichier
│   └── GITHUB-ACTIONS-MACOS.md    # Guide compilation macOS
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration
│   ├── websocket/
│   │   ├── server.go              # Serveur WebSocket
│   │   └── messages.go            # Structures JSON
│   ├── printer/
│   │   ├── manager.go             # Interface
│   │   ├── manager_windows.go     # Implémentation Windows
│   │   ├── manager_darwin.go      # Implémentation macOS
│   │   └── manager_linux.go       # Implémentation Linux
│   ├── cashdrawer/
│   │   └── drawer.go              # Ouverture tiroir
│   └── systray/
│       └── tray.go                # Icône barre des tâches
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Intégration avec Billettera

### Page de test créée

Une page de test a été ajoutée à la plateforme Billettera :

**URL** : `https://billettera.com/print-agent-test`

**Fichiers créés :**
- `app/Http/Controllers/PrintAgentTestController.php`
- `resources/views/print-agent-test.blade.php`
- Route dans `routes/web.php`

### Intégration future

Une fois l'agent validé, il remplacera QZ Tray dans le module Encaissement :
- `app/Http/Livewire/Encaissement.php`
- `app/Services/DirectPrintService.php`

---

## Références

### Matériel de test recommandé

| Type | Modèle recommandé |
|------|-------------------|
| Imprimante thermique | EPSON TM-T20III |
| Tiroir-caisse | Tout modèle avec câble RJ11 |

### Documentation ESC/POS

- [EPSON ESC/POS Reference](https://reference.epson-biz.com/modules/ref_escpos/index.php)
- [QZ Tray (pour référence)](https://qz.io/docs/)

### Librairies Go utilisées

| Librairie | Usage |
|-----------|-------|
| `github.com/gorilla/websocket` | Serveur WebSocket |
| `github.com/getlantern/systray` | Icône barre des tâches |
| `github.com/alexbrainman/printer` | Impression Windows |

---

## Contact

**Projet** : Billettera Print Agent
**Plateforme** : https://billettera.com
**Support** : support@billettera.com
