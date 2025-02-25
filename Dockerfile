FROM golang:alpine AS build
WORKDIR /go/src/github.com/koooyooo/soi-go
COPY . .
ENV GO111MODULE=on
RUN apk update && apk add ca-certificates git && rm -rf /var/cache/apk/*
RUN go mod download
RUN go mod verify
RUN GOOS=linux go build -o soi soi.go


FROM golang:alpine
WORKDIR /
COPY --from=build /go/src/github.com/koooyooo/soi-go/soi soi
ENTRYPOINT ["./soi"]
