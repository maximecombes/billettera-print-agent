// Package websocket gère le serveur WebSocket
package websocket

// Action types pour les requêtes
const (
	ActionListPrinters = "list_printers"
	ActionPrintRaw     = "print_raw"
	ActionOpenDrawer   = "open_drawer"
	ActionPing         = "ping"
)

// Action types pour les réponses
const (
	ActionPrintersList  = "printers_list"
	ActionPrintResult   = "print_result"
	ActionDrawerResult  = "drawer_result"
	ActionPong          = "pong"
	ActionError         = "error"
)

// Request représente une requête entrante
type Request struct {
	Action  string `json:"action"`
	Printer string `json:"printer,omitempty"`
	Data    string `json:"data,omitempty"` // Données ESC/POS en base64 ou string
}

// Response représente une réponse sortante
type Response struct {
	Action   string   `json:"action"`
	Success  bool     `json:"success,omitempty"`
	Printers []string `json:"printers,omitempty"`
	Error    string   `json:"error,omitempty"`
	Version  string   `json:"version,omitempty"`
}

// NewSuccessResponse crée une réponse de succès
func NewSuccessResponse(action string) *Response {
	return &Response{
		Action:  action,
		Success: true,
	}
}

// NewErrorResponse crée une réponse d'erreur
func NewErrorResponse(message string) *Response {
	return &Response{
		Action:  ActionError,
		Success: false,
		Error:   message,
	}
}

// NewPrintersListResponse crée une réponse avec la liste des imprimantes
func NewPrintersListResponse(printers []string) *Response {
	return &Response{
		Action:   ActionPrintersList,
		Success:  true,
		Printers: printers,
	}
}
