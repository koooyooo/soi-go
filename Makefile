MAIN_MOD="cmd/soi2/soi2.go"
EX_FILE="./soi2"

run:
	@ go run $(MAIN_MOD)
build:
	@ go build -o $(EX_FILE) $(MAIN_MOD)
install:
	@ go install $(MAIN_MOD)
clean:
	@ rm $(EX_FILE)
