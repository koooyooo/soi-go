CLI_MOD="soi.go"

CLI_BIN="./soi"

PROJECT_NAME=${SOI_PROJECT_NAME}
BUCKET_NAME=${SOI_BUCKET_NAME}
DOCKER_IMAGE_NAME=${SOI_DOCKER_IMAGE_NAME}

# ------------------------
# Local - CLI
# ------------------------
.PHONY: run build install clean test
run:
	@ go run "$(CLI_MOD)"

build: clean
	@ go build -o "$(CLI_BIN)" "$(CLI_MOD)"

install: build
	@ go build -o "$(CLI_BIN)" "$(CLI_MOD)"
	@ GOBIN="$$(go env GOPATH)/bin"; \
	  cp "$(CLI_BIN)" "$$GOBIN/soi"; \
	  cp "$(CLI_BIN)" "$$GOBIN/soi-go"; \
	  echo "installed $$GOBIN/soi and $$GOBIN/soi-go"

clean:
	@ rm -f "$(CLI_BIN)"

test:
	@ go test ./...

# ------------------------
# Cloud Run (optional image of CLI)
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
