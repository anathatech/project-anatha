PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := 0.1.1
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=Anatha \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=anathad \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=anathacli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/anathad-manager
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/anathad
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/anathacli

build: go.sum
	go build -mod=readonly $(BUILD_FLAGS) -o build/anathad-manager ./cmd/anathad-manager
	go build -mod=readonly $(BUILD_FLAGS) -o build/anathad ./cmd/anathad
	go build -mod=readonly $(BUILD_FLAGS) -o build/anathacli ./cmd/anathacli

build-linux: go.sum
	GOOS=linux GOARCH=amd64 $(MAKE) build

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	#GO111MODULE=on go mod verify

# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify

# Clean up the build directory
clean:
	rm -rf build/

# Local validator nodes using Docker
build-docker:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: build-linux localnet-stop
	@if ! [ -f build/node0/anathad/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/anathad:Z anatha/core testnet --v 6 -o . --starting-ip-address 192.168.10.2 --keyring-backend=test --chain-id test ; fi
	docker-compose up

# Stop testnet
localnet-stop:
	docker-compose down

localnet: clean build-linux build-docker localnet-start

devnet-prepare:
	./scripts/prepare-test.sh

devnet-start:
	DAEMON_NAME=anathad DAEMON_HOME=~/.anathad DAEMON_ALLOW_DOWNLOAD_BINARIES=on DAEMON_RESTART_AFTER_UPGRADE=on \
	anathad-manager start --pruning="nothing" --log_level "main:info,state:info,x/crisis:info,x/hra:info,x/upgrade:info,x/gov:info,x/governance:info,x/treasury:info,x/distribution:debug,x/mint:debug,x/astaking:debug,*:error"

devnet: clean install devnet-prepare devnet-start

devnet-reset: clean devnet-prepare devnet-start

# Create log files
log-files:
	sudo mkdir -p /var/log/anathad && sudo touch /var/log/anathad/anathad.log && sudo touch /var/log/anathad/anathad_error.log

# Create service file
create-service:
	envsubst < ./scripts/anathad.service > ./anathad.service