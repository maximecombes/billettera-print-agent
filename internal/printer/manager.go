// Package printer gère l'impression sur les imprimantes thermiques
package printer

// Manager définit l'interface pour gérer les imprimantes
type Manager interface {
	// ListPrinters retourne la liste des imprimantes disponibles
	ListPrinters() ([]string, error)

	// PrintRaw envoie des données RAW (ESC/POS) à l'imprimante
	PrintRaw(printerName string, data []byte) error
}

// NewManager crée un nouveau gestionnaire d'imprimantes
// L'implémentation dépend du système d'exploitation (build tags)
func NewManager() Manager {
	return newPlatformManager()
}
