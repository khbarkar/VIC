BINARY   := vic
SRC_PATH := $(CURDIR)/$(BINARY)
BIN_PATH := $(HOME)/.local/bin/$(BINARY)

.PHONY: build run install test lint clean

build:
	go build -o $(BINARY) .

run: build
	./$(BINARY)

install: build
	mkdir -p $(HOME)/.local/bin
	ln -sf $(SRC_PATH) $(BIN_PATH)
	@echo "Symlinked $(SRC_PATH) -> $(BIN_PATH)"

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -f $(BINARY)
