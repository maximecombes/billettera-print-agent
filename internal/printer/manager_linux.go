//go:build linux

package printer

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// linuxManager implémente Manager pour Linux (même approche que macOS via CUPS)
type linuxManager struct{}

// newPlatformManager crée le gestionnaire Linux
func newPlatformManager() Manager {
	return &linuxManager{}
}

// ListPrinters retourne la liste des imprimantes Linux via CUPS
func (m *linuxManager) ListPrinters() ([]string, error) {
	cmd := exec.Command("lpstat", "-p")
	output, err := cmd.Output()
	if err != nil {
		return []string{}, nil
	}

	var printers []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "printer ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				printers = append(printers, parts[1])
			}
		}
	}

	return printers, nil
}

// PrintRaw envoie des données RAW à l'imprimante Linux via lp
func (m *linuxManager) PrintRaw(printerName string, data []byte) error {
	cmd := exec.Command("lp", "-d", printerName, "-o", "raw")
	cmd.Stdin = bytes.NewReader(data)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erreur impression: %w - %s", err, string(output))
	}

	return nil
}
