CC = go
APP_NAME = julien
BUILD_DIR = build
BUILD_NAME = app
INSTALL_PATH = /usr/bin/$(APP_NAME)
BUILD_CMD = $(CC) build
BUILD_PATH = ./$(BUILD_DIR)/$(BUILD_NAME)

julien:
	$(BUILD_CMD) -o $(BUILD_PATH)
	mv $(BUILD_PATH) .

install:
	$(BUILD_CMD) -o $(BUILD_PATH)
	mv $(BUILD_PATH) $(INSTALL_PATH)

.PHONY: julien

