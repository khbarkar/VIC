BINARY   := vic
SRC_PATH := $(CURDIR)/$(BINARY)
BIN_PATH := $(HOME)/.local/bin/$(BINARY)
VERSION  := $(shell cat VERSION)
COMMIT   := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
LDFLAGS  := -X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)'

.PHONY: build run install update test lint clean

build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) .

run: build
	./$(BINARY)

install: build
	mkdir -p $(HOME)/.local/bin
	ln -sf $(SRC_PATH) $(BIN_PATH)
	@echo "Symlinked $(SRC_PATH) -> $(BIN_PATH)"

update:
	./$(BINARY) update

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -f $(BINARY)
