package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"billettera-print-agent/internal/cashdrawer"
	"billettera-print-agent/internal/config"
	"billettera-print-agent/internal/printer"

	"github.com/gorilla/websocket"
)

// Server représente le serveur WebSocket
type Server struct {
	config   *config.Config
	upgrader websocket.Upgrader
	printer  printer.Manager
	drawer   *cashdrawer.Drawer
	clients  map[*websocket.Conn]bool
	mu       sync.RWMutex
	version  string
}

// NewServer crée un nouveau serveur WebSocket
func NewServer(cfg *config.Config, version string) *Server {
	s := &Server{
		config:  cfg,
		printer: printer.NewManager(),
		drawer:  cashdrawer.NewDrawer(),
		clients: make(map[*websocket.Conn]bool),
		version: version,
	}

	// Configurer l'upgrader avec vérification CORS
	s.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return true // Pas d'origine = connexion locale
			}
			return cfg.IsOriginAllowed(origin)
		},
	}

	return s
}

// Start démarre le serveur WebSocket
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)

	http.HandleFunc("/", s.handleConnection)

	log.Printf("Serveur WebSocket démarré sur ws://localhost%s", addr)
	return http.ListenAndServe(addr, nil)
}

// handleConnection gère une nouvelle connexion WebSocket
func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Erreur upgrade WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Enregistrer le client
	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	log.Printf("Nouveau client connecté depuis %s", r.RemoteAddr)

	// Supprimer le client à la déconnexion
	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		log.Printf("Client déconnecté: %s", r.RemoteAddr)
	}()

	// Boucle de lecture des messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Erreur lecture: %v", err)
			}
			break
		}

		// Traiter le message
		response := s.handleMessage(message)

		// Envoyer la réponse
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("Erreur serialisation réponse: %v", err)
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
			log.Printf("Erreur envoi réponse: %v", err)
			break
		}
	}
}

// handleMessage traite un message entrant et retourne la réponse
func (s *Server) handleMessage(message []byte) *Response {
	var req Request
	if err := json.Unmarshal(message, &req); err != nil {
		return NewErrorResponse("Format JSON invalide")
	}

	log.Printf("Action reçue: %s", req.Action)

	switch req.Action {
	case ActionPing:
		return s.handlePing()

	case ActionListPrinters:
		return s.handleListPrinters()

	case ActionPrintRaw:
		return s.handlePrintRaw(&req)

	case ActionOpenDrawer:
		return s.handleOpenDrawer(&req)

	default:
		return NewErrorResponse(fmt.Sprintf("Action inconnue: %s", req.Action))
	}
}

// handlePing répond à un ping
func (s *Server) handlePing() *Response {
	return &Response{
		Action:  ActionPong,
		Success: true,
		Version: s.version,
	}
}

// handleListPrinters retourne la liste des imprimantes
func (s *Server) handleListPrinters() *Response {
	printers, err := s.printer.ListPrinters()
	if err != nil {
		return NewErrorResponse(fmt.Sprintf("Erreur listing imprimantes: %v", err))
	}
	return NewPrintersListResponse(printers)
}

// handlePrintRaw envoie des données RAW à l'imprimante
func (s *Server) handlePrintRaw(req *Request) *Response {
	if req.Printer == "" {
		return NewErrorResponse("Imprimante non spécifiée")
	}
	if req.Data == "" {
		return NewErrorResponse("Données d'impression vides")
	}

	// Les données peuvent être en base64 ou en string directe
	data := []byte(req.Data)

	if err := s.printer.PrintRaw(req.Printer, data); err != nil {
		return NewErrorResponse(fmt.Sprintf("Erreur impression: %v", err))
	}

	return NewSuccessResponse(ActionPrintResult)
}

// handleOpenDrawer ouvre le tiroir-caisse
func (s *Server) handleOpenDrawer(req *Request) *Response {
	if req.Printer == "" {
		return NewErrorResponse("Imprimante non spécifiée")
	}

	if err := s.drawer.Open(req.Printer); err != nil {
		return NewErrorResponse(fmt.Sprintf("Erreur ouverture tiroir: %v", err))
	}

	return NewSuccessResponse(ActionDrawerResult)
}

// ClientCount retourne le nombre de clients connectés
func (s *Server) ClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}
