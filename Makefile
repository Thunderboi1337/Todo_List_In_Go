# Name of your program
APP_NAME = Todo_List

# Paths
INSTALL_DIR = /usr/local/bin
CONFIG_DIR = ~/.config/$(APP_NAME)
DATA_DIR = ~/.local/share/$(APP_NAME)
CACHE_DIR = ~/.cache/$(APP_NAME)
DESKTOP_FILE = /usr/share/applications/$(APP_NAME).desktop
#TODO_FILE = $(DATA_DIR)/todo_list.csv

# Icon (For later if I want to add Icon)
ICON_PATH = /path/to/icon.png

# Build the Go program
build:
	@echo "Building $(APP_NAME)..."
	GOOS=linux GOARCH=amd64 go build -o $(APP_NAME) main.go

# Install the binary and set up directories
install: build
	@echo "Installing $(APP_NAME)..."
	# Move the binary to /usr/local/bin
	sudo mv $(APP_NAME) $(INSTALL_DIR)/
	mkdir -p $(CONFIG_DIR)
	mkdir -p $(DATA_DIR)
	mkdir -p $(CACHE_DIR)
	

# Create a .desktop file for the application menu
desktop-entry:
	@echo "Creating desktop entry..."
	echo "[Desktop Entry]" | sudo tee $(DESKTOP_FILE)
	echo "Version=1.0" | sudo tee -a $(DESKTOP_FILE)
	echo "Name=$(APP_NAME)" | sudo tee -a $(DESKTOP_FILE)
	echo "Comment=Save your tasks" | sudo tee -a $(DESKTOP_FILE)
	echo "Exec=$(INSTALL_DIR)/$(APP_NAME)" | sudo tee -a $(DESKTOP_FILE)
	echo "Icon=$(ICON_PATH)" | sudo tee -a $(DESKTOP_FILE)
	echo "Terminal=true" | sudo tee -a $(DESKTOP_FILE)
	echo "Type=Application" | sudo tee -a $(DESKTOP_FILE)
	echo "Categories=Utility;" | sudo tee -a $(DESKTOP_FILE)
	# Make the .desktop file executable
	sudo chmod +x $(DESKTOP_FILE)

# Clean the build and installation
clean:
	@echo "Cleaning up..."
	# Remove binary
	rm -f $(INSTALL_DIR)/$(APP_NAME)
	# Remove desktop entry
	sudo rm -f $(DESKTOP_FILE)
	# Optionally remove config, data, and cache directories
	rm -rf $(CONFIG_DIR)
	rm -rf $(DATA_DIR)
	rm -rf $(CACHE_DIR)

# Default target: build and install
all: install desktop-entry

.PHONY: build install desktop-entry clean all
