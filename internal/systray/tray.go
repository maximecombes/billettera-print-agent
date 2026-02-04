// Package systray gère l'icône dans la barre des tâches
package systray

import (
	"log"
	"os"

	"github.com/getlantern/systray"
)

// Tray représente l'icône dans la barre des tâches
type Tray struct {
	version    string
	onQuit     func()
	statusItem *systray.MenuItem
}

// NewTray crée une nouvelle instance du system tray
func NewTray(version string, onQuit func()) *Tray {
	return &Tray{
		version: version,
		onQuit:  onQuit,
	}
}

// Run démarre le system tray (bloquant)
// Cette fonction doit être appelée depuis le thread principal sur macOS
func (t *Tray) Run() {
	systray.Run(t.onReady, t.onExit)
}

// onReady est appelé quand le system tray est prêt
func (t *Tray) onReady() {
	// Définir l'icône (utilise une icône par défaut si non trouvée)
	systray.SetIcon(getIcon())
	systray.SetTitle("Billettera")
	systray.SetTooltip("Billettera Print Agent - Connecté")

	// Menu
	systray.AddMenuItem("Billettera Print Agent v"+t.version, "").Disable()
	systray.AddSeparator()

	t.statusItem = systray.AddMenuItem("✓ Serveur actif", "Statut du serveur")
	t.statusItem.Disable()

	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quitter", "Fermer l'agent")

	// Gérer les clics sur le menu
	go func() {
		<-mQuit.ClickedCh
		log.Println("Demande de fermeture depuis le menu")
		if t.onQuit != nil {
			t.onQuit()
		}
		systray.Quit()
	}()
}

// onExit est appelé quand le system tray se ferme
func (t *Tray) onExit() {
	log.Println("System tray fermé")
}

// SetStatus met à jour le texte de statut
func (t *Tray) SetStatus(status string) {
	if t.statusItem != nil {
		t.statusItem.SetTitle(status)
	}
}

// getIcon retourne l'icône de l'application
// Retourne une icône vide si le fichier n'est pas trouvé
func getIcon() []byte {
	// Essayer de charger l'icône depuis le fichier
	iconPaths := []string{
		"icon.ico",
		"assets/icon.ico",
		"icon.png",
		"assets/icon.png",
	}

	for _, path := range iconPaths {
		if data, err := os.ReadFile(path); err == nil {
			return data
		}
	}

	// Icône par défaut (carré blanc 16x16 en ICO minimal)
	// En production, remplacer par une vraie icône
	return defaultIcon()
}

// defaultIcon retourne une icône par défaut (carré bleu 16x16 PNG)
func defaultIcon() []byte {
	// PNG 16x16 bleu simple - icône minimale mais valide
	return []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x91, 0x68, 0x36, 0x00, 0x00, 0x00,
		0x01, 0x73, 0x52, 0x47, 0x42, 0x00, 0xae, 0xce, 0x1c, 0xe9, 0x00, 0x00,
		0x00, 0x04, 0x67, 0x41, 0x4d, 0x41, 0x00, 0x00, 0xb1, 0x8f, 0x0b, 0xfc,
		0x61, 0x05, 0x00, 0x00, 0x00, 0x09, 0x70, 0x48, 0x59, 0x73, 0x00, 0x00,
		0x0e, 0xc3, 0x00, 0x00, 0x0e, 0xc3, 0x01, 0xc7, 0x6f, 0xa8, 0x64, 0x00,
		0x00, 0x00, 0x1d, 0x49, 0x44, 0x41, 0x54, 0x38, 0x4f, 0x63, 0x64, 0x60,
		0x60, 0xf8, 0xcf, 0xc0, 0xc0, 0xc0, 0xc4, 0x80, 0x04, 0x18, 0x49, 0x51,
		0x30, 0x6a, 0xc0, 0xa8, 0x01, 0x14, 0x1a, 0x00, 0x00, 0x9e, 0x0c, 0x03,
		0x01, 0xd6, 0x4a, 0xa5, 0x1e, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e,
		0x44, 0xae, 0x42, 0x60, 0x82,
	}
}
