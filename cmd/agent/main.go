// Billettera Print Agent
// Agent d'impression local pour la plateforme Billettera
//
// Ce programme fait le pont entre le navigateur (billettera.com)
// et les périphériques locaux (imprimante thermique, tiroir-caisse).
//
// Usage:
//
//	go run cmd/agent/main.go
//
// Build:
//
//	GOOS=windows go build -ldflags="-H windowsgui" -o BilletteraPrintAgent.exe cmd/agent/main.go
//	GOOS=darwin go build -o BilletteraPrintAgent cmd/agent/main.go
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"billettera-print-agent/internal/config"
	"billettera-print-agent/internal/systray"
	"billettera-print-agent/internal/websocket"
)

const (
	// Version de l'agent
	Version = "1.0.0"

	// Nom de l'application
	AppName = "Billettera Print Agent"
)

func main() {
	// Afficher la bannière en mode console
	printBanner()

	// Charger la configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Avertissement: erreur configuration: %v", err)
	}

	log.Printf("Configuration chargée: port=%d, origines=%v", cfg.Port, cfg.AllowedOrigins)

	// Créer le serveur WebSocket
	server := websocket.NewServer(cfg, Version)

	// Canal pour signaler l'arrêt
	quit := make(chan struct{})

	// Démarrer le serveur WebSocket dans une goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Erreur serveur WebSocket: %v", err)
			close(quit)
		}
	}()

	// Gérer les signaux d'arrêt (Ctrl+C)
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Signal d'arrêt reçu")
		close(quit)
	}()

	// Démarrer le system tray
	// Note: Sur macOS, systray.Run doit être appelé depuis le thread principal
	tray := systray.NewTray(Version, func() {
		close(quit)
	})

	log.Println("Agent prêt. L'icône apparaît dans la barre des tâches.")

	// Démarrer le system tray (bloquant)
	// Quand l'utilisateur clique sur "Quitter", cette fonction retourne
	tray.Run()

	log.Println("Arrêt de l'agent...")
}

// printBanner affiche la bannière de démarrage
func printBanner() {
	banner := `
╔═══════════════════════════════════════════════════════╗
║                                                       ║
║      BILLETTERA PRINT AGENT  v%-10s             ║
║                                                       ║
║      Agent d'impression pour billettera.com           ║
║                                                       ║
╚═══════════════════════════════════════════════════════╝
`
	fmt.Printf(banner, Version)
	fmt.Println()
}
