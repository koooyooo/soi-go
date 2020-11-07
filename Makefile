CLI_MOD="cmd/cli/soi.go"
SERV_MOD="cmd/srv/soi-server.go"
CLI_EX="./soi"
SERV_EX="./soi-server"

.PHONY: run
run:
	@ go run "$(CLI_MOD)"

.PHONY: build
build:
	@ go build -o "$(CLI_EX)" "$(CLI_MOD)"

.PHONY: install
install:
	@ go install "$(CLI_MOD)"

.PHONY: clean
clean:
	@ rm "$(CLI_EX)"

.PHONY: run-server
run-server:
	@ go run "$(SERV_MOD)"

.PHONY: send-request
send-request:
	@ curl -X POST -d '{"name":"Name","title":"Title","uri":"URI","tags":["tag1","tag2"],"created":"created","path":"/path"}' http://localhost:8080/store

.PHONY: build-server
build-server:
	@ go build -o "$(SERV_EX)" "$(SERV_MOD)"

.PHONY: push-image
push-image:
	@ docker build -t gcr.io/soi-cloud/soi-server:latest . ;\
	  docker push gcr.io/soi-cloud/soi-server:latest