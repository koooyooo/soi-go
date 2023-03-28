CLI_MOD="cmd/cli/soi.go"
SERV_MOD="cmd/srv/soi-server.go"

CLI_BIN="./soi"
SERV_BIN="./soi-server"

PROJECT_NAME=${SOI_PROJECT_NAME}
BUCKET_NAME=${SOI_BUCKET_NAME}
DOCKER_IMAGE_NAME=${SOI_DOCKER_IMAGE_NAME}

# ------------------------
# Local - CLI
# ------------------------
.PHONY: run build install clean
run:
	@ go run "$(CLI_MOD)"

build:
	@ go build -o "$(CLI_BIN)" "$(CLI_MOD)"

#   go install always use "main.go"'s "main" as a binary name
install:
	@ go install "$(CLI_MOD)"

clean:
	@ rm "$(CLI_BIN)"


# ------------------------
# Local - Simple
# ------------------------
.PHONY: build-simple install-simple
build-simple:
	@ go build -o soi-simple "cmd/simple/soi-simple.go"

install-simple:
	@ go install "cmd/simple/soi-simple.go"

# ------------------------
# Local - API
# ------------------------
.PHONY: run-server send-request build-server
run-server:
	@ go run "$(SERV_MOD)"

send-request:
	@ curl -X POST -d '{"name":"Name","title":"Title","uri":"URI","tags":["tag1","tag2"],"created_at":"2021-01-01T00:00:00+09:00","path":"/path"}' http://localhost:8080/api/v1/sample_user/sample_bucket/sois

build-server:
	@ go build -o "$(SERV_BIN)" "$(SERV_MOD)"

# ------------------------
# Cloud Run
# ------------------------
.PHONY: push-image deploy-image
push-image:
	@ gcloud auth login; \
	  gcloud config set project $(PROJECT_NAME); \
      gcloud auth configure-docker; \
      docker rmi gcr.io/$(PROJECT_NAME)/$(DOCKER_IMAGE_NAME):latest; \
      docker build -t gcr.io/$(PROJECT_NAME)/$(DOCKER_IMAGE_NAME):latest . ;\
	  docker push gcr.io/$(PROJECT_NAME)/$(DOCKER_IMAGE_NAME):latest

deploy-image: push-image
	@ gcloud beta run deploy soi-cloud \
	--image gcr.io/$(PROJECT_NAME)/$(DOCKER_IMAGE_NAME):latest \
	--port 8080 \
	--platform=managed \
	--region=asia-northeast1 \
	--set-env-vars=SOI_BUCKET_NAME=$(BUCKET_NAME)

# ------------------------
# gRPC
# ------------------------
.PHONY: gen-grpc
gen-grpc:
	@ protoc --go_out=./soipb --go_opt=paths=source_relative --go-grpc_out=./soipb --go-grpc_opt=paths=source_relative pkg/srv/server/grpc/soi.proto
