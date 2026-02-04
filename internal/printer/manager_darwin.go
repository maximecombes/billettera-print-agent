//go:build darwin

package printer

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// darwinManager implémente Manager pour macOS
type darwinManager struct{}

// newPlatformManager crée le gestionnaire macOS
func newPlatformManager() Manager {
	return &darwinManager{}
}

// ListPrinters retourne la liste des imprimantes macOS via CUPS
func (m *darwinManager) ListPrinters() ([]string, error) {
	// Utiliser lpstat pour lister les imprimantes
	cmd := exec.Command("lpstat", "-p")
	output, err := cmd.Output()
	if err != nil {
		// Si lpstat échoue, essayer avec une liste vide
		return []string{}, nil
	}

	var printers []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// Format: "printer NOM_IMPRIMANTE is idle..."
		if strings.HasPrefix(line, "printer ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				printers = append(printers, parts[1])
			}
		}
	}

	return printers, nil
}

// PrintRaw envoie des données RAW à l'imprimante macOS via lp
func (m *darwinManager) PrintRaw(printerName string, data []byte) error {
	// Utiliser lp avec l'option -o raw pour envoyer des données brutes
	cmd := exec.Command("lp", "-d", printerName, "-o", "raw")
	cmd.Stdin = bytes.NewReader(data)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erreur impression: %w - %s", err, string(output))
	}

	return nil
}
