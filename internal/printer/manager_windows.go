//go:build windows

package printer

import (
	"fmt"

	"github.com/alexbrainman/printer"
)

// windowsManager implémente Manager pour Windows
type windowsManager struct{}

// newPlatformManager crée le gestionnaire Windows
func newPlatformManager() Manager {
	return &windowsManager{}
}

// ListPrinters retourne la liste des imprimantes Windows
func (m *windowsManager) ListPrinters() ([]string, error) {
	printers, err := printer.ReadNames()
	if err != nil {
		return nil, fmt.Errorf("impossible de lister les imprimantes: %w", err)
	}
	return printers, nil
}

// PrintRaw envoie des données RAW à l'imprimante Windows
func (m *windowsManager) PrintRaw(printerName string, data []byte) error {
	// Ouvrir l'imprimante
	p, err := printer.Open(printerName)
	if err != nil {
		return fmt.Errorf("impossible d'ouvrir l'imprimante %s: %w", printerName, err)
	}
	defer p.Close()

	// Démarrer un document RAW
	err = p.StartRawDocument("Billettera Ticket")
	if err != nil {
		return fmt.Errorf("impossible de démarrer le document: %w", err)
	}

	// Écrire les données
	_, err = p.Write(data)
	if err != nil {
		p.EndDocument()
		return fmt.Errorf("erreur d'écriture: %w", err)
	}

	// Terminer le document
	err = p.EndDocument()
	if err != nil {
		return fmt.Errorf("erreur de fin de document: %w", err)
	}

	return nil
}
