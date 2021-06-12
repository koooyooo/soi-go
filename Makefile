CLI_MOD="cmd/cli/soi.go"
SERV_MOD="cmd/srv/soi-server.go"

CLI_BIN="./soi"
SERV_BIN="./soi-server"

PROJECT_NAME=${SOI_PROJECT_NAME}
BUCKET_NAME=${SOI_BUCKET_NAME}
DOCKER_IMAGE_NAME=${SOI_DOCKER_IMAGE_NAME}

.PHONY: run
run:
	@ go run "$(CLI_MOD)"

.PHONY: build
build:
	@ go build -o "$(CLI_BIN)" "$(CLI_MOD)"

.PHONY: install
install:
	@ go install "$(CLI_MOD)"

.PHONY: clean
clean:
	@ rm "$(CLI_BIN)"

.PHONY: run-server
run-server:
	@ go run "$(SERV_MOD)"

.PHONY: send-request
send-request:
	@ curl -X POST -d '{"name":"Name","title":"Title","uri":"URI","tags":["tag1","tag2"],"created":"created","path":"/path"}' http://localhost:8080/store

.PHONY: build-server
build-server:
	@ go build -o "$(SERV_BIN)" "$(SERV_MOD)"

.PHONY: push-image
push-image:
	@ gcloud auth login; \
	  gcloud config set project $(PROJECT_NAME); \
      gcloud auth configure-docker; \
      docker rmi gcr.io/$(PROJECT_NAME)/$(DOCKER_IMAGE_NAME):latest; \
      docker build -t gcr.io/$(PROJECT_NAME)/$(DOCKER_IMAGE_NAME):latest . ;\
	  docker push gcr.io/$(PROJECT_NAME)/$(DOCKER_IMAGE_NAME):latest

.PHONY: deploy-image
deploy-image: push-image
	@ gcloud beta run deploy soi-cloud \
	--image gcr.io/$(PROJECT_NAME)/$(DOCKER_IMAGE_NAME):latest \
	--port 8080 \
	--platform=managed \
	--region=asia-northeast1 \
	--set-env-vars=SOI_BUCKET_NAME=$(BUCKET_NAME)

.PHONY: gen-grpc
gen-grpc:
	@ protoc --go_out=./soipb --go_opt=paths=source_relative --go-grpc_out=./soipb --go-grpc_opt=paths=source_relative pkg/srv/server/grpc/soi.proto
