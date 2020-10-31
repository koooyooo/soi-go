MAIN_MOD="cmd/cli/soi.go"
SERV_MOD="cmd/srv/soi-server.go"
EX_FILE="./soi"

.PHONY: run
run:
	@ go run "$(MAIN_MOD)"

.PHONY: build
build:
	@ go build -o "$(EX_FILE)" "$(MAIN_MOD)"

.PHONY: install
install:
	@ go install "$(MAIN_MOD)"

.PHONY: clean
clean:
	@ rm "$(EX_FILE)"

.PHONY: run-server
run-server:
	@ go run "$(SERV_MOD)"