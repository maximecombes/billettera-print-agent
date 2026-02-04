# Makefile pour Billettera Print Agent

VERSION := 1.0.0
APP_NAME := BilletteraPrintAgent
BUILD_DIR := build

# Couleurs pour le terminal
GREEN := \033[0;32m
NC := \033[0m

.PHONY: all clean deps build-windows build-darwin build-linux run test

# Build par défaut : toutes les plateformes
all: deps build-windows build-darwin

# Installer les dépendances
deps:
	@echo "$(GREEN)Installation des dépendances...$(NC)"
	go mod download
	go mod tidy

# Build Windows
build-windows:
	@echo "$(GREEN)Build Windows amd64...$(NC)"
	@mkdir -p $(BUILD_DIR)/windows
	GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui -s -w" -o $(BUILD_DIR)/windows/$(APP_NAME).exe ./cmd/agent
	@cp $(BUILD_DIR)/LISEZ-MOI-Windows.txt $(BUILD_DIR)/windows/LISEZ-MOI.txt 2>/dev/null || true
	@echo "✓ $(BUILD_DIR)/windows/$(APP_NAME).exe"

# Build macOS Intel
build-darwin:
	@echo "$(GREEN)Build macOS amd64...$(NC)"
	@mkdir -p $(BUILD_DIR)/darwin-amd64
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/darwin-amd64/$(APP_NAME) ./cmd/agent
	@cp $(BUILD_DIR)/LISEZ-MOI-macOS.txt $(BUILD_DIR)/darwin-amd64/LISEZ-MOI.txt 2>/dev/null || true
	@echo "✓ $(BUILD_DIR)/darwin-amd64/$(APP_NAME)"

# Build macOS ARM (M1/M2)
build-darwin-arm64:
	@echo "$(GREEN)Build macOS arm64...$(NC)"
	@mkdir -p $(BUILD_DIR)/darwin-arm64
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/darwin-arm64/$(APP_NAME) ./cmd/agent
	@cp $(BUILD_DIR)/LISEZ-MOI-macOS.txt $(BUILD_DIR)/darwin-arm64/LISEZ-MOI.txt 2>/dev/null || true
	@echo "✓ $(BUILD_DIR)/darwin-arm64/$(APP_NAME)"

# Build Linux (optionnel)
build-linux:
	@echo "$(GREEN)Build Linux amd64...$(NC)"
	@mkdir -p $(BUILD_DIR)/linux
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/linux/$(APP_NAME) ./cmd/agent
	@echo "✓ $(BUILD_DIR)/linux/$(APP_NAME)"

# Créer les ZIP de distribution
dist: all
	@echo "$(GREEN)Création des archives de distribution...$(NC)"
	@cd $(BUILD_DIR)/windows && zip -r ../billettera-print-agent-windows-v$(VERSION).zip .
	@cd $(BUILD_DIR)/darwin-amd64 && zip -r ../billettera-print-agent-macos-v$(VERSION).zip .
	@echo "✓ Archives créées dans $(BUILD_DIR)/"

# Lancer en mode développement
run:
	go run ./cmd/agent

# Lancer les tests
test:
	go test -v ./...

# Nettoyer les builds
clean:
	@echo "$(GREEN)Nettoyage...$(NC)"
	rm -rf $(BUILD_DIR)/windows/$(APP_NAME).exe
	rm -rf $(BUILD_DIR)/darwin-amd64
	rm -rf $(BUILD_DIR)/darwin-arm64
	rm -rf $(BUILD_DIR)/linux
	rm -f $(BUILD_DIR)/*.zip

# Afficher l'aide
help:
	@echo "Commandes disponibles :"
	@echo "  make deps          - Installer les dépendances Go"
	@echo "  make build-windows - Compiler pour Windows"
	@echo "  make build-darwin  - Compiler pour macOS Intel"
	@echo "  make build-darwin-arm64 - Compiler pour macOS ARM"
	@echo "  make all           - Compiler Windows + macOS"
	@echo "  make dist          - Créer les ZIP de distribution"
	@echo "  make run           - Lancer en mode développement"
	@echo "  make test          - Lancer les tests"
	@echo "  make clean         - Nettoyer les builds"
