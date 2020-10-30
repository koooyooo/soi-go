MAIN_MOD="cmd/soi2/soi2.go"
EX_FILE="./soi2"

.PHONY: run
run:
	@ go run $(MAIN_MOD)

.PHONY: build
build:
	@ go build -o $(EX_FILE) $(MAIN_MOD)

.PHONY: install
install:
	@ go install $(MAIN_MOD)

.PHONY: clean
clean:
	@ rm $(EX_FILE)
