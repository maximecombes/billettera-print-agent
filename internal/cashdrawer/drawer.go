// Package cashdrawer gère l'ouverture du tiroir-caisse
package cashdrawer

import (
	"billettera-print-agent/internal/printer"
)

// Commandes ESC/POS pour ouvrir le tiroir-caisse
// Ces commandes sont standard et fonctionnent avec la plupart des tiroirs
var (
	// Commande standard : ESC p 0 25 25 (pulse sur pin 2)
	OpenDrawerCommand = []byte{0x1B, 0x70, 0x00, 0x19, 0x19}

	// Alternative : ESC p 1 25 25 (pulse sur pin 5)
	OpenDrawerCommandAlt = []byte{0x1B, 0x70, 0x01, 0x19, 0x19}
)

// Drawer gère les opérations sur le tiroir-caisse
type Drawer struct {
	printer printer.Manager
}

// NewDrawer crée un nouveau gestionnaire de tiroir-caisse
func NewDrawer() *Drawer {
	return &Drawer{
		printer: printer.NewManager(),
	}
}

// Open ouvre le tiroir-caisse connecté à l'imprimante spécifiée
// Le tiroir-caisse est généralement connecté à l'imprimante thermique via un câble RJ11
func (d *Drawer) Open(printerName string) error {
	// Envoyer la commande ESC/POS d'ouverture du tiroir
	return d.printer.PrintRaw(printerName, OpenDrawerCommand)
}

// OpenWithAltPin ouvre le tiroir-caisse en utilisant le pin alternatif
// Utilisé si le tiroir est connecté sur le pin 5 au lieu du pin 2
func (d *Drawer) OpenWithAltPin(printerName string) error {
	return d.printer.PrintRaw(printerName, OpenDrawerCommandAlt)
}
